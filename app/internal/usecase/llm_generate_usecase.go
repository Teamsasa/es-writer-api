package usecase

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"es-api/app/internal/entity/model"
	db "es-api/app/internal/repository/db"
	gemini "es-api/app/internal/repository/gemini"
	tavily "es-api/app/internal/repository/tavily"

	"github.com/labstack/echo/v4"
)

type LLMGenerateUsecase interface {
	LLMGenerate(c echo.Context, req model.LLMGenerateRequest) ([]model.GeneratedAnswer, error)
}

type llmGenerateUsecase struct {
	htmlExtractUsecase HTMLExtractUsecase
	llmService         gemini.GeminiRepository
	companyInfoRepo    tavily.TavilyRepository
	experienceRepo     db.ExperienceRepository
}

// NewLLMGenerateUsecase は新しいLLMGenerateUsecaseを作成
func NewLLMGenerateUsecase(
	htmlExtractUsecase HTMLExtractUsecase,
	llmService gemini.GeminiRepository,
	companyInfoRepo tavily.TavilyRepository,
	experienceRepo db.ExperienceRepository,
) LLMGenerateUsecase {
	return &llmGenerateUsecase{
		htmlExtractUsecase: htmlExtractUsecase,
		llmService:         llmService,
		companyInfoRepo:    companyInfoRepo,
		experienceRepo:     experienceRepo,
	}
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
	answers := make([]model.LLMGeneratedResponse, 0, len(questions))
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

// buildPrompt はLLMへのプロンプトを構築
func buildPrompt(question string, companyInfo *model.CompanyInfo, experience *model.Experiences, companyName string) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("あなたはエントリーシートのプロ作成者です。以下の質問に回答してください：「%s」\n\n", question))

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

	sb.WriteString("【回答作成の重要なガイドライン】\n")
	sb.WriteString("1. 日本語のみを使用し、外国語や特殊な記号（アスタリスク*, 波ダッシュ~, ハット^ など）は一切使用しないでください\n")
	sb.WriteString("2. 具体例や部活動名、数値などを示す際も、記号で囲まずに自然な文章で書いてください\n")
	sb.WriteString("3. 箇条書きやマークダウン記法は使用せず、通常の文章として書いてください\n")
	sb.WriteString("4. 「〜です。〜ます。」といった丁寧な文体を一貫して使用してください\n")
	sb.WriteString("5. 質問文や「自己PRは以下の通りです」などの形式的な表現は回答に含めないでください\n")
	sb.WriteString("6. 企業名や企業理念の過度な繰り返しを避け、自然な頻度で言及してください\n")
	sb.WriteString("7. 具体的なエピソードや経験を含め、説得力のある内容にしてください\n")
	sb.WriteString("8. 「〜と思います」「〜と考えます」「〜できると思います」など、主観的で自然な表現を適切に使ってください\n")
	sb.WriteString("9. 専門用語の使用は適度に控え、一般的な表現を心がけてください\n")
	sb.WriteString("10. 起承転結を意識した、読みやすく自然な文章構成にしてください\n")

	sb.WriteString("11. 出力前に回答の文字数をカウントして、文字数制限を超えていたら再度回答を作り直してください\n")
	sb.WriteString("12. 実際の文字数が文字数制限に近づくようにしてください（少なすぎても不自然です）\n")

	// 最終出力についての指示
	sb.WriteString("\n【最終出力の注意事項】\n")
	sb.WriteString("・特殊記号や外国語を含まない、純粋な日本語の文章を出力してください\n")
	sb.WriteString("・アスタリスク(*)や他の記号で強調したり囲んだりしないでください\n")
	sb.WriteString("・回答文のみを出力し、質問の繰り返しや前置き・説明は含めないでください\n")
	sb.WriteString("・空行は段落の区切りのみに使用し、冒頭や末尾の余計な空行は入れないでください\n")

	return sb.String()
}
