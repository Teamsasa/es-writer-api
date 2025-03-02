package repository

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        string    `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type DBAuthRepository interface {
	FindUser(userID string) (bool, error)
	CreateUser(userID string) error
}

type dbAuthRepository struct {
	db *gorm.DB
}

func NewDBAuthRepository(db *gorm.DB) DBAuthRepository {
	return &dbAuthRepository{
		db: db,
	}
}

func (r *dbAuthRepository) FindUser(userID string) (bool, error) {
	var count int64
	if err := r.db.Model(&User{}).Where("id = ?", userID).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}

	return count > 0, nil
}

func (r *dbAuthRepository) CreateUser(userID string) error {
	user := User{
		ID: userID,
	}

	if err := r.db.Create(&user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}
