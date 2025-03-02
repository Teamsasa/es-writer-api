package mock

import (
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"

	"es-api/app/internal/entity/model"
)

type ExperienceRepositoryMock struct {
	mock.Mock
}

func (m *ExperienceRepositoryMock) GetExperienceByUserID(c echo.Context) (model.Experiences, error) {
	args := m.Called(c)
	return args.Get(0).(model.Experiences), args.Error(1)
}

func (m *ExperienceRepositoryMock) FindExperienceByUserID(c echo.Context) (bool, error) {
	args := m.Called(c)
	return args.Bool(0), args.Error(1)
}

func (m *ExperienceRepositoryMock) PostExperience(c echo.Context, experience model.InputExperience) (model.Experiences, error) {
	args := m.Called(c, experience)
	return args.Get(0).(model.Experiences), args.Error(1)
}
