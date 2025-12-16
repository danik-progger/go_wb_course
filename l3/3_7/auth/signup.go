package auth

import (
	"net/http"
	"warehouseControl/db"
	model "warehouseControl/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	// Get email and password from req body
	var body struct {
		Email    string
		Password string
		Role     string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body"})
		return
	}

	// Check role
	if body.Role != "admin" && body.Role != "manager" && body.Role != "user" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role"})
		return
	}

	// Encrypt a password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to hash pasword"})
		return
	}

	// Add user to DB
	user := model.User{Email: body.Email, Password: string(hash), Role: body.Role}
	result := db.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
