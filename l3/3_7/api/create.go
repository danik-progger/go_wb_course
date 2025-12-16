package api

import (
	"fmt"
	"time"
	"warehouseControl/db"
	"warehouseControl/models"

	"github.com/gin-gonic/gin"
)

func CreateItem(c *gin.Context) {
	var item models.Inventory
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Validate required fields
	if item.Name == "" {
		c.JSON(400, gin.H{"error": "Name is required"})
		return
	}

	// Set created/updated times
	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()

	if err := db.DB.Create(&item).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create item"})
		return
	}

	// Log the creation in history
	history := models.HistoryRecord{
		ItemID:    item.ID,
		Action:    "INSERT",
		NewValues: fmt.Sprintf("%+v", item),
		ChangedBy: "system",
		ChangedAt: time.Now(),
	}
	db.DB.Create(&history)

	c.JSON(201, item)
}
