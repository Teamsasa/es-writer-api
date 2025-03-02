package test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"es-api/app/infrastructure/db"
	"es-api/app/infrastructure/migrate"
)

func loadEnvFile(t *testing.T, path string) {
	err := godotenv.Load(path)
	if err != nil {
		log.Fatalln(err)
	}
}

func SetupTestDB(t *testing.T, path string) *gorm.DB {
	loadEnvFile(t, path)

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

func SetupEchoContext(userID string) echo.Context {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	req.Header.Set("idp", "test")
	ctx := e.NewContext(req, rec)
	ctx.Set("userID", userID)
	return ctx
}
