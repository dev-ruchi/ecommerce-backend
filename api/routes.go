package api

import (
	"e-store-backend/handlers"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	jwt "github.com/golang-jwt/jwt/v5"
)

func SetupRoutes() {
	router := gin.Default()

	router.MaxMultipartMemory = 8 << 20 // 8 MiB

	router.Use(corsMiddleware())

	router.POST("/login", handlers.HandleLogin)
	router.POST("/users", handlers.HandleAddUsers)

	router.GET("/products", handlers.HandleFetchProducts)
	router.GET("/products/:id", handlers.HandleFetchProduct)
	router.POST("/products", handlers.HandleAddProducts)
	router.PUT("/products", handlers.HandleUpdateProducts)
	router.DELETE("/products/:id", handlers.HandleDeleteProducts)

	guardedRoutes := router.Group("/", authMiddleware)

	guardedRoutes.POST("/orders", handlers.HandleAddOrders)
	guardedRoutes.GET("/orders", handlers.HandleFetchOrders)

	guardedRoutes.POST("/address", handlers.HandleAddAddress)
	guardedRoutes.GET("/addresses", handlers.HandleFetchAddresses)

	guardedRoutes.POST("/files/upload", handlers.HandleFilesUpload)
	router.GET("/files/:fileName", handlers.HandleFetchFile)

	router.GET("/payment/:order_id", handlers.HandlePayment)

	router.POST("/reviews", handlers.HandleAddReviews)

	router.Run()
}

func authMiddleware(c *gin.Context) {
	// Get the Authorization header value
	tokenString := c.GetHeader("Authorization")

	// Check if the token is missing
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		fmt.Println("Token not found")
		return
	}

	secretKey := os.Getenv("JWT_SECRET")

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	// Check for parsing errors
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		fmt.Println("Unable to parse token string")
		return
	}

	// Check if the token is valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Extract userId from the token claims
		userId, ok := claims["userId"]
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			fmt.Println("Invalid token")
			return
		}

		// Add userId to the Gin context for further use
		c.Set("userId", userId)

		// Continue with the next middleware or handler
		c.Next()
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
	}
}

func corsMiddleware() gin.HandlerFunc {
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
