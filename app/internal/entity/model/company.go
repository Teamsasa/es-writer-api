package model

import "time"

// CompanyResearch - 企業情報のキャッシュ用モデル
type CompanyResearch struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	CompanyID   string    `json:"company_id" gorm:"unique;not null"` // gBizINFOの法人番号
	CompanyName string    `json:"company_name" gorm:"not null"`      // 企業名
	Philosophy  string    `json:"philosophy" gorm:"not null"`        // 企業理念
	CareerPath  string    `json:"career_path" gorm:"not null"`       // キャリアパス
	TalentNeeds string    `json:"talent_needs" gorm:"not null"`      // 求める人材像
	CreatedAt   time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"not null"`
}

// CompanyBasicInfo - 企業検索結果用の基本情報
type CompanyBasicInfo struct {
	CompanyID   string `json:"companyId"`
	CompanyName string `json:"companyName"`
}

// GBizInfoResponse - gBizINFO APIのレスポンス
type GBizInfoResponse struct {
	Response []struct {
		CorporateNumber string `json:"corporate_number"`
		Name            string `json:"name"`
	} `json:"hojin-infos"`
}
