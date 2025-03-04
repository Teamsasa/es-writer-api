package mock

import (
	"es-api/app/internal/entity/model"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
)

type GeminiRepositoryMock struct {
	mock.Mock
}

func (m *GeminiRepositoryMock) GetGeminiRequest(c echo.Context, input model.GeminiInput) (model.GeminiResponse, error) {
	args := m.Called(c, input)
	return args.Get(0).(model.GeminiResponse), args.Error(1)
}
