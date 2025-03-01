package repository

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"es-api/app/internal/entity/model"
)

type WelcomeRepository interface {
	GetAllUsers(c echo.Context) ([]entity.User, error)
}

type welcomeRepository struct {
	db *gorm.DB
}

func NewWelcomeRepository(db *gorm.DB) WelcomeRepository {
	return &welcomeRepository{db: db}
}

func (r *welcomeRepository) GetAllUsers(c echo.Context) ([]entity.User, error) {
	var users []entity.User
	result := r.db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}
