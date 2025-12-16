package main

import (
	"log"
	"warehouseControl/api"
	"warehouseControl/auth"
	"warehouseControl/db"
	"warehouseControl/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	if godotenv.Load() != nil {
		log.Fatal("Failed to load enviroment variables")
	}
	db.ConnectToDb()
}

func main() {
	r := gin.Default()

	// Serve static files
	r.Static("/static", "./static")

	// Serve the main page
	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	// Auth routes
	r.POST("/auth/signup", auth.Signup)
	r.POST("/auth/login", auth.Login)

	// Item routes
	r.POST("/items", middleware.RequireAdmin, api.CreateItem)
	r.GET("/items", api.GetItems)
	r.GET("/items/:id", api.GetItem)
	r.PUT("/items/:id", middleware.RequireManager, api.UpdateItem)
	r.DELETE("/items/:id", middleware.RequireManager, api.DeleteItem)

	// History routes
	r.GET("/items/:id/history", api.GetHistoryForItem)

	// Run the server
	r.Run(":8080")
}
