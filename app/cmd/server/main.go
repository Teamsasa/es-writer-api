package main

import (
	"log"

	"github.com/joho/godotenv"

	"es-api/app/infrastructure/db"
	"es-api/app/internal/handler"
	"es-api/app/internal/repository"
	"es-api/app/internal/router"
	"es-api/app/internal/usecase"
)

func main() {
	err := godotenv.Load("../../../.env")
	if err != nil {
		log.Fatalln(err)
	}
	dbConnection := db.NewDB()
	welcomeRepository := repository.NewWelcomeRepository(dbConnection)
	welcomeUsecase := usecase.NewWelcomeUsecase(welcomeRepository)
	welcomeHandler := handler.NewWelcomeHandler(welcomeUsecase)
	e := router.NewRouter(welcomeHandler)
	e.Logger.Fatal(e.Start(":8080"))
}
