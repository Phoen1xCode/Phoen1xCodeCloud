package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/phoen1xcode/phoen1xcodecloud/internal/models"
	"gorm.io/gorm"
)

type AdminHandler struct {
	db *gorm.DB
}

func NewAdminHandler(db *gorm.DB) *AdminHandler {
	return &AdminHandler{db: db}
}

func (h *AdminHandler) GetStats(c *gin.Context) {
	var userCount, shareCount, fileShareCount, textShareCount int64
	var totalFileSize int64

	h.db.Model(&models.User{}).Count(&userCount)
	h.db.Model(&models.Share{}).Count(&shareCount)
	h.db.Model(&models.Share{}).Where("type = ?", "file").Count(&fileShareCount)
	h.db.Model(&models.Share{}).Where("type = ?", "text").Count(&textShareCount)
	h.db.Model(&models.Share{}).Where("type = ?", "file").Select("COALESCE(SUM(file_size), 0)").Scan(&totalFileSize)

	c.JSON(http.StatusOK, gin.H{
		"users":           userCount,
		"total_shares":    shareCount,
		"file_shares":     fileShareCount,
		"text_shares":     textShareCount,
		"total_file_size": totalFileSize,
	})
}

func (h *AdminHandler) ListAllShares(c *gin.Context) {
	var shares []models.Share
	if err := h.db.Preload("User").Order("created_at desc").Find(&shares).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch shares"})
		return
	}
	c.JSON(http.StatusOK, shares)
}

func (h *AdminHandler) ListUsers(c *gin.Context) {
	var users []models.User
	if err := h.db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	c.JSON(http.StatusOK, users)
}
