package test

import (
	"context"
	"log"
	"testing"

	"github.com/joho/godotenv"
	"gorm.io/gorm"

	"es-api/app/infrastructure/db"
	"es-api/app/infrastructure/migrate"
	"es-api/app/internal/contextKey"
)

func LoadEnvFile(t *testing.T, path string) {
	err := godotenv.Load(path)
	if err != nil {
		log.Fatalln(err)
	}
}

func SetupTestDB(t *testing.T, path string) *gorm.DB {
	LoadEnvFile(t, path)

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

func SetupContextContext(userID string) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, contextKey.UserIDKey, userID)
	ctx = context.WithValue(ctx, contextKey.IDPKey, "test")
	return ctx
}
