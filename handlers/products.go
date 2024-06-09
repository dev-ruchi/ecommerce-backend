package handlers

import (
	"fmt"
	"log"

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
        INSERT INTO products (title, price, description, rating)
        VALUES ($1, $2, $3, $4)
        RETURNING id`

	err = app.Db.QueryRow(query, product.Title, product.Price, product.Description, product.Rating).Scan(
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
