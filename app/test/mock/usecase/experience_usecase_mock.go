package mock

import (
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"

	"es-api/app/internal/entity/model"
)

type ExperienceUsecaseMock struct {
	mock.Mock
}

func (m *ExperienceUsecaseMock) GetExperienceByUserID(c echo.Context) (*model.Experiences, error) {
	args := m.Called(c)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*model.Experiences), args.Error(1)
}

func (m *ExperienceUsecaseMock) PostExperience(c echo.Context, inputExperience model.InputExperience) (*model.Experiences, error) {
	args := m.Called(c, inputExperience)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*model.Experiences), args.Error(1)
}
