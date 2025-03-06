package model

type AnswerItem struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type ESGenerateRequest struct {
	Questions []string `json:"questions"`
	Company   string   `json:"company"`
	HTML      string   `json:"html"`
	Model     string   `json:"model"`
}
