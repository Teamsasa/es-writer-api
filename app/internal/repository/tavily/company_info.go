package repository

import (
	"context"
	"fmt"
	"strings"
	"sync"
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
	for range resultCh {
		success++
	}

	// エラーがあっても部分的な結果を返す
	return info, nil
}

// 企業情報を文字列として整形して返す
func FormatCompanyInfo(info *CompanyInfo) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# %s の企業理念と求める人材\n\n", info.Name))

	sb.WriteString("## 企業理念・バリュー\n")
	if info.Philosophy != "" {
		paragraphs := strings.Split(info.Philosophy, "\n")
		for _, p := range paragraphs {
			if strings.TrimSpace(p) != "" {
				sb.WriteString(p + "\n\n")
			}
		}
	} else {
		sb.WriteString("企業理念に関する情報は取得できませんでした。\n\n")
	}

	// キャリアパス
	sb.WriteString("## 社員のキャリアパス\n")
	if info.CareerPath != "" {
		paragraphs := strings.Split(info.CareerPath, "\n")
		for _, p := range paragraphs {
			if strings.TrimSpace(p) != "" {
				sb.WriteString(p + "\n\n")
			}
		}
	} else {
		sb.WriteString("キャリアパスに関する情報は取得できませんでした。\n\n")
	}

	// 求める人材
	sb.WriteString("## 求める人材像\n")
	if info.TalentNeeds != "" {
		paragraphs := strings.Split(info.TalentNeeds, "\n")
		for _, p := range paragraphs {
			if strings.TrimSpace(p) != "" {
				sb.WriteString(p + "\n\n")
			}
		}
	} else {
		sb.WriteString("求める人材に関する情報は取得できませんでした。\n\n")
	}

	return sb.String()
}
