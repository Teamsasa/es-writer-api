package model

// 企業情報をまとめた構造体
type CompanyInfo struct {
	Name        string `json:"name"`
	Philosophy  string `json:"philosophy"`   // 企業理念
	CareerPath  string `json:"career_path"`  // キャリアパス
	TalentNeeds string `json:"talent_needs"` // 求める人材像
}

// Tavily APIへのリクエストパラメータ
type TavilySearchParams struct {
	Query         string `json:"query"`
	SearchDepth   string `json:"search_depth"`
	MaxResults    int    `json:"max_results"`
	ApiKey        string `json:"api_key"`
	IncludeAnswer bool   `json:"include_answer"` // AIによる要約を含めるかのフラグ
}

// Tavily APIからの検索結果
type TavilySearchResult struct {
	Results []struct {
		Title   string `json:"title"`
		URL     string `json:"url"`
		Content string `json:"content"`
	} `json:"results"`
	Answer string `json:"answer,omitempty"` // AIによる要約
}
