package app

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var Db *sql.DB

func SetupDatabase() {
	connectionString := fmt.Sprintf(
		"user=%s dbname=%s password=%s host=%s port=%s sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
	)

	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		log.Fatal(err)
	}

	Db = db

	createUsersTable()
	createProductTable()
}

func createUsersTable() {
	_, err := Db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username TEXT,
			email TEXT,
			password TEXT
		)
	`)

	if err != nil {
		fmt.Println(err)
	}
}

func createProductTable() {
	fmt.Println("SHOULD CREATE PRODUCT TABLE")
	_, err := Db.Exec(`
        CREATE TABLE IF NOT EXISTS products (
            id SERIAL PRIMARY KEY,
            title TEXT,
            price DECIMAL(10, 2),
            description TEXT,
            rating DECIMAL(3, 2)	
		)
    `)

	if err != nil {
		fmt.Println(err)
	}
}
