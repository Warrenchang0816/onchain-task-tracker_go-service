package handler

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "file is required",
		})
		return
	}

	// 確保 uploads 資料夾存在
	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to create uploads directory",
		})
		return
	}

	filename := filepath.Base(file.Filename)
	filePath := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to save file",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"url":     "http://localhost:8081/uploads/" + filename,
		"message": "upload success",
	})
}