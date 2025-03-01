package main

import (
	"log"

	"github.com/joho/godotenv"

	"tmp/app/infrastructure/db"
	"tmp/app/infrastructure/migrate"
)

func main() {
	err := godotenv.Load("../../../.env")
	if err != nil {
		log.Fatalln(err)
	}
	dbConnection := db.NewDB()
	migrate.RunMigrations(dbConnection)
	db.CloseDB(dbConnection)
	log.Println("🟢 Migrations completed")
}
