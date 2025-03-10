package model

import "time"

// CompanyResearch - 企業情報のキャッシュ用モデル
type CompanyResearch struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CompanyID   string    `json:"company_id"`   // gBizINFOの法人番号
	CompanyName string    `json:"company_name"` // 企業名
	Philosophy  string    `json:"philosophy"`   // 企業理念
	CareerPath  string    `json:"career_path"`  // キャリアパス
	TalentNeeds string    `json:"talent_needs"` // 求める人材像
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CompanyBasicInfo - 企業検索結果用の基本情報
type CompanyBasicInfo struct {
	CompanyID   string `json:"corporate_number"` // 法人番号
	CompanyName string `json:"name"`             // 企業名
}

// GBizInfoResponse - gBizINFO APIのレスポンス
type GBizInfoResponse struct {
	Response []struct {
		CorporateNumber string `json:"corporate_number"`
		Name            string `json:"name"`
	} `json:"hojin-infos"`
}
