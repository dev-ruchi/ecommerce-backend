package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"e-store-backend/app"
	"e-store-backend/models"

	"github.com/gin-gonic/gin"
)

func HandleAddProducts(context *gin.Context) {
	var product models.Product

	err := context.BindJSON(&product)

	if err != nil {
		fmt.Println(err)
		context.JSON(400, gin.H{
			"message": "Bad request",
		})
		return
	}

	query := `
        INSERT INTO products (title, price, description, rating, images)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id`

	err = app.Db.QueryRow(query, product.Title, product.Price, product.Description, product.Rating, product.Images).Scan(
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

func HandleFetchProducts(context *gin.Context) {
	rows, err := app.Db.Query("SELECT * FROM products")

	if err != nil {

		log.Fatal(err)

		context.JSON(500, gin.H{
			"message": "Something went wrong",
		})
	}

	defer rows.Close()

	var products []models.Product

	for rows.Next() {

		var product models.Product

		if err := rows.Scan(&product.Id, &product.Title, &product.Price, &product.Description, &product.Rating, &product.Images); err != nil {

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

	if (products == nil) {
		context.JSON(200, []models.Product{})
		return
	}

	context.JSON(200, products)
}

func HandleUpdateProducts(context *gin.Context) {
	context.JSON(200, gin.H{
		"message": "Updated successfully",
	})
}

func HandleDeleteProducts(context *gin.Context) {

	query := `
      DELETE FROM products WHERE id=$1;`

	_, err := app.Db.Query(query, context.Param("id"))

	if err != nil {
		fmt.Println(err)
		context.JSON(500, gin.H{
			"message": "Something went wrong",
		})
		return
	}

	context.Status(204)
}

func HandleFetchProduct(context *gin.Context) {
	// Get the product ID from the URL parameters
	productID := context.Param("id")

	// Prepare the SQL query
	query := "SELECT id, title, price, description, rating, images FROM products WHERE id = $1"

	// Query the database for the product
	var product models.Product

	err := app.Db.QueryRow(query, productID).Scan(&product.Id, &product.Title, &product.Price, &product.Description, &product.Rating, &product.Images)
	if err != nil {
		if err == sql.ErrNoRows {
			// If no product is found, return a 404 Not Found response
			context.JSON(http.StatusNotFound, gin.H{
				"message": "Product not found",
			})
		} else {
			// Log the error and return a 500 Internal Server Error response
			log.Println("Error fetching product:", err)
			context.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to fetch product",
			})
		}
		return
	}

	// Return the product as JSON
	context.JSON(http.StatusOK, product)
}
