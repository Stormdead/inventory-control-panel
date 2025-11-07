package controllers

import (
	"net/http"
	"time"

	"github.com/Stormdead/inventory-control-panel/backend/config"
	"github.com/Stormdead/inventory-control-panel/backend/models"
	"github.com/gin-gonic/gin"
)

// GET /api/dashboard/stats - Estadísticas generales del inventario
func GetDashboardStats(c *gin.Context) {
	var stats struct {
		TotalProducts    int64   `json:"total_products"`
		TotalCategories  int64   `json:"total_categories"`
		TotalStock       int     `json:"total_stock"`
		LowStockProducts int64   `json:"low_stock_products"`
		TotalUsers       int64   `json:"total_users"`
		TotalValue       float64 `json:"total_inventory_value"`
	}

	config.DB.Model(&models.Product{}).Count(&stats.TotalProducts)
	config.DB.Model(&models.Category{}).Count(&stats.TotalCategories)
	config.DB.Model(&models.User{}).Count(&stats.TotalUsers)
	config.DB.Model(&models.Product{}).Where("stock < ?", 10).Count(&stats.LowStockProducts)

	var products []models.Product
	config.DB.Find(&products)

	for _, product := range products {
		stats.TotalStock += product.Stock
		stats.TotalValue += float64(product.Stock) * product.Price
	}

	c.JSON(http.StatusOK, gin.H{
		"stats": stats,
	})
}

// GET /api/dashboard/recent-movements - Movimientos recientes (últimos 10)
func GetRecentMovements(c *gin.Context) {
	var movements []models.Movement

	if err := config.DB.Preload("Product").Preload("Product.Category").Preload("User").
		Order("movement_date DESC").
		Limit(10).
		Find(&movements).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener movimientos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"movements": movements,
		"total":     len(movements),
	})
}

// GET /api/dashboard/low-stock-alerts - Alertas de productos con stock bajo
func GetLowStockAlerts(c *gin.Context) {
	var products []models.Product

	if err := config.DB.Preload("Category").
		Where("stock < ?", 10).
		Order("stock ASC").
		Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener productos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": products,
		"total":    len(products),
		"message":  "Productos que requieren reabastecimiento",
	})
}

// GET /api/dashboard/movement-summary - Resumen de movimientos (último mes)
func GetMovementSummary(c *gin.Context) {
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)

	var summary struct {
		TotalEntradas    int64 `json:"total_entradas"`
		TotalSalidas     int64 `json:"total_salidas"`
		CantidadEntradas int   `json:"cantidad_entradas"`
		CantidadSalidas  int   `json:"cantidad_salidas"`
	}

	config.DB.Model(&models.Movement{}).
		Where("type = ? AND movement_date >= ?", "entrada", thirtyDaysAgo).
		Count(&summary.TotalEntradas)

	config.DB.Model(&models.Movement{}).
		Where("type = ? AND movement_date >= ?", "salida", thirtyDaysAgo).
		Count(&summary.TotalSalidas)

	var entradas []models.Movement
	config.DB.Where("type = ? AND movement_date >= ?", "entrada", thirtyDaysAgo).Find(&entradas)
	for _, mov := range entradas {
		summary.CantidadEntradas += mov.Quantity
	}

	var salidas []models.Movement
	config.DB.Where("type = ? AND movement_date >= ?", "salida", thirtyDaysAgo).Find(&salidas)
	for _, mov := range salidas {
		summary.CantidadSalidas += mov.Quantity
	}

	c.JSON(http.StatusOK, gin.H{
		"summary": summary,
		"period":  "Últimos 30 días",
	})
}

// GET /api/dashboard/top-products - Top 5 productos con más movimientos
func GetTopProducts(c *gin.Context) {
	type ProductMovement struct {
		ProductID      uint   `json:"product_id"`
		ProductName    string `json:"product_name"`
		TotalMovements int64  `json:"total_movements"`
		CurrentStock   int    `json:"current_stock"`
		CategoryName   string `json:"category_name"`
	}

	var results []ProductMovement

	// Subconsulta para contar movimientos
	subQuery := config.DB.Model(&models.Movement{}).
		Select("product_id, COUNT(*) as movement_count").
		Group("product_id")

	// Query principal con joins elegantes
	err := config.DB.Table("products as p").
		Select(`
			p.id as product_id,
			p.name as product_name,
			COALESCE(m.movement_count, 0) as total_movements,
			p.stock as current_stock,
			COALESCE(c.name, 'Sin categoría') as category_name
		`).
		Joins("LEFT JOIN (?) as m ON p.id = m.product_id", subQuery).
		Joins("LEFT JOIN categories c ON p.category_id = c.id").
		Where("p.deleted_at IS NULL").
		Order("total_movements DESC").
		Limit(5).
		Scan(&results).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener estadísticas"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": results,
		"total":    len(results),
	})
}
