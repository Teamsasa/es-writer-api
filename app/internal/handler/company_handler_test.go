package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"

	"es-api/app/internal/entity/model"
	"es-api/app/internal/handler"
)

type mockCompanyUsecase struct {
	testifymock.Mock
}

func (m *mockCompanyUsecase) SearchCompanies(ctx echo.Context, keyword string) ([]model.CompanyBasicInfo, error) {
	args := m.Called(ctx, keyword)
	return args.Get(0).([]model.CompanyBasicInfo), args.Error(1)
}

func TestCompanyHandler_SearchCompanies(t *testing.T) {
	mockUsecase := new(mockCompanyUsecase)
	h := handler.NewCompanyHandler(mockUsecase)

	t.Run("正常系:検索結果あり", func(t *testing.T) {
		expectedResponse := []model.CompanyBasicInfo{
			{
				CompanyID:   "1234567890123",
				CompanyName: "株式会社テスト",
			},
		}
		keyword := "株式会社テスト"

		mockUsecase.On("SearchCompanies", testifymock.Anything, keyword).Return(expectedResponse, nil)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/api/companies/search", nil)
		q := req.URL.Query()
		q.Add("keyword", keyword)
		req.URL.RawQuery = q.Encode()
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.SearchCompanies(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response []model.CompanyBasicInfo
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("正常系:検索結果なし", func(t *testing.T) {
		expectedResponse := []model.CompanyBasicInfo{}
		keyword := "存在しない会社"

		mockUsecase.On("SearchCompanies", testifymock.Anything, keyword).Return(expectedResponse, nil)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/api/companies/search", nil)
		q := req.URL.Query()
		q.Add("keyword", keyword)
		req.URL.RawQuery = q.Encode()
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.SearchCompanies(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response []model.CompanyBasicInfo
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("異常系:検索キーワード未指定", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/api/companies/search", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.SearchCompanies(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		mockUsecase.AssertExpectations(t)
	})
}
