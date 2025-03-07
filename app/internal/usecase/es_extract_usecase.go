package usecase

import (
	"fmt"
	"strings"

	"es-api/app/internal/entity/model"
	gemini "es-api/app/internal/repository/gemini"

	"github.com/labstack/echo/v4"
)

// HTMLExtractUsecase はHTMLから質問を抽出するユースケース
type HTMLExtractUsecase interface {
	ExtractQuestions(c echo.Context, html string) ([]string, error)
}

// htmlExtractUsecase はHTMLExtractUsecaseの実装
type htmlExtractUsecase struct {
	geminiRepository gemini.GeminiRepository
}

// NewHTMLExtractUsecase は新しいHTMLExtractUsecaseを作成
func NewHTMLExtractUsecase(geminiRepository gemini.GeminiRepository) HTMLExtractUsecase {
	return &htmlExtractUsecase{
		geminiRepository: geminiRepository,
	}
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

	if len(questions) == 0 {
		return nil, fmt.Errorf("質問が見つかりませんでした")
	}

	return questions, nil
}

func (u *htmlExtractUsecase) parseQuestionList(text string) []string {
	// 質問リストを取得
	questions := strings.Split(text, "*#*")

	return questions
}
