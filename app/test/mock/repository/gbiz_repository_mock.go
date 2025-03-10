package mock

import (
	"es-api/app/internal/entity/model"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
)

type GBizInfoRepositoryMock struct {
	mock.Mock
}

func (m *GBizInfoRepositoryMock) SearchCompanies(c echo.Context, keyword string) ([]model.CompanyBasicInfo, error) {
	args := m.Called(c, keyword)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.CompanyBasicInfo), args.Error(1)
}
