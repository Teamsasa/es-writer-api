package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"es-api/app/internal/entity/model"
	repository "es-api/app/internal/repository/db"
	"es-api/app/test"
	"es-api/app/test/factory"
)

func TestCompanyResearchRepository_FindByCompanyID(t *testing.T) {
	db := test.SetupTestDB(t, "../../../../.env")
	defer test.CleanupDB(t, db)

	repo := repository.NewCompanyResearchRepository(db)

	t.Run("異常系:企業情報が存在しない場合", func(t *testing.T) {
		ctx := context.Background()
		research, err := repo.FindByCompanyID(ctx, "non-existent-id")

		assert.NoError(t, err)
		assert.Nil(t, research)
	})

	t.Run("正常系:企業情報が存在する場合", func(t *testing.T) {
		dummyResearch := factory.CreateCompanyResearch(t, db)

		ctx := context.Background()
		research, err := repo.FindByCompanyID(ctx, dummyResearch.CompanyID)
		assert.NoError(t, err)
		assert.NotNil(t, research)
		assert.Equal(t, dummyResearch.CompanyID, research.CompanyID)
		assert.Equal(t, dummyResearch.CompanyName, research.CompanyName)
		assert.Equal(t, dummyResearch.Philosophy, research.Philosophy)
		assert.Equal(t, dummyResearch.CareerPath, research.CareerPath)
		assert.Equal(t, dummyResearch.TalentNeeds, research.TalentNeeds)
	})
}

func TestCompanyResearchRepository_Create(t *testing.T) {
	db := test.SetupTestDB(t, "../../../../.env")
	defer test.CleanupDB(t, db)

	repo := repository.NewCompanyResearchRepository(db)

	t.Run("正常系:新規企業情報の作成", func(t *testing.T) {
		newResearch := &model.CompanyResearch{
			CompanyID:   "9999999999999",
			CompanyName: "テスト株式会社2",
			Philosophy:  "テスト企業理念2",
			CareerPath:  "テストキャリアパス2",
			TalentNeeds: "テスト求める人材像2",
		}

		ctx := context.Background()
		err := repo.Create(ctx, newResearch)

		assert.NoError(t, err)
		assert.NotZero(t, newResearch.ID)
	})
}
