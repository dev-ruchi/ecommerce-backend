package handlers

import (
	"fmt"
	"log"

	"e-store-backend/app"
	"e-store-backend/models"

	"github.com/gin-gonic/gin"
)

func HandleAddOrders(context *gin.Context) {
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
        INSERT INTO orders (user_id, product_id, quantity, total_price, status)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id`

	err = app.Db.QueryRow(query, order.UserId, order.ProductId, order.Quantity, order.TotalPrice, order.Status).Scan(
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

func HandleFetchOrders(context *gin.Context) {
	userId := uint(context.GetFloat64("userId"))
	rows, err := app.Db.Query(`SELECT id, user_id, product_id, quantity, total_price, status FROM orders WHERE user_id=$1`, userId)

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
		if err := rows.Scan(&order.Id, &order.UserId, &order.ProductId, &order.Quantity, &order.TotalPrice, &order.Status); err != nil {
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
