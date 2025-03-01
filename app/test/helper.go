package test

import (
	"log"
	"testing"
	"time"

	"gorm.io/gorm"

	"es-api/app/infrastructure/db"
	"es-api/app/infrastructure/migrate"
	"es-api/app/internal/entity/model"

	"github.com/joho/godotenv"
)

func loadEnvFile(t *testing.T) {
	err := godotenv.Load("../../../.env")
	if err != nil {
		log.Fatalln(err)
	}
}

func SetupTestDB(t *testing.T) *gorm.DB {
	loadEnvFile(t)

	dbConn := db.NewTestDB()
	migrate.RunMigrations(dbConn)
	db.CleanupTestDB(dbConn)
	return dbConn
}

func CleanupDB(t *testing.T, dbConn *gorm.DB) {
	db.CleanupTestDB(dbConn)
	sqlDB, err := dbConn.DB()
	if err != nil {
		t.Fatalf("Error getting DB instance: %v", err)
	}
	sqlDB.Close()
}

func SeedTestUser(t *testing.T, dbConn *gorm.DB) entity.User {
	user := entity.User{
		ID:        "test-id",
		Name:      "Test User",
		Email:     "test@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := dbConn.Create(&user).Error; err != nil {
		t.Fatalf("Error creating test user: %v", err)
	}

	return user
}
