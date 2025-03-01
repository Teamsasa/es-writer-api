package mock

import (
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"

	"es-api/app/internal/entity/model"
)

type WelcomeRepositoryMock struct {
	mock.Mock
}

func (m *WelcomeRepositoryMock) GetAllUsers(c echo.Context) ([]entity.User, error) {
	args := m.Called(c)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]entity.User), args.Error(1)
}
