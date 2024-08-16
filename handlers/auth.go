package handlers

import (
	"fmt"
	"net/http"
	"os"

	"e-store-backend/app"
	"e-store-backend/models"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
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

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))

	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	tokenStr, err := generateJWT(user.Id)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		fmt.Println("Failed to generate the token")
		return
	}

	// If the password is correct, return a success response
	context.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   tokenStr,
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

	password, err := (bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost))

	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
		})
	}

	user.Password = string(password)

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

	tokenStr, err := generateJWT(user.Id)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		fmt.Println("Failed to generate the token")
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": "Signup successful",
		"token":   tokenStr,
		"user": gin.H{
			"id":         user.Id,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"email":      user.Email,
		},
	})
}

func generateJWT(userId int) (string, error) {
	claims := &jwt.MapClaims{
		"expiresAt": 15000,
		"userId":    userId,
	}

	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString([]byte(secret))
	return tokenStr, err
}
