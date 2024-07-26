package handlers

import (
	"crypto/sha512"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"e-store-backend/app"
	"e-store-backend/models"

	"github.com/gin-gonic/gin"
	payu "github.com/payu-india/web-sdk-go"
)

// calculateHash computes the SHA-512 hash based on the given input parameters.
func calculateHash(key, txnid, amount, productId, firstname, email, salt string) string {
	// Concatenate the input string as per the formula
	udf1 := ""
	udf2 := ""
	udf3 := ""
	udf4 := ""
	udf5 := ""
	input := key + "|" + txnid + "|" + amount + "|" + productId + "|" + firstname + "|" + email + "|" + udf1 + "|" + udf2 + "|" + udf3 + "|" + udf4 + "|" + udf5 + "||||||" + salt

	// Calculate the SHA-512 hash
	hash := sha512.Sum512([]byte(input))

	// Convert the hash to a hexadecimal string
	hashString := hex.EncodeToString(hash[:])
	return hashString
}

// HandleInitPayment handles the payment initialization and generates the payment form.
func HandlePayment(c *gin.Context) {

	payuClient, err := payu.NewClient(
		os.Getenv("PAYU_KEY"),
		os.Getenv("PAYU_SALT"),
		os.Getenv("PAYU_MODE"),
	)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong.",
		})
		return
	}

	var order models.Order

	query := `SELECT user_id, product_id, txn_id FROM orders WHERE id=$1`

	row := app.Db.QueryRow(query, c.Param("order_id"))

	row.Scan(&order.UserId, &order.ProductId, &order.TxnId)

	var user models.User

	query = `SELECT first_name, phone, email from users WHERE id = $1`

	if row.Err() != nil {
		c.JSON(400, gin.H{
			"message": "Bad request",
		})
		return
	}

	row = app.Db.QueryRow(query, order.UserId)

	if row.Err() != nil {
		c.JSON(401, gin.H{
			"message": "Unauthorized",
		})
		return
	}

	row.Scan(
		&user.FirstName,
		&user.Phone,
		&user.Email,
	)

	var product models.Product

	query = "SELECT id, price FROM products WHERE id = $1"

	err = app.Db.QueryRow(query, order.ProductId).Scan(&product.Id, &product.Price)

	if err != nil {
		if err == sql.ErrNoRows {
			// If no product is found, return a 404 Not Found response
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Product not found",
			})
		}

		return
	}

	// Define the input values
	txnid := strconv.Itoa(order.TxnId)
	amount := fmt.Sprintf("%.2f", product.Price)
	productId := strconv.Itoa(product.Id)

	// Calculate the hash
	hash := calculateHash(
		os.Getenv("PAYU_KEY"),
		txnid,
		amount,
		productId,
		user.FirstName,
		user.Email,
		os.Getenv("PAYU_SALT"),
	)

	// Generate the payment form with the calculated hash
	form, err := payuClient.GeneratePaymentForm(map[string]interface{}{
		"txnid":       txnid,
		"amount":      amount,
		"productinfo": productId,
		"firstname":   user.FirstName,
		"email":       user.Email,
		"phone":       user.Phone,
		"surl":        "http://localhost:5173/orders?orderId=" + c.Param("order_id"),
		"furl":        "http://localhost:5173/500",
		"hash":        hash,
	})

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong.",
		})
		return
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(form))
}
