package usecase

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"es-api/app/internal/entity/model"
	db "es-api/app/internal/repository/db"
	gemini "es-api/app/internal/repository/gemini"
	tavily "es-api/app/internal/repository/tavily"
)

type LLMGenerateUsecase interface {
	LLMGenerate(ctx context.Context, req model.LLMGenerateRequest) ([]model.LLMGeneratedResponse, error)
}

// llmGenerateUsecase はLLMGenerateUsecaseの実装
type llmGenerateUsecase struct {
	geminiRepo          gemini.GeminiRepository
	companyInfoRepo     tavily.TavilyRepository
	experienceRepo      db.ExperienceRepository
	companyResearchRepo db.CompanyResearchRepository
}

// NewLLMGenerateUsecase は新しいLLMGenerateUsecaseを作成
func NewLLMGenerateUsecase(
	geminiRepo gemini.GeminiRepository,
	companyInfoRepo tavily.TavilyRepository,
	experienceRepo db.ExperienceRepository,
	companyResearchRepo db.CompanyResearchRepository,
) LLMGenerateUsecase {
	return &llmGenerateUsecase{
		geminiRepo:          geminiRepo,
		companyInfoRepo:     companyInfoRepo,
		experienceRepo:      experienceRepo,
		companyResearchRepo: companyResearchRepo,
	}
}

// LLMGenerate はHTMLから質問を抽出し、企業情報とユーザーの経験に基づいて回答を生成
func (u *llmGenerateUsecase) LLMGenerate(ctx context.Context, req model.LLMGenerateRequest) ([]model.LLMGeneratedResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// 1. HTMLから質問を抽出
	questions, err := u.extractQuestionsFromHTML(ctx, req.HTML)
	if err != nil {
		return nil, fmt.Errorf("質問抽出に失敗しました: %w", err)
	}
	if len(questions) == 0 {
		return nil, fmt.Errorf("質問が見つかりませんでした")
	}

	// 2. 企業情報を取得
	companyInfo, err := u.getCompanyInfo(ctx, req.CompanyID, req.CompanyName)
	if err != nil {
		// 企業情報がなくても回答を生成したいので、エラーはログに記録するのみ
		log.Printf("企業情報の取得に失敗しました: %v", err)
	}

	// 3. ユーザーの経験情報をデータベースから取得
	experience, err := u.experienceRepo.GetExperienceByUserID(ctx)
	if err != nil {
		// 経験情報がなくても回答を生成したいので、エラーはログに記録するのみ
		log.Printf("経験情報の取得に失敗しました: %v", err)
	}

	// 4. 質問ごとに回答を生成
	var wg sync.WaitGroup

	type indexedResponse struct {
		index int
		resp  model.LLMGeneratedResponse
	}
	responseCh := make(chan indexedResponse, len(questions))
	errorCh := make(chan error, len(questions))

	llmModel := model.GeminiFlashLite
	if model.LLMModel(req.Model) != "" {
		llmModel = model.LLMModel(req.Model)
	}

	for i, question := range questions {
		wg.Add(1)
		go func(idx int, q string) {
			defer wg.Done()

			defer func() {
				if r := recover(); r != nil {
					errorCh <- fmt.Errorf("質問「%s」の処理中にパニックが発生: %v", q, r)
				}
			}()

			prompt := u.buildPrompt(q, companyInfo, &experience, req.CompanyName)
			llmInput := model.GeminiInput{
				Model: llmModel,
				Text:  prompt,
			}

			done := make(chan struct{})
			var resp model.GeminiResponse
			var err error

			go func() {
				resp, err = u.geminiRepo.GetGeminiRequest(ctx, llmInput)
				close(done)
			}()

			// タイムアウト検出
			select {
			case <-done:
				if err != nil {
					errorCh <- fmt.Errorf("質問「%s」への回答生成に失敗: %v", q, err)
					return
				}
				responseCh <- indexedResponse{
					index: idx,
					resp: model.LLMGeneratedResponse{
						Question: q,
						Answer:   resp.Text,
					},
				}
			case <-ctx.Done():
				errorCh <- fmt.Errorf("質問「%s」の回答生成がタイムアウトまたはキャンセルされました: %v", q, ctx.Err())
			}
		}(i, question)
	}

	go func() {
		wg.Wait()
		close(responseCh)
		close(errorCh)
	}()

	answers := make([]model.LLMGeneratedResponse, len(questions))
	validAnswers := 0

	for resp := range responseCh {
		answers[resp.index] = resp.resp
		validAnswers++
	}

	if validAnswers != len(questions) {
		if err, ok := <-errorCh; ok {
			return nil, err
		}
		return nil, fmt.Errorf("回答を生成できませんでした")
	}

	return answers, nil
}

func (u *llmGenerateUsecase) extractQuestionsFromHTML(ctx context.Context, html string) ([]string, error) {
	// HTMLが空の場合はエラー
	if html == "" {
		return nil, fmt.Errorf("HTMLが空です")
	}

	promptTemplate, err := loadPromptFromFile("extract_questions.txt")
	if err != nil {
		return nil, fmt.Errorf("プロンプトファイルの読み込みに失敗: %v", err)
	}
	prompt := promptTemplate + html

	// Gemini APIを呼び出す
	geminiInput := model.GeminiInput{
		Model: model.GeminiFlashLite, // HTML解析は軽量モデルで十分
		Text:  prompt,
	}

	// geminiリポジトリを利用して質問抽出
	geminiResponse, err := u.geminiRepo.GetGeminiRequest(ctx, geminiInput)
	if err != nil {
		return nil, fmt.Errorf("質問抽出エラー: %w", err)
	}

	// Geminiの応答から質問リストを抽出
	questions := strings.Split(geminiResponse.Text, "*#*")

	// 空白の質問をフィルタリング
	var filteredQuestions []string
	for _, q := range questions {
		trimmed := strings.TrimSpace(q)
		if trimmed != "" {
			filteredQuestions = append(filteredQuestions, trimmed)
		}
	}

	return filteredQuestions, nil
}

func (u *llmGenerateUsecase) buildPrompt(question string, companyInfo *model.CompanyInfo, experience *model.Experiences, companyName string) string {
	promptTemplate, err := loadPromptFromFile("es_generation.txt")
	if err != nil {
		log.Printf("プロンプトファイルの読み込みに失敗: %v, デフォルトのプロンプトを使用します", err)
		return fmt.Sprintf("企業%sについての質問です。一般的な応募者として回答してください。\n\n%s", companyName, question)
	}

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(promptTemplate, question))

	// 企業情報の追加
	sb.WriteString("【企業情報】\n")
	if companyInfo != nil && companyInfo.Name != "" {
		if companyInfo.Philosophy != "" {
			sb.WriteString("■企業理念・バリュー\n")
			sb.WriteString(companyInfo.Philosophy)
			sb.WriteString("\n\n")
		}

		if companyInfo.TalentNeeds != "" {
			sb.WriteString("■求める人材像\n")
			sb.WriteString(companyInfo.TalentNeeds)
			sb.WriteString("\n\n")
		}

		if companyInfo.CareerPath != "" {
			sb.WriteString("■キャリアパス\n")
			sb.WriteString(companyInfo.CareerPath)
			sb.WriteString("\n\n")
		}
	} else {
		sb.WriteString(fmt.Sprintf("%sという企業についての質問です。一般的な応募者として回答してください。\n\n", companyName))
	}

	// 応募者の経験情報の追加
	sb.WriteString("【応募者の経歴情報】\n")
	if experience != nil {
		if experience.Work != "" {
			sb.WriteString("■職務経歴\n")
			sb.WriteString(experience.Work)
			sb.WriteString("\n\n")
		}

		if experience.Skills != "" {
			sb.WriteString("■スキル\n")
			sb.WriteString(experience.Skills)
			sb.WriteString("\n\n")
		}

		if experience.SelfPR != "" {
			sb.WriteString("■自己PR\n")
			sb.WriteString(experience.SelfPR)
			sb.WriteString("\n\n")
		}

		if experience.FutureGoals != "" {
			sb.WriteString("■将来の目標\n")
			sb.WriteString(experience.FutureGoals)
			sb.WriteString("\n\n")
		}
	}

	return sb.String()
}

func loadPromptFromFile(filename string) (string, error) {
	paths := []string{
		filename,
		filepath.Join("../../internal/usecase/prompts", filename),
	}

	var content []byte
	var err error

	for _, path := range paths {
		content, err = os.ReadFile(path)
		if err == nil {
			return string(content), nil
		}
	}

	return "", err
}

func (u *llmGenerateUsecase) getCompanyInfo(ctx context.Context, companyID string, companyName string) (*model.CompanyInfo, error) {
	// キャッシュから企業情報を検索
	research, err := u.companyResearchRepo.FindByCompanyID(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("企業情報のキャッシュ検索中にエラーが発生しました: %w", err)
	}

	// キャッシュがある場合はそれを返す
	if research != nil {
		log.Printf("企業情報のキャッシュを利用します: %s", research.CompanyName)
		return &model.CompanyInfo{
			Name:        research.CompanyName,
			Philosophy:  research.Philosophy,
			CareerPath:  research.CareerPath,
			TalentNeeds: research.TalentNeeds,
		}, nil
	}

	// APIキーを設定
	apiKey := os.Getenv("TAVILY_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("TAVILY_API_KEYが設定されていません")
	}

	// 企業情報を検索
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	companyInfo, err := u.searchCompanyInfoParallel(ctx, apiKey, companyName)
	if err != nil {
		return nil, fmt.Errorf("企業情報の検索中にエラーが発生しました: %w", err)
	}

	// 検索結果をキャッシュに保存
	research = &model.CompanyResearch{
		CompanyID:   companyID,
		CompanyName: companyName,
		Philosophy:  companyInfo.Philosophy,
		CareerPath:  companyInfo.CareerPath,
		TalentNeeds: companyInfo.TalentNeeds,
	}
	if err := u.companyResearchRepo.Create(ctx, research); err != nil {
		log.Printf("企業情報のキャッシュ保存中にエラーが発生しました: %v", err)
	}

	return companyInfo, nil
}

func (u *llmGenerateUsecase) searchCompanyInfoParallel(ctx context.Context, apiKey string, companyName string) (*model.CompanyInfo, error) {
	info := &model.CompanyInfo{
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
		result, err := u.companyInfoRepo.SearchWithAnswer(ctx, apiKey, philosophyQuery)

		if err != nil || result == nil || result.Answer == "" {
			// バックアップクエリでリトライ
			backupQuery := fmt.Sprintf("%s 理念 目指すもの", companyName)
			result, err = u.companyInfoRepo.SearchWithAnswer(ctx, apiKey, backupQuery)
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
		result, err := u.companyInfoRepo.SearchWithAnswer(ctx, apiKey, careerQuery)

		if err != nil || result == nil || result.Answer == "" {
			// バックアップクエリ
			careerQuery = fmt.Sprintf("%s 社員インタビュー キャリア", companyName)
			result, err = u.companyInfoRepo.SearchWithAnswer(ctx, apiKey, careerQuery)
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
		result, err := u.companyInfoRepo.SearchWithAnswer(ctx, apiKey, talentQuery)

		if err != nil || result == nil || result.Answer == "" {
			// バックアップクエリ
			talentQuery = fmt.Sprintf("%s 採用情報 募集要項", companyName)
			result, err = u.companyInfoRepo.SearchWithAnswer(ctx, apiKey, talentQuery)
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
