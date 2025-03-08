package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"es-api/app/internal/entity/model"
)

type TavilyRepository interface {
	SearchWithAnswer(ctx context.Context, apiKey string, query string) (*model.TavilySearchResult, error)
}

type tavilyRepository struct{}

func NewTavilyRepository() TavilyRepository {
	return &tavilyRepository{}
}

// 検索結果とAI要約を返す
func (r *tavilyRepository) SearchWithAnswer(ctx context.Context, apiKey string, query string) (*model.TavilySearchResult, error) {
	var result *model.TavilySearchResult
	var lastErr error

	// 最大3回のリトライを実行
	for attempt := 0; attempt < 3; attempt++ {
		if attempt > 0 {
			log.Printf("検索クエリ「%s」のリトライ中... (%d/3)", query, attempt+1)
			time.Sleep(1 * time.Second)
		}

		result, lastErr = doSearch(ctx, apiKey, query)

		// エラーがなく、結果とAI要約がある場合
		if lastErr == nil && result != nil && result.Answer != "" {
			return result, nil
		}
	}

	// 全リトライが失敗した場合でも最後の結果を返す
	if lastErr != nil {
		return nil, lastErr
	}
	return result, nil
}

// 実際の検索リクエストを実行する内部関数
func doSearch(ctx context.Context, apiKey string, query string) (*model.TavilySearchResult, error) {
	params := model.TavilySearchParams{
		Query:         query,
		SearchDepth:   "advanced",
		MaxResults:    5,
		ApiKey:        apiKey,
		IncludeAnswer: true,
	}

	jsonData, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.tavily.com/search", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %s - %s", resp.Status, string(bodyBytes))
	}

	var result model.TavilySearchResult
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, fmt.Errorf("JSON解析エラー: %v - レスポンス: %s", err, string(bodyBytes))
	}

	return &result, nil
}
