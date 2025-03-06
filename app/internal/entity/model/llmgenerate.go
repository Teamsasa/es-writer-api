package model

type GeneratedAnswer struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type LLMGenerateRequest struct {
	Questions []string `json:"questions"`
	Company   string   `json:"company"`
	HTML      string   `json:"html"`
	Model     string   `json:"model"`
}
