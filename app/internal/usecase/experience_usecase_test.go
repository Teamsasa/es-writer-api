package usecase_test

import (
	"errors"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"

	"es-api/app/internal/entity/model"
	"es-api/app/internal/usecase"
	appmock "es-api/app/test/mock/repository"
)

// (getExperienceByUserID - 正常系 - データが存在する場合)
func TestExperienceUsecase_GetExperienceByUserID(t *testing.T) {
	mockRepo := new(appmock.ExperienceRepositoryMock)

	experience := model.Experiences{
		ID:        "test-id-1",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	e := echo.New()
	ctx := e.NewContext(nil, nil)
	mockRepo.On("GetExperienceByUserID", testifymock.Anything).Return(experience, nil)

	uc := usecase.NewExperienceUsecase(mockRepo)

	res, err := uc.GetExperienceByUserID(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, experience.ID, res.ID)
	mockRepo.AssertExpectations(t)
}

// (getExperienceByUserID - 異常系 - データが存在しない場合)
func TestExperienceUsecase_GetExperienceByUserID_Error(t *testing.T) {
	mockRepo := new(appmock.ExperienceRepositoryMock)

	e := echo.New()
	ctx := e.NewContext(nil, nil)
	expectedErr := errors.New("repository error")
	mockRepo.On("GetExperienceByUserID", testifymock.Anything).Return(model.Experiences{}, expectedErr)

	uc := usecase.NewExperienceUsecase(mockRepo)

	res, err := uc.GetExperienceByUserID(ctx)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, expectedErr, err)
	mockRepo.AssertExpectations(t)
}

// (postExperience - 正常系 - 経験が存在しない場合)
func TestExperienceUsecase_PostExperience_Success(t *testing.T) {
	mockRepo := new(appmock.ExperienceRepositoryMock)

	inputExperience := model.InputExperience{
		Work:        "test work",
		Skills:      "test skills",
		SelfPR:      "test self PR",
		FutureGoals: "test future goals",
	}

	createdExperience := model.Experiences{
		ID:          "test-id-1",
		UserID:      "test-user-id",
		Work:        inputExperience.Work,
		Skills:      inputExperience.Skills,
		SelfPR:      inputExperience.SelfPR,
		FutureGoals: inputExperience.FutureGoals,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	e := echo.New()
	ctx := e.NewContext(nil, nil)
	mockRepo.On("FindExperienceByUserID", testifymock.Anything).Return(false, nil)
	mockRepo.On("PostExperience", testifymock.Anything, inputExperience).Return(createdExperience, nil)

	uc := usecase.NewExperienceUsecase(mockRepo)

	res, err := uc.PostExperience(ctx, inputExperience)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, createdExperience.ID, res.ID)
	assert.Equal(t, createdExperience.Work, res.Work)
	assert.Equal(t, createdExperience.Skills, res.Skills)
	assert.Equal(t, createdExperience.SelfPR, res.SelfPR)
	assert.Equal(t, createdExperience.FutureGoals, res.FutureGoals)
	mockRepo.AssertExpectations(t)
}

// (postExperience - 異常系 - 既にデータが存在する場合)
func TestExperienceUsecase_PostExperience_AlreadyExists(t *testing.T) {
	mockRepo := new(appmock.ExperienceRepositoryMock)

	inputExperience := model.InputExperience{
		Work:        "test work",
		Skills:      "test skills",
		SelfPR:      "test self PR",
		FutureGoals: "test future goals",
	}

	e := echo.New()
	ctx := e.NewContext(nil, nil)
	mockRepo.On("FindExperienceByUserID", testifymock.Anything).Return(true, nil)

	uc := usecase.NewExperienceUsecase(mockRepo)

	res, err := uc.PostExperience(ctx, inputExperience)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "experience already exists", err.Error())
	mockRepo.AssertExpectations(t)
}
