package routes

import (
	"github.com/Stormdead/inventory-control-panel/backend/controllers"
	"github.com/Stormdead/inventory-control-panel/backend/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		// Rutas de autenticación (públicas)
		auth := api.Group("/auth")
		{
			auth.POST("/register", controllers.Register)
			auth.POST("/login", controllers.Login)
			auth.GET("/profile", middleware.AuthMiddleware(), controllers.GetProfile)
		}

		// Rutas de categorías
		categories := api.Group("/categories")
		categories.Use(middleware.AuthMiddleware())
		{
			categories.GET("", controllers.GetCategories)
			categories.GET("/:id", controllers.GetCategory)
			categories.POST("", middleware.AdminMiddleware(), controllers.CreateCategory)
			categories.PUT("/:id", middleware.AdminMiddleware(), controllers.UpdateCategory)
			categories.DELETE("/:id", middleware.AdminMiddleware(), controllers.DeleteCategory)
		}

		// Rutas de dashboard (requieren autenticación)
		dashboard := api.Group("/dashboard")
		dashboard.Use(middleware.AuthMiddleware())
		{
			dashboard.GET("/stats", controllers.GetDashboardStats)
			dashboard.GET("/recent-movements", controllers.GetRecentMovements)
			dashboard.GET("/low-stock-alerts", controllers.GetLowStockAlerts)
			dashboard.GET("/movement-summary", controllers.GetMovementSummary)
			dashboard.GET("/top-products", controllers.GetTopProducts)
		}

		// Rutas de productos
		products := api.Group("/products")
		products.Use(middleware.AuthMiddleware())
		{
			products.GET("", controllers.GetProducts)
			products.GET("/low-stock", controllers.GetLowStockProducts)
			products.GET("/category/:category_id", controllers.GetProductsByCategory)
			products.GET("/:id", controllers.GetProduct)
			products.POST("", middleware.AdminMiddleware(), controllers.CreateProduct)
			products.PUT("/:id", controllers.UpdateProduct)
			products.DELETE("/:id", middleware.AdminMiddleware(), controllers.DeleteProduct)
		}

		// Rutas de movimientos
		movements := api.Group("/movements")
		movements.Use(middleware.AuthMiddleware())
		{
			movements.GET("", controllers.GetMovements)                                        // Listar todos
			movements.GET("/type/:type", controllers.GetMovementsByType)                       // Por tipo
			movements.GET("/product/:product_id", controllers.GetMovementsByProduct)           // Por producto
			movements.GET("/:id", controllers.GetMovement)                                     // Obtener uno
			movements.POST("", controllers.CreateMovement)                                     // Crear movimiento
			movements.DELETE("/:id", middleware.AdminMiddleware(), controllers.DeleteMovement) // Eliminar (admin)
		}
	}
}
