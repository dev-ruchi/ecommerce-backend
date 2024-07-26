package handlers

import (
	"fmt"
	"net/http"

	"e-store-backend/app"
	"e-store-backend/models"
	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func HandleLogin(context *gin.Context) {
	var loginRequest LoginRequest

	err := context.BindJSON(&loginRequest)

	if err != nil {
		fmt.Println(err)
		context.JSON(400, gin.H{
			"message": "Bad request",
		})
		return
	}

	var user models.User

	query := `
        SELECT id, first_name, last_name, phone, email, password
        FROM users
        WHERE email = $1`

	// Query the database to find the user by username
	row := app.Db.QueryRow(query, loginRequest.Email)

	if row.Err() != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Invalid email or password",
		})
		return
	}

	row.Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
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
			"id":         user.Id,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"email":      user.Email,
		},
	})
}

func HandleAddUsers(context *gin.Context) {
	var user models.User

	err := context.BindJSON(&user)

	if err != nil {
		fmt.Println(err)
		context.JSON(400, gin.H{
			"message": "Bad request",
		})
		return
	}

	query := `
        INSERT INTO users (first_name, last_name, phone, email, password)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id`

	err = app.Db.QueryRow(query, user.FirstName, user.LastName, user.Phone, user.Email, user.Password).Scan(
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
