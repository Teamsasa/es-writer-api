package handler_test

import (
	"bytes"
	"encoding/json"
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
	h := handler.NewExperienceHandler(mockUsecase)

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

	t.Run("正常系:経験を取得できる", func(t *testing.T) {
		mockUsecase.On("GetExperienceByUserID", testifymock.Anything).Return(experience, nil)

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
	})
}

func TestExperienceHandler_PostExperience(t *testing.T) {
	mockUsecase := new(appmock.ExperienceUsecaseMock)
	h := handler.NewExperienceHandler(mockUsecase)

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

	inputExperience := model.InputExperience{
		Work:        "test-work",
		Skills:      "test-skills",
		SelfPR:      "test-self-pr",
		FutureGoals: "test-future-goals",
	}

	t.Run("正常系:経験を作成できる", func(t *testing.T) {
		mockUsecase.On("PostExperience", testifymock.Anything, inputExperience).Return(experience, nil)

		e := echo.New()
		reqBody, _ := json.Marshal(inputExperience)
		req := httptest.NewRequest(http.MethodPost, "/api/experience", bytes.NewBuffer(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
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
	})
}
