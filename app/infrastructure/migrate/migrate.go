package migrate

import (
	"log"

	"gorm.io/gorm"

	"tmp/app/internal/entity"
)

func RunMigrations(db *gorm.DB) {
	err := db.AutoMigrate(&entity.User{})
	if err != nil {
		log.Fatalf("🔴 Error migrating User model: %s", err)
	}
	log.Println("🟢 User model migrated")
}
