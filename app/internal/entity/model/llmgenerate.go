package model

type LLMGeneratedResponse struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type LLMGenerateRequest struct {
	Questions   []string `json:"questions"`
	CompanyName string   `json:"companyName"`
	CompanyID   string   `json:"companyId"`
	HTML        string   `json:"html"`
	Model       string   `json:"model"`
}
