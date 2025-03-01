package migrate

import (
	"log"

	"es-api/app/internal/entity/model"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) {
	err := db.AutoMigrate(&entity.User{})
	if err != nil {
		log.Fatalf("🔴 Error migrating User model: %s", err)
	}
	log.Println("🟢 User model migrated")
}
