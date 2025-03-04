package repository_test

import (
	"testing"

	"es-api/app/internal/entity/model"
	"es-api/app/test"
	"es-api/app/test/mock/repository"
	// "es-api/app/internal/repository/gemini"

	"github.com/stretchr/testify/assert"
)

func TestGeminiRepository_GetGeminiRequest(t *testing.T) {
	tests := []struct {
		name    string
		input   model.GeminiInput
		wantErr bool
	}{
		{
			name: "正常系：基本的なテキスト生成",
			input: model.GeminiInput{
				Model: model.GeminiFlashLite,
				Text:  "Hello, what is your model version?",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := test.SetupEchoContext("")

			// モックを使用したテスト
			mockRepo := new(mock.GeminiRepositoryMock)
			geminiRes := model.GeminiResponse{
				Text:         "I am Gemini Flash Lite model.",
				InputTokens:  5,
				OutputTokens: 7,
			}
			mockRepo.On("GetGeminiRequest", ctx, tt.input).Return(geminiRes, nil)
			response, err := mockRepo.GetGeminiRequest(ctx, tt.input)

			// 実際のリポジトリを使用したテスト
			// repo := repository.NewGeminiRepository()
			// test.LoadEnvFile(t, "../../../../.env")
			// response, err := repo.GetGeminiRequest(ctx, tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, response.Text)
			assert.NotZero(t, response.InputTokens)
			assert.NotZero(t, response.OutputTokens)
		})
	}
}
