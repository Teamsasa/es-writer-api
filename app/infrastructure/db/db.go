package db

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBConnectionManager interface {
	GetConnection(idp string) *gorm.DB
}

type dbConnectionManager struct {
	prodDB    *gorm.DB
	swaggerDB *gorm.DB
}

func NewDBConnectionManager() DBConnectionManager {
	prodDB := NewDB()
	swaggerDB := NewSwaggerDB()

	return &dbConnectionManager{
		prodDB:    prodDB,
		swaggerDB: swaggerDB,
	}
}

func (m *dbConnectionManager) GetConnection(idp string) *gorm.DB {
	switch idp {
	case "clerk":
		return m.prodDB
	default:
		return m.swaggerDB
	}
}

func NewDB() *gorm.DB {
	url := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  url,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		log.Fatalf("🔴 Error connecting to database: %s", err)
	}
	log.Println("🟢 Connected to database")
	return db
}

func CloseDB(db *gorm.DB) {
	sqlDB, _ := db.DB()
	if err := sqlDB.Close(); err != nil {
		log.Fatalf("🔴 Error closing to database: %s", err)
	}
	log.Println("🟢 Database connection closed")
}

func NewSwaggerDB() *gorm.DB {
	url := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("SWAGGER_DB_HOST"),
		os.Getenv("SWAGGER_DB_PORT"),
		os.Getenv("SWAGGER_DB_USER"),
		os.Getenv("SWAGGER_DB_PASSWORD"),
		os.Getenv("SWAGGER_DB_NAME"))

	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		log.Fatalf("🔴 Error connecting to swagger database: %s", err)
	}
	log.Println("🟢 Connected to swagger database")
	return db
}

func CleanupSwaggerDB(db *gorm.DB) {
	db.Exec("DELETE FROM users")
}

func NewTestDB() *gorm.DB {
	url := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("TEST_DB_HOST"),
		os.Getenv("TEST_DB_PORT"),
		os.Getenv("TEST_DB_USER"),
		os.Getenv("TEST_DB_PASSWORD"),
		os.Getenv("TEST_DB_NAME"))

	// テスト時はログを無効化する
	config := &gorm.Config{
		Logger: logger.New(
			log.New(io.Discard, "", log.LstdFlags),
			logger.Config{
				SlowThreshold:             200 * time.Millisecond,
				LogLevel:                  logger.Silent,
				IgnoreRecordNotFoundError: true,
				Colorful:                  false,
			},
		),
	}

	db, err := gorm.Open(postgres.Open(url), config)
	if err != nil {
		log.Fatalf("🔴 Error connecting to test database: %s", err)
	}
	log.Println("🟢 Connected to test database")
	return db
}

func CleanupTestDB(db *gorm.DB) {
	db.Exec("DELETE FROM users")
}
