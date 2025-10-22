// main.go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"

	"todolist/backend/database"
	"todolist/backend/models"
	"todolist/backend/routes"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using default environment variables")
	}
}

// SeedDatabase creates initial data if it doesn't exist
func SeedDatabase() {
	var user models.User
	// Check if the first user exists
	if err := database.DB.First(&user).Error; err == gorm.ErrRecordNotFound {
		fmt.Println("No users found. Seeding database...")
		// Create a default user
		defaultUser := models.User{
			Username:     "default_user",
			Email:        "default@example.com",
			PasswordHash: "not_a_real_hash", // In a real app, this would be a proper hash
		}
		if err := database.DB.Create(&defaultUser).Error; err != nil {
			log.Fatalf("could not seed database: %v", err)
		}
		fmt.Println("Database seeded successfully.")
	}
}

// --- Main Function ---
func main() {
	LoadEnv()

	database.ConnectToDatabase()
	database.MigrateDatabase()
	// SeedDatabase()

	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		os.Mkdir("uploads", os.ModePerm)
	}

	r := gin.Default()

	r.Static("/uploads", "./uploads")

	// --- เพิ่ม CORS Middleware ---
	// อนุญาตให้ Frontend ที่รันบน localhost:3000 (Next.js) สามารถเรียก API นี้ได้
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"} // Port มาตรฐานของ Next.js dev server
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(config))

	routes.SetupRoutes(r)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	fmt.Println("Starting server on port 8080...")
	err := r.Run(":8080")
	if err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
