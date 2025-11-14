package handlers

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/phoen1xcode/phoen1xcodecloud/internal/models"
	"github.com/phoen1xcode/phoen1xcodecloud/pkg/storage"
	"github.com/phoen1xcode/phoen1xcodecloud/pkg/utils"
	"gorm.io/gorm"
)

const (
	maxFileSize         = 100 * 1024 * 1024 // 100MB
	maxTextContentSize  = 1024 * 1024       // 1MB
)

var allowedFileExtensions = map[string]bool{
	".txt": true, ".pdf": true, ".doc": true, ".docx": true,
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
	".zip": true, ".tar": true, ".gz": true,
	".mp4": true, ".mp3": true, ".wav": true,
	".csv": true, ".json": true, ".xml": true,
	".go": true, ".js": true, ".py": true, ".java": true,
	".html": true, ".css": true, ".md": true,
}

type ShareHandler struct {
	db      *gorm.DB
	storage storage.Storage
}

func NewShareHandler(db *gorm.DB, storage storage.Storage) *ShareHandler {
	return &ShareHandler{db: db, storage: storage}
}

func (h *ShareHandler) UploadFile(c *gin.Context) {
	userID := c.GetUint("user_id")
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	// Validate file size
	if file.Size > maxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("File size exceeds maximum allowed size of %d MB", maxFileSize/(1024*1024))})
		return
	}

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedFileExtensions[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File type not allowed"})
		return
	}

	shareCode, err := utils.GenerateShareCode()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate share code"})
		return
	}
	key := fmt.Sprintf("%d/%s/%s", userID, shareCode, file.Filename)

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer src.Close()

	if err := h.storage.Upload(c.Request.Context(), key, src); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
		return
	}

	share := models.Share{
		ShareCode: shareCode,
		UserID:    userID,
		Type:      "file",
		FileName:  file.Filename,
		FileSize:  file.Size,
		FilePath:  key,
	}

	if err := h.db.Create(&share).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create share"})
		return
	}

	c.JSON(http.StatusCreated, share)
}

func (h *ShareHandler) CreateTextShare(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate text content size
	if len(req.Content) > maxTextContentSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Text content exceeds maximum allowed size of %d KB", maxTextContentSize/1024)})
		return
	}

	shareCode, err := utils.GenerateShareCode()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate share code"})
		return
	}
	share := models.Share{
		ShareCode:   shareCode,
		UserID:      userID,
		Type:        "text",
		TextContent: req.Content,
	}

	if err := h.db.Create(&share).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create share"})
		return
	}

	c.JSON(http.StatusCreated, share)
}

func (h *ShareHandler) GetShare(c *gin.Context) {
	shareCode := c.Param("code")
	var share models.Share

	if err := h.db.Where("share_code = ?", shareCode).First(&share).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Share not found"})
		return
	}

	if share.ExpiresAt != nil && share.ExpiresAt.Before(time.Now()) {
		c.JSON(http.StatusGone, gin.H{"error": "Share has expired"})
		return
	}

	h.db.Model(&share).Update("downloads", share.Downloads+1)

	if share.Type == "file" {
		reader, err := h.storage.Download(c.Request.Context(), share.FilePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to download file"})
			return
		}
		defer reader.Close()

		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", share.FileName))
		c.DataFromReader(http.StatusOK, share.FileSize, "application/octet-stream", reader, nil)
		return
	}

	c.JSON(http.StatusOK, share)
}

func (h *ShareHandler) ListUserShares(c *gin.Context) {
	userID := c.GetUint("user_id")
	var shares []models.Share

	if err := h.db.Where("user_id = ?", userID).Order("created_at desc").Find(&shares).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch shares"})
		return
	}

	c.JSON(http.StatusOK, shares)
}

func (h *ShareHandler) DeleteShare(c *gin.Context) {
	userID := c.GetUint("user_id")
	shareCode := c.Param("code")
	var share models.Share

	if err := h.db.Where("share_code = ? AND user_id = ?", shareCode, userID).First(&share).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Share not found"})
		return
	}

	if share.Type == "file" {
		if err := h.storage.Delete(c.Request.Context(), share.FilePath); err != nil {
			log.Printf("Failed to delete file from storage: %v", err)
			// Continue with database deletion even if storage deletion fails
		}
	}

	if err := h.db.Delete(&share).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete share"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Share deleted"})
}
