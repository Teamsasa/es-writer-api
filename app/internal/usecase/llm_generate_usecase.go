package usecase

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"es-api/app/internal/entity/model"
	db "es-api/app/internal/repository/db"
	gemini "es-api/app/internal/repository/gemini"
	tavily "es-api/app/internal/repository/tavily"

	"github.com/labstack/echo/v4"
)

type LLMGenerateUsecase interface {
	LLMGenerate(c echo.Context, req model.LLMGenerateRequest) ([]model.LLMGeneratedResponse, error)
}

// llmGenerateUsecase はLLMGenerateUsecaseの実装
type llmGenerateUsecase struct {
	geminiRepo      gemini.GeminiRepository
	companyInfoRepo tavily.TavilyRepository
	experienceRepo  db.ExperienceRepository
}

// NewLLMGenerateUsecase は新しいLLMGenerateUsecaseを作成
func NewLLMGenerateUsecase(
	geminiRepo gemini.GeminiRepository,
	companyInfoRepo tavily.TavilyRepository,
	experienceRepo db.ExperienceRepository,
) LLMGenerateUsecase {
	return &llmGenerateUsecase{
		geminiRepo:      geminiRepo,
		companyInfoRepo: companyInfoRepo,
		experienceRepo:  experienceRepo,
	}
}

// LLMGenerate はHTMLから質問を抽出し、企業情報とユーザーの経験に基づいて回答を生成
func (u *llmGenerateUsecase) LLMGenerate(c echo.Context, req model.LLMGenerateRequest) ([]model.LLMGeneratedResponse, error) {
	// 1. HTMLから質問を抽出
	questions, err := u.extractQuestionsFromHTML(c, req.HTML)
	if err != nil {
		return nil, fmt.Errorf("質問抽出に失敗しました: %w", err)
	}
	if len(questions) == 0 {
		return nil, fmt.Errorf("質問が見つかりませんでした")
	}

	// 2. 企業情報を取得
	ctx, cancel := context.WithTimeout(c.Request().Context(), 20*time.Second)
	defer cancel()

	companyInfo, err := u.getCompanyInfo(ctx, req.Company)
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
			defer func() {
				if r := recover(); r != nil {
					errorCh <- fmt.Errorf("質問「%s」の処理中にパニックが発生: %v", q, r)
					wg.Done()
				}
			}()

			defer wg.Done()

			type apiResponse struct {
				response model.GeminiResponse
				err      error
			}
			resultCh := make(chan apiResponse, 1)

			go func() {
				prompt := u.buildPrompt(q, companyInfo, &experience, req.Company)
				llmInput := model.GeminiInput{
					Model: llmModel,
					Text:  prompt,
				}

				resp, err := u.geminiRepo.GetGeminiRequest(c, llmInput)
				resultCh <- apiResponse{
					response: resp,
					err:      err,
				}
			}()

			select {
			case result := <-resultCh:
				if result.err != nil {
					errorCh <- fmt.Errorf("質問「%s」への回答生成に失敗: %v", q, result.err)
					return
				}

				responseCh <- indexedResponse{
					index: idx,
					resp: model.LLMGeneratedResponse{
						Question: q,
						Answer:   result.response.Text,
					},
				}

			case <-time.After(20 * time.Second):
				errorCh <- fmt.Errorf("質問「%s」の回答生成がタイムアウトしました（20秒）", q)
			}
		}(i, question)
	}

	go func() {
		wg.Wait()
		close(responseCh)
		close(errorCh)
	}()

	answers := make([]model.LLMGeneratedResponse, len(questions))
	for i := range answers {
		answers[i] = model.LLMGeneratedResponse{}
	}
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

func (u *llmGenerateUsecase) extractQuestionsFromHTML(c echo.Context, html string) ([]string, error) {
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
	geminiResponse, err := u.geminiRepo.GetGeminiRequest(c, geminiInput)
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

func (u *llmGenerateUsecase) getCompanyInfo(ctx context.Context, companyName string) (*model.CompanyInfo, error) {
	// APIキーを設定
	apiKey := os.Getenv("TAVILY_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("TAVILY_API_KEYが設定されていません")
	}

	// 企業情報を検索
	companyInfo, err := u.searchCompanyInfoParallel(ctx, apiKey, companyName)
	if err != nil {
		return nil, fmt.Errorf("企業情報の検索中にエラーが発生しました: %w", err)
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
