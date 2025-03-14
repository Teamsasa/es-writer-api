package repository

import (
	"context"
	"fmt"
	"os"

	"es-api/app/internal/entity/model"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiRepository interface {
	GetGeminiRequest(ctx context.Context, input model.GeminiInput) (model.GeminiResponse, error)
}

type geminiRepository struct{}

func NewGeminiRepository() GeminiRepository {
	return &geminiRepository{}
}

func (r *geminiRepository) GetGeminiRequest(ctx context.Context, input model.GeminiInput) (model.GeminiResponse, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		return model.GeminiResponse{}, err
	}
	defer client.Close()

	gemModel := client.GenerativeModel(string(input.Model))
	text := input.Text

	response, err := gemModel.GenerateContent(ctx, genai.Text(text))
	if err != nil {
		return model.GeminiResponse{}, err
	}

	if len(response.Candidates) == 0 {
		return model.GeminiResponse{}, fmt.Errorf("no response generated")
	}

	var responseText string
	for _, part := range response.Candidates[0].Content.Parts {
		responseText += fmt.Sprintf("%v", part)
	}

	result := model.GeminiResponse{
		Text:         responseText,
		InputTokens:  response.UsageMetadata.PromptTokenCount,
		OutputTokens: response.UsageMetadata.CandidatesTokenCount,
	}
	return result, nil
}
