package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// 相対パスから.envファイルを読み込む
func loadEnvFile() bool {
	envPath := "../../../../.env"

	// ファイルが存在するか確認
	_, err := os.Stat(envPath)
	if os.IsNotExist(err) {
		log.Printf("Warning: .envファイルが見つかりません: %s", envPath)
		return false
	}

	err = godotenv.Load(envPath)
	if err != nil {
		log.Printf("Warning: .envファイルの読み込みに失敗: %v", err)
		return false
	}

	return true
}

// 検索結果とAI要約を返す
func searchWithAnswer(ctx context.Context, apiKey string, query string) (*TavilySearchResult, error) {
	var result *TavilySearchResult
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
func doSearch(ctx context.Context, apiKey string, query string) (*TavilySearchResult, error) {
	params := TavilySearchParams{
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

	var result TavilySearchResult
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, fmt.Errorf("JSON解析エラー: %v - レスポンス: %s", err, string(bodyBytes))
	}

	return &result, nil
}

func search_company(companyName string) {
	// APIキーを設定
	apiKey := os.Getenv("TAVILY_API_KEY")
	if apiKey == "" {
		log.Fatal("APIキーが設定されていません")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	companyInfo, err := searchCompanyInfoParallel(ctx, apiKey, companyName)
	if err != nil {
		log.Fatalf("企業情報の検索中にエラーが発生しました: %v", err)
	}

	outputResult(companyInfo)
}

func main() {
	// .envファイルを読み込む
	loadEnvFile()

	// コマンドライン引数から企業名を取得
	if len(os.Args) < 2 {
		log.Fatal("使用方法: go run search_company.go \"企業名\"")
	}
	companyName := os.Args[1]

	search_company(companyName)
}
