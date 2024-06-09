package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Product struct {
	Id          int     `json:"id"`
	Title       string  `json:"title"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	Rating      float64 `json:"rating"`
}

var Db *sql.DB

func setupDatabase() {
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

func handleAddProducts(context *gin.Context) {
	var product Product

	err := context.BindJSON(&product)

	if err != nil {
		fmt.Println(err)
		context.JSON(400, gin.H{
			"message": "Bad request",
		})
		return
	}

	query := `
        INSERT INTO products (title, price, description, rating)
        VALUES ($1, $2, $3, $4)
        RETURNING id`

	err = Db.QueryRow(query, product.Title, product.Price, product.Description, product.Rating).Scan(
		&product.Id,
	)

	if err != nil {
		fmt.Println(err)
		context.JSON(500, gin.H{
			"message": "Something went wrong",
		})
		return
	}

	context.JSON(201, product)
}

func handleFetchProducts(context *gin.Context) {
	rows, err := Db.Query("SELECT * FROM products")

	if err != nil {

		log.Fatal(err)

		context.JSON(500, gin.H{
			"message": "Something went wrong",
		})
	}

	defer rows.Close()

	var products []Product

	for rows.Next() {

		var product Product

		if err := rows.Scan(&product.Id, &product.Title, &product.Price, &product.Description, &product.Rating); err != nil {

			log.Fatal(err)

			context.JSON(500, gin.H{
				"message": "Something went wrong",
			})
		}

		products = append(products, product)
	}

	if err = rows.Err(); err != nil {

		log.Fatal(err)

		context.JSON(500, gin.H{
			"message": "Something went wrong",
		})
	}

	context.JSON(200, products)
}

func handleUpdateProducts(context *gin.Context) {
	context.JSON(200, gin.H{
		"message": "Updated successfully",
	})
}

func handleDeleteProducts(context *gin.Context) {

	query := `
      DELETE FROM products WHERE id=$1;`

	_, err := Db.Query(query, context.Param("id"))

	if err != nil {
		fmt.Println(err)
		context.JSON(500, gin.H{
			"message": "Something went wrong",
		})
		return
	}

	context.Status(204)
}

func handleAddUsers(context *gin.Context) {
	var user User

	err := context.BindJSON(&user)

	if err != nil {
		fmt.Println(err)
		context.JSON(400, gin.H{
			"message": "Bad request",
		})
		return
	}

	query := `
        INSERT INTO users (username, email, password)
        VALUES ($1, $2, $3)
        RETURNING id`

	err = Db.QueryRow(query, user.Username, user.Email, user.Password).Scan(
		&user.Id,
	)

	if err != nil {
		fmt.Println(err)
		context.JSON(500, gin.H{
			"message": "Something went wrong",
		})
		return
	}

	context.JSON(201, user)
}

func handleLogin(context *gin.Context) {
	var loginRequest LoginRequest

	err := context.BindJSON(&loginRequest)

	if err != nil {
		fmt.Println(err)
		context.JSON(400, gin.H{
			"message": "Bad request",
		})
		return
	}

	var user User

	query := `
        SELECT id, username, email, password
        FROM users
        WHERE email = $1`

	// Query the database to find the user by username
	row := Db.QueryRow(query, loginRequest.Email)

	if row.Err() != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Invalid email or password",
		})
		return
	}

	row.Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.Password,
	)

	if loginRequest.Password != user.Password {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Invalid email and password combination",
		})
		return
	}

	// If the password is correct, return a success response
	context.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user": gin.H{
			"id":       user.Id,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	fmt.Println("Welcome to book app!")

	err := godotenv.Load()

	if err != nil {
		fmt.Println("Error loading .env file")
	}

	setupDatabase()

	if err != nil {
		fmt.Println(err)
	}

	router := gin.Default()

	router.Use(CORSMiddleware())

	router.POST("/login", handleLogin)
	router.POST("/users", handleAddUsers)
	router.GET("/products", handleFetchProducts)
	router.POST("/products", handleAddProducts)
	router.PUT("/products", handleUpdateProducts)
	router.DELETE("/products/:id", handleDeleteProducts)

	router.Run()

	defer Db.Close()
}
