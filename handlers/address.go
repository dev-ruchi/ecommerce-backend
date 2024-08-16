package handlers

import (
	"fmt"
	"log"

	"e-store-backend/app"
	"e-store-backend/models"

	"github.com/gin-gonic/gin"
)

func HandleAddAddress(context *gin.Context) {
	var address models.Address

	address.UserId =  int(context.GetFloat64("userId"))

	err := context.BindJSON(&address)

	if err != nil {
		fmt.Println(err)
		context.JSON(400, gin.H{
			"message": "Bad request",
		})
		return
	}

	query := `
        INSERT INTO addresses (user_id, street, city, state, pin_code)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id`

	err = app.Db.QueryRow(query, address.UserId, address.Street, address.City, address.State, address.PinCode).Scan(
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

func HandleFetchAddresses(context *gin.Context) {
	userId := uint(context.GetFloat64("userId"))
	rows, err := app.Db.Query("SELECT id, user_id, street, city, state, pin_code FROM addresses WHERE user_id=$1", userId)

	if err != nil {

		log.Fatal(err)

		context.JSON(500, gin.H{
			"message": "Something went wrong",
		})
	}

	defer rows.Close()

	var addresses []models.Address

	for rows.Next() {

		var address models.Address

		if err := rows.Scan(&address.Id, &address.UserId, &address.Street, &address.City, &address.State, &address.PinCode); err != nil {

			log.Fatal(err)

			context.JSON(500, gin.H{
				"message": "Something went wrong",
			})
		}

		addresses = append(addresses, address)
	}

	if err = rows.Err(); err != nil {

		log.Fatal(err)

		context.JSON(500, gin.H{
			"message": "Something went wrong",
		})
	}

	if addresses == nil {
		context.JSON(200, []models.Address{})
		return
	}

	context.JSON(200, addresses)
}
