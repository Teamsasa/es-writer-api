package migrate

import (
	"log"

	"es-api/app/internal/entity/model"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) {
	err := db.AutoMigrate(&model.Users{})
	if err != nil {
		log.Fatalf("🔴 Error migrating User model: %s", err)
	}
	err = db.AutoMigrate(&model.Experiences{})
	if err != nil {
		log.Fatalf("🔴 Error migrating Experience model: %s", err)
	}
	err = db.AutoMigrate(&model.CompanyResearch{})
	if err != nil {
		log.Fatalf("🔴 Error migrating CompanyResearch model: %s", err)
	}
	log.Println("🟢 Migrations completed")
}
