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
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

// 企業情報をまとめた構造体
type CompanyInfo struct {
	Name        string `json:"name"`
	Philosophy  string `json:"philosophy"`   // 企業理念
	CareerPath  string `json:"career_path"`  // キャリアパス
	TalentNeeds string `json:"talent_needs"` // 求める人材像
}

// Tavily APIへのリクエストパラメータ
type TavilySearchParams struct {
	Query         string `json:"query"`
	SearchDepth   string `json:"search_depth"`
	MaxResults    int    `json:"max_results"`
	ApiKey        string `json:"api_key"`
	IncludeAnswer bool   `json:"include_answer"` // AIによる要約を含めるかのフラグ
}

// Tavily APIからの検索結果
type TavilySearchResult struct {
	Results []struct {
		Title   string `json:"title"`
		URL     string `json:"url"`
		Content string `json:"content"`
	} `json:"results"`
	Answer string `json:"answer,omitempty"` // AIによる要約
}

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

func main() {
	// .envファイルを読み込む
	loadEnvFile()

	// コマンドライン引数から企業名を取得
	if len(os.Args) < 2 {
		log.Fatal("使用方法: go run search_company.go \"企業名\"")
	}
	companyName := os.Args[1]

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

func searchCompanyInfoParallel(ctx context.Context, apiKey string, companyName string) (*CompanyInfo, error) {
	info := &CompanyInfo{
		Name: companyName,
	}

	var wg sync.WaitGroup
	errCh := make(chan error, 3)
	resultCh := make(chan bool, 3)

	wg.Add(3)

	// 1. 企業理念の検索
	go func() {
		defer wg.Done()

		philosophyQuery := fmt.Sprintf("%s 企業理念 ミッション 価値観 経営理念", companyName)
		result, err := searchWithAnswer(ctx, apiKey, philosophyQuery)

		if err != nil || result == nil || result.Answer == "" {
			// バックアップクエリでリトライ
			backupQuery := fmt.Sprintf("%s 理念 目指すもの", companyName)
			result, err = searchWithAnswer(ctx, apiKey, backupQuery)
		}

		if err != nil {
			errCh <- fmt.Errorf("企業理念の検索エラー: %v", err)
			resultCh <- false
			return
		}

		if result != nil && result.Answer != "" {
			info.Philosophy = result.Answer
		}

		resultCh <- true
	}()

	// 2. キャリアパスの検索
	go func() {
		defer wg.Done()

		careerQuery := fmt.Sprintf("%s 社員 キャリアパス キャリア形成 成長機会 研修", companyName)
		result, err := searchWithAnswer(ctx, apiKey, careerQuery)

		if err != nil || result == nil || result.Answer == "" {
			// バックアップクエリ
			careerQuery = fmt.Sprintf("%s 社員インタビュー キャリア", companyName)
			result, err = searchWithAnswer(ctx, apiKey, careerQuery)
		}

		if err != nil {
			errCh <- fmt.Errorf("キャリアパスの検索エラー: %v", err)
			resultCh <- false
			return
		}

		if result != nil && result.Answer != "" {
			info.CareerPath = result.Answer
		}

		resultCh <- true
	}()

	// 3. 求める人材像の検索
	go func() {
		defer wg.Done()

		talentQuery := fmt.Sprintf("%s 求める人材 採用 人物像 採用基準", companyName)
		result, err := searchWithAnswer(ctx, apiKey, talentQuery)

		if err != nil || result == nil || result.Answer == "" {
			// バックアップクエリ
			talentQuery = fmt.Sprintf("%s 採用情報 募集要項", companyName)
			result, err = searchWithAnswer(ctx, apiKey, talentQuery)
		}

		if err != nil {
			errCh <- fmt.Errorf("求める人材像の検索エラー: %v", err)
			resultCh <- false
			return
		}

		if result != nil && result.Answer != "" {
			info.TalentNeeds = result.Answer
		}

		resultCh <- true
	}()

	// 全ての検索が完了するのを待機
	go func() {
		wg.Wait()
		close(errCh)
		close(resultCh)
	}()

	success := 0
	for _ = range resultCh {
		success++
	}

	// エラーがあっても部分的な結果を返す
	return info, nil
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

func outputResult(info *CompanyInfo) {
	fmt.Printf("# %s の企業理念と求める人材\n\n", info.Name)

	fmt.Println("## 企業理念・バリュー")
	if info.Philosophy != "" {
		paragraphs := strings.Split(info.Philosophy, "\n")
		for _, p := range paragraphs {
			if strings.TrimSpace(p) != "" {
				fmt.Println(p)
				fmt.Println()
			}
		}
	} else {
		fmt.Println("企業理念に関する情報は取得できませんでした。")
		fmt.Println()
	}

	// キャリアパス
	fmt.Println("## 社員のキャリアパス")
	if info.CareerPath != "" {
		paragraphs := strings.Split(info.CareerPath, "\n")
		for _, p := range paragraphs {
			if strings.TrimSpace(p) != "" {
				fmt.Println(p)
				fmt.Println()
			}
		}
	} else {
		fmt.Println("キャリアパスに関する情報は取得できませんでした。")
		fmt.Println()
	}

	// 求める人材
	fmt.Println("## 求める人材像")
	if info.TalentNeeds != "" {
		paragraphs := strings.Split(info.TalentNeeds, "\n")
		for _, p := range paragraphs {
			if strings.TrimSpace(p) != "" {
				fmt.Println(p)
				fmt.Println()
			}
		}
	} else {
		fmt.Println("求める人材に関する情報は取得できませんでした。")
		fmt.Println()
	}
}
