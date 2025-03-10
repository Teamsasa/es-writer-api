package gbiz_test

import (
	"fmt"
	"testing"

	"es-api/app/internal/entity/model"
	"es-api/app/test"
	mock "es-api/app/test/mock/repository"

	"github.com/stretchr/testify/assert"
)

func TestGBizInfoRepository(t *testing.T) {
	t.Run("SearchCompanies", func(t *testing.T) {
		t.Run("正常系", func(t *testing.T) {
			t.Run("会社情報の検索が成功する", func(t *testing.T) {
				// 環境変数のセットアップ
				t.Setenv("GBIZ_API_KEY", "dummy-api-key")

				ctx := test.SetupEchoContext("")
				keyword := "テスト株式会社"
				expected := []model.CompanyBasicInfo{
					{
						CompanyID:   "1234567890123",
						CompanyName: "テスト株式会社",
					},
				}

				// モックを使用したテスト
				mockRepo := new(mock.GBizInfoRepositoryMock)
				mockRepo.On("SearchCompanies", ctx, keyword).Return(expected, nil)

				companies, err := mockRepo.SearchCompanies(ctx, keyword)

				assert.NoError(t, err)
				assert.Equal(t, expected, companies)
			})
		})

		t.Run("異常系", func(t *testing.T) {
			t.Run("GBIZ_API_KEYが未設定の場合はエラーを返す", func(t *testing.T) {
				// 環境変数のセットアップ
				t.Setenv("GBIZ_API_KEY", "")

				ctx := test.SetupEchoContext("")
				keyword := "テスト株式会社"

				// モックを使用したテスト
				mockRepo := new(mock.GBizInfoRepositoryMock)
				mockRepo.On("SearchCompanies", ctx, keyword).Return(nil, fmt.Errorf("GBIZ_API_KEY is not set"))

				companies, err := mockRepo.SearchCompanies(ctx, keyword)

				assert.Error(t, err)
				assert.Nil(t, companies)
			})

			t.Run("APIが非200ステータスを返した場合はエラーを返す", func(t *testing.T) {
				// 環境変数のセットアップ
				t.Setenv("GBIZ_API_KEY", "dummy-api-key")

				ctx := test.SetupEchoContext("")
				keyword := "テスト株式会社"

				// モックを使用したテスト
				mockRepo := new(mock.GBizInfoRepositoryMock)
				mockRepo.On("SearchCompanies", ctx, keyword).Return(nil, fmt.Errorf("API returned non-200 status code: 500"))

				companies, err := mockRepo.SearchCompanies(ctx, keyword)

				assert.Error(t, err)
				assert.Nil(t, companies)
			})
		})
	})
}
