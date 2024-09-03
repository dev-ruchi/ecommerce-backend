package handlers

import (
	"fmt"

	"e-store-backend/app"
	"e-store-backend/models"

	"github.com/gin-gonic/gin"
)

func HandleAddReviews(context *gin.Context) {
	var review models.Review

	err := context.BindJSON(&review)

	if err != nil {
		fmt.Println(err)
		context.JSON(400, gin.H{
			"message": "Bad request",
		})
		return
	}

	query := `
        INSERT INTO reviews (rating, comment)
        VALUES ($1, $2)
        RETURNING id`

	err = app.Db.QueryRow(query, review.Rating, review.Comment).Scan(
		&review.Id,
	)

	if err != nil {
		fmt.Println(err)
		context.JSON(500, gin.H{
			"message": "Something went wrong",
		})
		return
	}

	context.JSON(201, review)
}
