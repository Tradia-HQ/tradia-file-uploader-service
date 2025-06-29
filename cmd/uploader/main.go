package main

import (
	"fmt"
	"github.com/Tradia-HQ/tradia-file-uploader-service/internalpackage/config"
	"github.com/Tradia-HQ/tradia-file-uploader-service/internalpackage/services"
	"github.com/Tradia-HQ/tradia-file-uploader-service/internalpackage/services/interfaces"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

const (
	// MaxFileSize is 10MB in bytes (10 * 1024 * 1024)
	MaxFileSize = 10 * 1024 * 1024
)

type App struct {
	uploader interfaces.FileUploader
}

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize uploader
	uploader, err := services.NewUploader(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize uploader: %v", err)
	}

	// Create application
	app := &App{uploader: uploader}

	// Initialize Gin router
	r := gin.Default()

	// Define routes
	r.POST("/upload", app.handleUpload)
	r.GET("/signed-url", app.handleSignedURL)
	r.DELETE("/delete", app.handleDelete)

	// Start server
	log.Printf("Starting server on port %d", cfg.Port)
	if err := r.Run(fmt.Sprintf(":%d", cfg.Port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func (app *App) handleUpload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read file: " + err.Error()})
		return
	}
	defer file.Close()

	// Check file size
	if header.Size > MaxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("File size exceeds 10MB limit (got %d bytes)", header.Size)})
		return
	}

	publicURL, objectName, err := app.uploader.Upload(c.Request.Context(), file, header.Filename, header.Header.Get("Content-Type"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "File uploaded successfully",
		"publicURL":  publicURL,
		"objectName": objectName,
	})
}

func (app *App) handleSignedURL(c *gin.Context) {
	filename := c.Query("filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Filename required"})
		return
	}

	signedURL, err := app.uploader.GetSignedURL(c.Request.Context(), filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get signed URL: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Signed URL generated successfully",
		"signedURL": signedURL,
	})
}

func (app *App) handleDelete(c *gin.Context) {
	fileURL := c.Query("fileURL")
	if fileURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File URL required"})
		return
	}

	if err := app.uploader.Delete(c.Request.Context(), fileURL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "File deleted successfully",
	})
}
