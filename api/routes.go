package api

import (
	"e-store-backend/handlers"

	"github.com/gin-gonic/gin"
)

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

func SetupRoutes() {
	router := gin.Default()

	router.Use(corsMiddleware())

	router.POST("/login", handlers.HandleLogin)
	router.POST("/users", handlers.HandleAddUsers)
	router.GET("/products", handlers.HandleFetchProducts)
	router.POST("/products", handlers.HandleAddProducts)
	router.PUT("/products", handlers.HandleUpdateProducts)
	router.DELETE("/products/:id", handlers.HandleDeleteProducts)

	router.Run()
}
