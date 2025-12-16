package api

import (
	"strconv"
	"warehouseControl/db"
	"warehouseControl/models"

	"github.com/gin-gonic/gin"
)

// Get all inventory items
func GetItems(c *gin.Context) {
	var items []models.Inventory
	if err := db.DB.Find(&items).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch items"})
		return
	}

	c.JSON(200, items)
}

// Get a specific inventory item
func GetItem(c *gin.Context) {
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

	c.JSON(200, item)
}
