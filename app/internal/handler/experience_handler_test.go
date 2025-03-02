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

	"es-api/app/internal/entity/model"
	"es-api/app/internal/handler"
	appmock "es-api/app/test/mock/usecase"
)

func TestExperienceHandler_GetExperienceByUserID(t *testing.T) {
	mockUsecase := new(appmock.ExperienceUsecaseMock)

	experience := &model.Experiences{
		ID:          "test-id",
		UserID:      "test-user-id",
		Work:        "test-work",
		Skills:      "test-skills",
		SelfPR:      "test-self-pr",
		FutureGoals: "test-future-goals",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockUsecase.On("GetExperienceByUserID", testifymock.Anything).Return(experience, nil)

	h := handler.NewExperienceHandler(mockUsecase)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/experience", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.GetExperienceByUserID(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	var responseExperience model.Experiences
	err = json.Unmarshal(rec.Body.Bytes(), &responseExperience)
	assert.NoError(t, err)
	assert.Equal(t, experience.ID, responseExperience.ID)
	mockUsecase.AssertExpectations(t)
}

func TestExperienceHandler_PostExperience(t *testing.T) {
	mockUsecase := new(appmock.ExperienceUsecaseMock)

	experience := &model.Experiences{
		ID:          "test-id",
		UserID:      "test-user-id",
		Work:        "test-work",
		Skills:      "test-skills",
		SelfPR:      "test-self-pr",
		FutureGoals: "test-future-goals",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockUsecase.On("PostExperience", testifymock.Anything, testifymock.Anything).Return(experience, nil)

	h := handler.NewExperienceHandler(mockUsecase)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/experience", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.PostExperience(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	var responseExperience model.Experiences
	err = json.Unmarshal(rec.Body.Bytes(), &responseExperience)
	assert.NoError(t, err)
	assert.Equal(t, experience.ID, responseExperience.ID)
	mockUsecase.AssertExpectations(t)
}

func TestExperienceHandler_PostExperience_Error(t *testing.T) {
	mockUsecase := new(appmock.ExperienceUsecaseMock)

	mockUsecase.On("PostExperience", testifymock.Anything, testifymock.Anything).Return(nil, errors.New("error"))

	h := handler.NewExperienceHandler(mockUsecase)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/experience", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.PostExperience(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockUsecase.AssertExpectations(t)
}
