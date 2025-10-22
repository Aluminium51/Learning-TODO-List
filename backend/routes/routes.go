// backend/routes/routes.go
package routes

import (
	"todolist/backend/controllers"
	"todolist/backend/middlewares" // 1. Import the middlewares package

	"github.com/gin-gonic/gin"
)

// SetupRoutes กำหนด routes ทั้งหมดสำหรับ API
func SetupRoutes(router *gin.Engine) {
	// จัดกลุ่ม API routes ภายใต้ /api/v1
	api := router.Group("/api/v1")
	{
		// --- Public Auth Routes ---
		// Routes เหล่านี้ไม่ต้องผ่าน Middleware
		api.POST("/register", controllers.Register)
		api.POST("/login", controllers.Login)

		// --- Protected Routes ---
		// 2. สร้าง Group ใหม่และใส่ Middleware เข้าไป
		// protected := api.Group("/").Use(middlewares.AuthMiddleware()) ผิด
		protected := api.Group("/", middlewares.AuthMiddleware())
		{
			// 3. ย้าย Todos group เข้ามาข้างในนี้
			protected.GET("/todos/stats", controllers.GetTodoStats)
			todos := protected.Group("/todos")
			{
				todos.GET("", controllers.GetTodos)
				todos.POST("", controllers.CreateTodo)
				todos.GET("/:id", controllers.GetTodoByID)
				todos.PUT("/:id", controllers.UpdateTodo)
				todos.DELETE("/:id", controllers.DeleteTodo)
				todos.POST("/:id/upload", controllers.UploadAttachment)
			}
		}
	}
}
