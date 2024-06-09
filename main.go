package main

import (
	"e-store-backend/api"
	"e-store-backend/app"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		panic("Error loading .env file")
	}

	app.SetupDatabase()

	api.SetupRoutes()

	defer app.Db.Close()
}
