package api

import (
	"fmt"
	"strconv"
	"time"
	"warehouseControl/db"
	"warehouseControl/models"

	"github.com/gin-gonic/gin"
)

// Delete an inventory item
func DeleteItem(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil || id == 0 {
		c.JSON(400, gin.H{"error": "Invalid item ID"})
		return
	}

	var item models.Inventory
	if err := db.DB.First(&item, uint(id)).Error; err != nil {
		c.JSON(404, gin.H{"error": "Item not found"})
		return
	}

	if err := db.DB.Delete(&item).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete item"})
		return
	}

	// Log the deletion in history
	history := models.HistoryRecord{
		ItemID:    item.ID,
		Action:    "DELETE",
		OldValues: fmt.Sprintf("%+v", item),
		ChangedBy: "system",
		ChangedAt: time.Now(),
	}
	db.DB.Create(&history)

	c.JSON(204, gin.H{}) // Return 204 No Content with empty JSON instead of nil
}
