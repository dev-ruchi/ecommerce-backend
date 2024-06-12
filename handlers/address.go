package handlers

import (
	"fmt"

	"e-store-backend/app"
	"e-store-backend/models"

	"github.com/gin-gonic/gin"
)

func HandleAddress(context *gin.Context) {
	var address models.Address

	err := context.BindJSON(&address)

	if err != nil {
		fmt.Println(err)
		context.JSON(400, gin.H{
			"message": "Bad request",
		})
		return
	}

	query := `
        INSERT INTO address (street, city, state, pin_code)
        VALUES ($1, $2, $3, $4)
        RETURNING id`

	err = app.Db.QueryRow(query, address.Street, address.City, address.State, address.PinCode).Scan(
		&address.Id,
	)

	if err != nil {
		fmt.Println(err)
		context.JSON(500, gin.H{
			"message": "Something went wrong",
		})
		return
	}

	context.JSON(201, address)
}
