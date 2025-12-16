package api

import (
	"fmt"
	"strconv"
	"time"
	"warehouseControl/db"
	"warehouseControl/models"

	"github.com/gin-gonic/gin"
)

// Update an existing inventory item
func UpdateItem(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil || id == 0 {
		c.JSON(400, gin.H{"error": "Invalid item ID"})
		return
	}

	var updatedItem models.Inventory
	if err := c.ShouldBindJSON(&updatedItem); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Validate required fields
	if updatedItem.Name == "" {
		c.JSON(400, gin.H{"error": "Name is required"})
		return
	}

	var existingItem models.Inventory
	if err := db.DB.First(&existingItem, uint(id)).Error; err != nil {
		c.JSON(404, gin.H{"error": "Item not found"})
		return
	}

	// Preserve the ID and update timestamps
	updatedItem.ID = existingItem.ID
	updatedItem.CreatedAt = existingItem.CreatedAt
	updatedItem.UpdatedAt = time.Now()

	if err := db.DB.Save(&updatedItem).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update item"})
		return
	}

	// Log the update in history
	history := models.HistoryRecord{
		ItemID:    updatedItem.ID,
		Action:    "UPDATE",
		OldValues: fmt.Sprintf("%+v", existingItem),
		NewValues: fmt.Sprintf("%+v", updatedItem),
		ChangedBy: "system",
		ChangedAt: time.Now(),
	}
	db.DB.Create(&history)

	c.JSON(200, updatedItem)
}
