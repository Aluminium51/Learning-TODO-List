// backend/controllers/todos_controller.go
package controllers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"
	"todolist/backend/database"
	"todolist/backend/models"

	"github.com/gin-gonic/gin"
)

// GetUserIDFromContext is a helper function to get user ID from context
func GetUserIDFromContext(c *gin.Context) (uint, error) {
	id, exists := c.Get("user_id")
	if !exists {
		return 0, fmt.Errorf("user_id not found in context")
	}
	userID, ok := id.(uint)
	if !ok {
		return 0, fmt.Errorf("user_id is of invalid type")
	}
	return userID, nil
}

// GET /todos - ดึงรายการ Todo เฉพาะของผู้ใช้ที่ login
func GetTodos(c *gin.Context) {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var todos []models.Todo
	result := database.DB.Where("user_id = ?", userID).Find(&todos)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, todos)

	// var todos []models.Todo
	// // ค้นหา todos ทั้งหมดจากฐานข้อมูล
	// result := database.DB.Find(&todos)
	// if result.Error != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
	// 	return
	// }
	// c.JSON(http.StatusOK, todos)
}

// POST /todos - สร้าง Todo ใหม่
func CreateTodo(c *gin.Context) {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var newTodo models.Todo
	if err := c.ShouldBindJSON(&newTodo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newTodo.UserID = userID // Set UserID from the logged-in user

	result := database.DB.Create(&newTodo)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, newTodo)
	// var newTodo models.Todo
	// // Bind JSON ที่ส่งเข้ามากับ newTodo struct
	// if err := c.ShouldBindJSON(&newTodo); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	// // TODO: ในอนาคต เราจะกำหนด UserID จาก user ที่ login อยู่
	// // ตอนนี้ hardcode เป็น 1 ไปก่อน
	// newTodo.UserID = 1

	// // สร้าง record ใหม่ในฐานข้อมูล
	// result := database.DB.Create(&newTodo)
	// if result.Error != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
	// 	return
	// }

	// c.JSON(http.StatusCreated, newTodo)
}

// GET /todos/:id - ดึง Todo ตาม ID
func GetTodoByID(c *gin.Context) {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var todo models.Todo
	id := c.Param("id")

	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&todo).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found or you don't have permission"})
		return
	}

	c.JSON(http.StatusOK, todo)
	// var todo models.Todo
	// id := c.Param("id")

	// result := database.DB.First(&todo, id)
	// if result.Error != nil {
	// 	c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
	// 	return
	// }

	// c.JSON(http.StatusOK, todo)
}

// PUT /todos/:id - อัปเดต Todo
func UpdateTodo(c *gin.Context) {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var todo models.Todo
	id := c.Param("id")

	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&todo).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found or you don't have permission"})
		return
	}

	var updatedTodo models.Todo
	if err := c.ShouldBindJSON(&updatedTodo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Model(&todo).Updates(updatedTodo)
	c.JSON(http.StatusOK, todo)
	// var todo models.Todo
	// id := c.Param("id")

	// // ค้นหา todo ที่มีอยู่ก่อน
	// if err := database.DB.First(&todo, id).Error; err != nil {
	// 	c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
	// 	return
	// }

	// // Bind JSON ที่ส่งเข้ามาเพื่ออัปเดต
	// var updatedTodo models.Todo
	// if err := c.ShouldBindJSON(&updatedTodo); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	// // อัปเดตข้อมูลในฐานข้อมูล
	// database.DB.Model(&todo).Updates(updatedTodo)

	// c.JSON(http.StatusOK, todo)
}

// DELETE /todos/:id - ลบ Todo
func DeleteTodo(c *gin.Context) {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var todo models.Todo
	id := c.Param("id")

	// First, check if the todo exists and belongs to the user
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&todo).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found or you don't have permission"})
		return
	}

	// If it exists, delete it
	database.DB.Delete(&todo)
	c.JSON(http.StatusOK, gin.H{"message": "Todo deleted successfully"})
	// var todo models.Todo
	// id := c.Param("id")

	// // ลบข้อมูลออกจากฐานข้อมูล
	// result := database.DB.Delete(&todo, id)
	// if result.Error != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
	// 	return
	// }

	// // ตรวจสอบว่ามีการลบเกิดขึ้นจริงหรือไม่
	// if result.RowsAffected == 0 {
	// 	c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
	// 	return
	// }

	// c.JSON(http.StatusOK, gin.H{"message": "Todo deleted successfully"})

}

func UploadAttachment(c *gin.Context) {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var todo models.Todo
	id := c.Param("id")

	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&todo).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found or you don't have permission"})
		return
	}

	file, err := c.FormFile("attachment")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File not provided"})
		return
	}

	uniqueFileName := fmt.Sprintf("%d-%s", time.Now().Unix(), filepath.Base(file.Filename))
	destination := filepath.Join("uploads", uniqueFileName)

	if err := c.SaveUploadedFile(file, destination); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	todo.AttachmentURL = uniqueFileName
	database.DB.Save(&todo)

	c.JSON(http.StatusOK, todo)
	// var todo models.Todo
	// id := c.Param("id")

	// // 1. ค้นหา Todo ที่ต้องการแนบไฟล์
	// if err := database.DB.First(&todo, id).Error; err != nil {
	// 	c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
	// 	return
	// }

	// // 2. รับไฟล์จาก request
	// file, err := c.FormFile("attachment")
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "File not provided"})
	// 	return
	// }

	// // 3. สร้างชื่อไฟล์ใหม่ที่ไม่ซ้ำกัน (เพื่อป้องกันการเขียนทับ)
	// // เช่น: 1724436000-original_filename.png
	// uniqueFileName := fmt.Sprintf("%d-%s", time.Now().Unix(), filepath.Base(file.Filename))
	// destination := filepath.Join("uploads", uniqueFileName)

	// // 4. บันทึกไฟล์ลงใน server
	// if err := c.SaveUploadedFile(file, destination); err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
	// 	return
	// }

	// // 5. อัปเดตฐานข้อมูลด้วย URL ของไฟล์
	// // เราจะเก็บแค่ชื่อไฟล์ เพราะเราสามารถสร้าง URL เต็มๆ ได้จาก Frontend
	// todo.AttachmentURL = uniqueFileName
	// database.DB.Save(&todo)

	// c.JSON(http.StatusOK, todo)
}

func GetTodoStats(c *gin.Context) {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var totalTodos int64
	var completedTodos int64

	// นับจำนวน Todo ทั้งหมดของผู้ใช้
	database.DB.Model(&models.Todo{}).Where("user_id = ?", userID).Count(&totalTodos)

	// นับจำนวน Todo ที่เสร็จแล้วของผู้ใช้
	database.DB.Model(&models.Todo{}).Where("user_id = ? AND completed = ?", userID, true).Count(&completedTodos)

	// คำนวณที่ยังไม่เสร็จ
	activeTodos := totalTodos - completedTodos

	// ส่งข้อมูลกลับไปเป็น JSON
	c.JSON(http.StatusOK, gin.H{
		"total":     totalTodos,
		"completed": completedTodos,
		"active":    activeTodos,
	})
}
