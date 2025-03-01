package handler_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"

	"tmp/app/internal/entity"
	"tmp/app/internal/handler"
	appmock "tmp/app/test/mock/usecase"
)

func TestWelcomeHandler_GetRandomUser(t *testing.T) {
	mockUsecase := new(appmock.WelcomeUsecaseMock)

	user := &entity.User{
		ID:        "test-id",
		Name:      "Test User",
		Email:     "test@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockUsecase.On("GetRandomUser", testifymock.Anything).Return(user, nil)

	h := handler.NewWelcomeHandler(mockUsecase)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/random-user", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.GetRandomUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	var responseUser entity.User
	err = json.Unmarshal(rec.Body.Bytes(), &responseUser)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, responseUser.ID)
	assert.Equal(t, user.Name, responseUser.Name)
	assert.Equal(t, user.Email, responseUser.Email)
	mockUsecase.AssertExpectations(t)
}

func TestWelcomeHandler_GetRandomUser_Error(t *testing.T) {
	mockUsecase := new(appmock.WelcomeUsecaseMock)

	expectedErr := errors.New("usecase error")
	mockUsecase.On("GetRandomUser", testifymock.Anything).Return(nil, expectedErr)

	h := handler.NewWelcomeHandler(mockUsecase)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/random-user", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.GetRandomUser(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	var response map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedErr.Error(), response["error"])
	mockUsecase.AssertExpectations(t)
}
