package handlers

import (
	"fmt"

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
