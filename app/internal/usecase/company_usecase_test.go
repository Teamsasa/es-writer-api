package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"

	"es-api/app/internal/entity/model"
	"es-api/app/internal/usecase"
	mock "es-api/app/test/mock/repository"
)

func TestCompanyUsecase_SearchCompanies(t *testing.T) {
	t.Run("正常系:検索結果がある場合", func(t *testing.T) {
		mockRepo := new(mock.GBizInfoRepositoryMock)

		expectedCompanies := []model.CompanyBasicInfo{
			{
				CompanyID:   "1234567890123",
				CompanyName: "株式会社テスト",
			},
		}

		ctx := context.Background()
		keyword := "株式会社テスト"

		mockRepo.On("SearchCompanies", testifymock.Anything, keyword).Return(expectedCompanies, nil)

		uc := usecase.NewCompanyUsecase(mockRepo)

		res, err := uc.SearchCompanies(ctx, keyword)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, expectedCompanies[0].CompanyID, res[0].CompanyID)
		assert.Equal(t, expectedCompanies[0].CompanyName, res[0].CompanyName)
		mockRepo.AssertExpectations(t)
	})

	t.Run("正常系:検索結果が空の場合", func(t *testing.T) {
		mockRepo := new(mock.GBizInfoRepositoryMock)
		ctx := context.Background()
		keyword := "存在しない会社"

		mockRepo.On("SearchCompanies", testifymock.Anything, keyword).Return([]model.CompanyBasicInfo{}, nil)

		uc := usecase.NewCompanyUsecase(mockRepo)

		res, err := uc.SearchCompanies(ctx, keyword)

		assert.NoError(t, err)
		assert.Empty(t, res)
		mockRepo.AssertExpectations(t)
	})

	t.Run("異常系:リポジトリでエラーが発生した場合", func(t *testing.T) {
		mockRepo := new(mock.GBizInfoRepositoryMock)

		ctx := context.Background()
		keyword := "エラーケース"
		expectedErr := errors.New("repository error")

		mockRepo.On("SearchCompanies", testifymock.Anything, keyword).Return(nil, expectedErr)

		uc := usecase.NewCompanyUsecase(mockRepo)

		res, err := uc.SearchCompanies(ctx, keyword)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, expectedErr, err)
		mockRepo.AssertExpectations(t)
	})
}
