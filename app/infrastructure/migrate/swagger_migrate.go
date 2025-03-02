package migrate

import (
	"log"

	"es-api/app/internal/entity/model"

	"gorm.io/gorm"
)

func RunSwaggerMigrations(db *gorm.DB) {
	err := db.AutoMigrate(&model.Users{})
	if err != nil {
		log.Fatalf("🔴 Error migrating User model: %s", err)
	}
	err = db.AutoMigrate(&model.Experiences{})
	if err != nil {
		log.Fatalf("🔴 Error migrating Experience model: %s", err)
	}
	log.Println("🟢 User and Experience models migrated")
}
