package main

import (
	"log"

	"github.com/joho/godotenv"

	"es-api/app/infrastructure/db"
	"es-api/app/infrastructure/migrate"
)

func main() {
	err := godotenv.Load("../../../.env")
	if err != nil {
		log.Fatalln(err)
	}
	dbConnection := db.NewDB()
	swaggerDbConnection := db.NewSwaggerDB()
	migrate.RunMigrations(dbConnection)
	migrate.RunMigrations(swaggerDbConnection)
	db.CloseDB(dbConnection)
	db.CloseDB(swaggerDbConnection)
	log.Println("🟢 Migrations completed")
}
