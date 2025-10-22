// backend/controllers/auth_controller.go
package controllers

import (
	"net/http"
	"todolist/backend/database"
	"todolist/backend/models"
	"todolist/backend/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// RegisterInput defines the structure for user registration data
type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// POST /register
func Register(c *gin.Context) {
	var input RegisterInput

	// 1. Validate input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// 3. Create the user
	user := models.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: string(hashedPassword),
	}
	result := database.DB.Create(&user)

	// Handle potential errors, like duplicate email/username
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email or username already exists"})
		return
	}

	// Don't send the password hash back
	user.PasswordHash = ""

	c.JSON(http.StatusCreated, user)
}

// LoginInput defines the structure for user login data
type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// POST /login
func Login(c *gin.Context) {
	var input LoginInput
	var user models.User

	// 1. Validate input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Find user by email
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		// ใช้ error message กลางๆ เพื่อความปลอดภัย
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// 3. Compare password with the hash in the database
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// 4. Generate JWT token
	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	userResponse := gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
	}

	// 5. Send token back to the user
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  userResponse,
	})
	
}
