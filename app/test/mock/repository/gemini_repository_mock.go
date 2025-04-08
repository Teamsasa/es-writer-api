package mock

import (
	"context"

	"es-api/app/internal/entity/model"

	"github.com/stretchr/testify/mock"
)

type GeminiRepositoryMock struct {
	mock.Mock
}

func (m *GeminiRepositoryMock) GetGeminiRequest(ctx context.Context, input model.GeminiInput) (model.GeminiResponse, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(model.GeminiResponse), args.Error(1)
}
