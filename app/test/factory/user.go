package factory

import (
	"testing"
	"time"

	"gorm.io/gorm"

	"es-api/app/internal/entity/model"
)

const (
	DummyUserID1 = "user_abcdefghijklmnopqrstuvwxyz"
	DummyUserID2 = "user_ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func CreateUser1(t *testing.T, dbConn *gorm.DB) model.Users {
	user := model.Users{
		ID:        "user_abcdefghijklmnopqrstuvwxyz",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := dbConn.Create(&user).Error; err != nil {
		t.Fatalf("Error creating test user: %v", err)
	}

	return user
}

func CreateUser2(t *testing.T, dbConn *gorm.DB) model.Users {
	user := model.Users{
		ID:        "user_ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := dbConn.Create(&user).Error; err != nil {
		t.Fatalf("Error creating test user: %v", err)
	}

	return user
}
