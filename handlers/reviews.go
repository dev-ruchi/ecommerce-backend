package handlers

import (
	"fmt"
	"log"

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

func HandleFetchReviews(context *gin.Context) {
	rows, err := app.Db.Query("SELECT * FROM reviews")

	if err != nil {

		log.Fatal(err)

		context.JSON(500, gin.H{
			"message": "Something went wrong",
		})
	}

	defer rows.Close()

	var reviews []models.Review

	for rows.Next() {

		var review models.Review

		if err := rows.Scan(&review.Id, &review.Rating, &review.Comment); err != nil {

			log.Fatal(err)

			context.JSON(500, gin.H{
				"message": "Something went wrong",
			})
		}

		reviews = append(reviews, review)
	}

	if err = rows.Err(); err != nil {

		log.Fatal(err)

		context.JSON(500, gin.H{
			"message": "Something went wrong",
		})
	}

	if reviews == nil {
		context.JSON(200, []models.Review{})
		return
	}

	context.JSON(200, reviews)
}
