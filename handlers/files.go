package handlers

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type FileUpload struct {
	Files []*multipart.FileHeader `form:"files[]" binding:"required"`
}

func HandleFilesUpload(c *gin.Context) {
	var form FileUpload

	// Parse the form data including the files
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var fileNames []string // Array to store filenames to be returned in response

	// Iterate through each file
	for _, fileHeader := range form.Files {
		// Access the original file name
		originalFileName := fileHeader.Filename
		prefix := strings.ReplaceAll(time.Now().UTC().Format(time.RFC3339), ":", "") // Create time string and remove :
		fileName := prefix + strings.ReplaceAll(originalFileName, " ", "") // Remove spaces from the filename
		fileNames = append(fileNames, fileName)
		err := c.SaveUploadedFile(fileHeader, "uploads/"+fileName) // Save the file
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}

		// Process the uploaded file as needed
		fmt.Printf("Original File Name: %s, File saved: %s\n", originalFileName, originalFileName)
	}

	// Send a response
	c.JSON(http.StatusCreated, gin.H{"files": fileNames})
}

func HandleFetchFile(c *gin.Context) {
	fileName := c.Param("fileName")
	c.File("uploads/" + fileName)
}

