package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"es-api/app/internal/entity/model"
	db "es-api/app/internal/repository/db"
	gemini "es-api/app/internal/repository/gemini"
	tavily "es-api/app/internal/repository/tavily"
	"es-api/app/internal/repository/parse_html"

	"github.com/labstack/echo/v4"
)

type ESGenerateUsecase interface {
	GenerateES(c echo.Context, req model.ESGenerateRequest) ([]model.AnswerItem, error)
}

type esGenerateUsecase struct {
	htmlAnalyzer    parse_html.HTMLAnalyzer
	llmService      gemini.GeminiRepository
	companyInfoRepo tavily.TavilyRepository
	experienceRepo  db.ExperienceRepository
	authRepo        db.DBAuthRepository
}

// NewESGenerateUsecase は新しいESGenerateUsecaseを作成
func NewESGenerateUsecase(
	htmlAnalyzer parse_html.HTMLAnalyzer,
	llmService gemini.GeminiRepository,
	companyInfoRepo tavily.TavilyRepository,
	experienceRepo db.ExperienceRepository,
	authRepo db.DBAuthRepository,
) ESGenerateUsecase {
	return &esGenerateUsecase{
		htmlAnalyzer:    htmlAnalyzer,
		llmService:      llmService,
		companyInfoRepo: companyInfoRepo,
		experienceRepo:  experienceRepo,
		authRepo:        authRepo,
	}
}

// GenerateES はHTMLから質問を抽出し、企業情報とユーザーの経験に基づいて回答を生成
func (u *esGenerateUsecase) GenerateES(c echo.Context, req model.ESGenerateRequest) ([]model.AnswerItem, error) {
	// ヘッダーからユーザー識別情報を取得
	idp := c.Request().Header.Get("idp")
	if idp == "" {
		return nil, fmt.Errorf("認証情報(idp)が見つかりません")
	}

	// コンテキストからユーザーIDを取得
	userIDVal := c.Get("userID")
	if userIDVal == nil {
		return nil, fmt.Errorf("ユーザー認証が必要です")
	}

	userID, ok := userIDVal.(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("無効なユーザーIDです")
	}

	// ユーザーの存在確認
	exists, err := u.authRepo.FindUser(userID)
	if err != nil {
		return nil, fmt.Errorf("ユーザー情報の取得中にエラーが発生しました: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("ユーザーが見つかりません")
	}

	// 1. HTMLから質問を抽出
	questions, err := u.htmlAnalyzer.ExtractQuestions(c, req.HTML)
	if err != nil {
		return nil, fmt.Errorf("質問抽出に失敗しました: %w", err)
	}
	if len(questions) == 0 {
		return nil, fmt.Errorf("質問が見つかりませんでした")
	}

	// 2. 企業情報を取得
	ctx, cancel := context.WithTimeout(c.Request().Context(), 20*time.Second)
	defer cancel()

	companyInfo, err := u.companyInfoRepo.GetCompanyInfo(ctx, req.Company)
	if err != nil {
		// 企業情報がなくても回答を生成したいので、エラーはログに記録するのみ
		log.Printf("企業情報の取得に失敗しました: %v", err)
	}

	// 3. ユーザーの経験情報をデータベースから取得
	experience, err := u.experienceRepo.GetExperienceByUserID(c)
	if err != nil {
		// 経験情報がなくても回答を生成したいので、エラーはログに記録するのみ
		log.Printf("経験情報の取得に失敗しました: %v", err)
	}

	// 4. 質問ごとに回答を生成
	answers := make([]model.AnswerItem, 0, len(questions))
	llmModel := model.LLMModel(req.Model)

	for _, question := range questions {
		// プロンプトを作成
		prompt := buildPrompt(question, companyInfo, &experience, req.Company)

		// LLMで回答を生成
		llmInput := model.GeminiInput{
			Model: llmModel,
			Text:  prompt,
		}

		geminiResponse, err := u.llmService.GetGeminiRequest(c, llmInput)
		if err != nil {
			log.Printf("質問「%s」への回答生成に失敗: %v", question, err)
			continue
		}

		// 回答が文字制限超えているかチェックしてもいいかも？

		// 結果を追加
		answerItem := model.AnswerItem{
			Question: question,
			Answer:   geminiResponse.Text,
		}
		answers = append(answers, answerItem)
	}

	// 回答が生成できなかった場合
	if len(answers) == 0 {
		return nil, fmt.Errorf("回答を生成できませんでした")
	}

	return answers, nil
}

// buildPrompt はLLMへのプロンプトを構築
func buildPrompt(question string, companyInfo *tavily.CompanyInfo, experience *model.Experiences, companyName string) string {
	return "省略"
}
