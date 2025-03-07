package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/labstack/echo/v4"

	"es-api/app/internal/entity/model"
	repository "es-api/app/internal/repository/db"
)

type ExperienceUsecase interface {
	GetExperienceByUserID(c echo.Context) (*model.Experiences, error)
	PostExperience(c echo.Context, experience model.InputExperience) (*model.Experiences, error)
}

type experienceUsecase struct {
	er repository.ExperienceRepository
}

func NewExperienceUsecase(r repository.ExperienceRepository) ExperienceUsecase {
	return &experienceUsecase{er: r}
}

func (u *experienceUsecase) GetExperienceByUserID(c echo.Context) (*model.Experiences, error) {
	experience, err := u.er.GetExperienceByUserID(c)
	if err != nil {
		return nil, err
	}

	return &experience, nil
}

func (u *experienceUsecase) PostExperience(c echo.Context, experience model.InputExperience) (*model.Experiences, error) {
	exists, err := u.er.FindExperienceByUserID(c)
	if err != nil {
		return nil, fmt.Errorf("failed to find experience by user ID: %w", err)
	}

	if exists {
		return nil, fmt.Errorf("experience already exists")
	}

	experiences, err := u.er.PostExperience(c, experience)
	if err != nil {
		return nil, fmt.Errorf("failed to post experience: %w", err)
	}

	return &experiences, nil
}

// LLMGenerate はHTMLから質問を抽出し、企業情報とユーザーの経験に基づいて回答を生成
func (u *llmGenerateUsecase) LLMGenerate(c echo.Context, req model.LLMGenerateRequest) ([]model.LLMGeneratedResponse, error) {
	// 1. HTMLから質問を抽出
	questions, err := u.htmlExtractUsecase.ExtractQuestions(c, req.HTML)
	if err != nil {
		return nil, fmt.Errorf("質問抽出に失敗しました: %w", err)
	}
	if len(questions) == 0 {
		return nil, fmt.Errorf("質問が見つかりませんでした")
	}

	// 2. 企業情報を取得
	ctx, cancel := context.WithTimeout(c.Request().Context(), 20*time.Second)
	defer cancel()

	companyInfo, err := u.GetCompanyInfo(ctx, req.Company)
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
	answers := make([]model.LLMGeneratedResponse, 0, len(questions))

	llmModel := model.GeminiFlashLite
	if model.LLMModel(req.Model) != "" {
		llmModel = model.LLMModel(req.Model)
	}

	for _, question := range questions {
		// プロンプトを作成
		prompt := buildPrompt(question, companyInfo, &experience, req.Company)

		// LLMで回答を生成
		llmInput := model.GeminiInput{
			Model: llmModel,
			Text:  prompt,
		}

		geminiResponse, err := u.geminiRepo.GetGeminiRequest(c, llmInput)
		if err != nil {
			log.Printf("質問「%s」への回答生成に失敗: %v", question, err)
			continue
		}

		// 結果を追加
		LLMGeneratedResponse := model.LLMGeneratedResponse{
			Question: question,
			Answer:   geminiResponse.Text,
		}
		answers = append(answers, LLMGeneratedResponse)
	}

	// 回答が生成できなかった場合
	if len(answers) == 0 {
		return nil, fmt.Errorf("回答を生成できませんでした")
	}

	return answers, nil
}

// ExtractQuestions はGemini APIを使用してHTMLから質問を抽出する
func (u *htmlExtractUsecase) ExtractQuestions(c echo.Context, html string) ([]string, error) {
	// HTMLが空の場合はエラー
	if html == "" {
		return nil, fmt.Errorf("HTMLが空です")
	}

	// Gemini用のプロンプトを構築
	prompt := `以下のHTMLはエントリーシート(ES)の入力フォームです。
このHTMLから入力欄に対応する質問文を抽出し、リストアップしてください。
質問に文字数制限がある場合は、「（300字以内）」のような形で質問文の末尾に追加してください。
質問のみをシンプルに抽出し、各質問の間には*#*を必ず挿入してください。
質問文には、質問番号やIDやHTMLタグなどは含めないでください。

例：
志望動機を教えてください。（400字以内）*#*学生時代に力を入れたことは何ですか？（300字以内）*#*あなたの強みを教えてください。

以下のHTMLを分析してください:
` + html

	// Gemini APIを呼び出す
	geminiInput := model.GeminiInput{
		Model: model.GeminiFlashLite, // HTML解析は軽量モデルで十分
		Text:  prompt,
	}

	// geminiリポジトリを利用して質問抽出
	geminiResponse, err := u.geminiRepository.GetGeminiRequest(c, geminiInput)
	if err != nil {
		return nil, fmt.Errorf("質問抽出エラー: %w", err)
	}

	// Geminiの応答から質問リストを抽出
	questions := u.parseQuestionList(geminiResponse.Text)

	return questions, nil
}
