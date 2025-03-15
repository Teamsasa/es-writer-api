package mock

import (
	"context"

	"github.com/stretchr/testify/mock"

	"es-api/app/internal/entity/model"
)

type ExperienceRepositoryMock struct {
	mock.Mock
}

func (m *ExperienceRepositoryMock) GetExperienceByUserID(ctx context.Context) (model.Experiences, error) {
	args := m.Called(ctx)
	return args.Get(0).(model.Experiences), args.Error(1)
}

func (m *ExperienceRepositoryMock) FindExperienceByUserID(ctx context.Context) (bool, error) {
	args := m.Called(ctx)
	return args.Bool(0), args.Error(1)
}

func (m *ExperienceRepositoryMock) PostExperience(ctx context.Context, experience model.InputExperience) (model.Experiences, error) {
	args := m.Called(ctx, experience)
	return args.Get(0).(model.Experiences), args.Error(1)
}

func (m *ExperienceRepositoryMock) PatchExperience(ctx context.Context, experience model.InputExperience) (model.Experiences, error) {
	args := m.Called(ctx, experience)
	return args.Get(0).(model.Experiences), args.Error(1)
}
