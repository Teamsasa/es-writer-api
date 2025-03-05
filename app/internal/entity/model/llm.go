package model

type LLMModel string

const (
	GeminiFlash         LLMModel = "gemini-2.0-flash"
	GeminiFlashLite     LLMModel = "gemini-2.0-flash-lite"
	GeminiFlashThinking LLMModel = "gemini-2.0-flash-thinking-exp"
)

type GeminiInput struct {
	Model LLMModel `json:"model"`
	Text  string   `json:"text"`
}

type GeminiResponse struct {
	Text         string `json:"text"`
	InputTokens  int32  `json:"input_tokens"`
	OutputTokens int32  `json:"output_tokens"`
}
