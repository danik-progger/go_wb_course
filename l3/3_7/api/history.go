package api

import (
	"strconv"
	"warehouseControl/db"
	"warehouseControl/models"

	"github.com/gin-gonic/gin"
)

// GetHistoryForItem gets history records for a specific item
func GetHistoryForItem(c *gin.Context) {
	itemId, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil || itemId == 0 {
		c.JSON(400, gin.H{"error": "Invalid item ID"})
		return
	}

	var history []models.HistoryRecord
	if err := db.DB.Where("item_id = ?", uint(itemId)).Order("changed_at DESC").Find(&history).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch history"})
		return
	}

	c.JSON(200, history)
}
