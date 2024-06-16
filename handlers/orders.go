package handlers

import (
	"fmt"
	"log"

	"e-store-backend/app"
	"e-store-backend/models"

	"github.com/gin-gonic/gin"
)

func HandleOrders(context *gin.Context) {
	var order models.Order

	err := context.BindJSON(&order)

	if err != nil {
		fmt.Println(err)
		context.JSON(400, gin.H{
			"message": "Bad request",
		})
		return
	}

	query := `
        INSERT INTO orders (user_id, product_id, quantity, total_price)
        VALUES ($1, $2, $3, $4)
        RETURNING id`

	err = app.Db.QueryRow(query, order.UserId, order.ProductId, order.Quantity, order.TotalPrice).Scan(
		&order.Id,
	)

	if err != nil {
		fmt.Println(err)
		context.JSON(500, gin.H{
			"message": "Something went wrong",
		})
		return
	}

	context.JSON(201, order)
}

func HandlePlaceOrders(context *gin.Context) {
	rows, err := app.Db.Query(`SELECT id, user_id, product_id, quantity, total_price FROM orders WHERE user_id=$1`, context.Param("user_id"))

	if err != nil {
		log.Fatal(err)
		context.JSON(500, gin.H{
			"message": "Something went wrong",
		})
		return
	}

	defer rows.Close()

	var orders []models.Order

	for rows.Next() {
		var order models.Order
		if err := rows.Scan(&order.Id, &order.UserId, &order.ProductId, &order.Quantity, &order.TotalPrice); err != nil {
			log.Fatal(err)
			context.JSON(500, gin.H{
				"message": "Something went wrong",
			})
			return
		}
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
		context.JSON(500, gin.H{
			"message": "Something went wrong",
		})
		return
	}

	context.JSON(200, orders)
}
