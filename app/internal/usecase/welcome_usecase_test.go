package usecase_test

import (
	"errors"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"

	"tmp/app/internal/entity"
	"tmp/app/internal/usecase"
	appmock "tmp/app/test/mock/repository"
)

func TestWelcomeUsecase_GetRandomUser(t *testing.T) {
	mockRepo := new(appmock.WelcomeRepositoryMock)

	users := []entity.User{
		{
			ID:        "test-id-1",
			Name:      "Test User 1",
			Email:     "test1@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        "test-id-2",
			Name:      "Test User 2",
			Email:     "test2@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	e := echo.New()
	ctx := e.NewContext(nil, nil)
	mockRepo.On("GetAllUsers", testifymock.Anything).Return(users, nil)

	uc := usecase.NewWelcomeUsecase(mockRepo)

	user, err := uc.GetRandomUser(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Contains(t, []string{"test-id-1", "test-id-2"}, user.ID)
	mockRepo.AssertExpectations(t)
}

func TestWelcomeUsecase_GetRandomUser_EmptyUsers(t *testing.T) {
	// モックの準備
	mockRepo := new(appmock.WelcomeRepositoryMock)

	// モックの振る舞いを設定
	e := echo.New()
	ctx := e.NewContext(nil, nil)
	mockRepo.On("GetAllUsers", testifymock.Anything).Return([]entity.User{}, nil)

	// テスト対象のユースケースを初期化
	uc := usecase.NewWelcomeUsecase(mockRepo)

	// テスト実行
	user, err := uc.GetRandomUser(ctx)

	// アサーション
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "users not found", err.Error())

	// モックが期待通り呼ばれたことを確認
	mockRepo.AssertExpectations(t)
}

func TestWelcomeUsecase_GetRandomUser_Error(t *testing.T) {
	// モックの準備
	mockRepo := new(appmock.WelcomeRepositoryMock)

	// モックの振る舞いを設定
	e := echo.New()
	ctx := e.NewContext(nil, nil)
	expectedErr := errors.New("database error")
	mockRepo.On("GetAllUsers", testifymock.Anything).Return([]entity.User{}, expectedErr)

	// テスト対象のユースケースを初期化
	uc := usecase.NewWelcomeUsecase(mockRepo)

	// テスト実行
	user, err := uc.GetRandomUser(ctx)

	// アサーション
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, expectedErr, err)

	// モックが期待通り呼ばれたことを確認
	mockRepo.AssertExpectations(t)
}
