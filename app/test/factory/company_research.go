package factory

import (
	"testing"

	"es-api/app/internal/entity/model"

	"gorm.io/gorm"
)

func CreateCompanyResearch(t *testing.T, db *gorm.DB) *model.CompanyResearch {
	research := &model.CompanyResearch{
		CompanyID:   "1234567890123",
		CompanyName: "テスト株式会社",
		Philosophy:  "テスト企業理念",
		CareerPath:  "テストキャリアパス",
		TalentNeeds: "テスト求める人材像",
	}

	err := db.Create(research).Error
	if err != nil {
		t.Fatal(err)
	}

	return research
}
