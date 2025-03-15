package mock

import (
	"context"

	"es-api/app/internal/entity/model"

	"github.com/stretchr/testify/mock"
)

type ExperienceUsecaseMock struct {
	mock.Mock
}

func (m *ExperienceUsecaseMock) GetExperienceByUserID(ctx context.Context) (*model.Experiences, error) {
	args := m.Called(ctx)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*model.Experiences), args.Error(1)
}

func (m *ExperienceUsecaseMock) PostExperience(ctx context.Context, inputExperience model.InputExperience) (*model.Experiences, error) {
	args := m.Called(ctx, inputExperience)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*model.Experiences), args.Error(1)
}
