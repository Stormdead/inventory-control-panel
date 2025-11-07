package controllers

import (
	"net/http"
	"time"

	"github.com/Stormdead/inventory-control-panel/backend/config"
	"github.com/Stormdead/inventory-control-panel/backend/models"
	"github.com/gin-gonic/gin"
)

// GET /api/movements - Listar todos los movimientos
func GetMovements(c *gin.Context) {
	var movements []models.Movement

	// Incluir relaciones con Product y User
	if err := config.DB.Preload("Product").Preload("Product.Category").Preload("User").Order("movement_date DESC").Find(&movements).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener movimientos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"movements": movements,
		"total":     len(movements),
	})
}

// GET /api/movements/:id - Obtener un movimiento por ID
func GetMovement(c *gin.Context) {
	id := c.Param("id")
	var movement models.Movement

	if err := config.DB.Preload("Product").Preload("Product.Category").Preload("User").First(&movement, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movimiento no encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"movement": movement,
	})
}

// POST /api/movements - Crear nuevo movimiento (entrada o salida)
func CreateMovement(c *gin.Context) {
	var movement models.Movement

	if err := c.ShouldBindJSON(&movement); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Obtener user_id del contexto (del token JWT)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}
	movement.UserID = userID.(uint)

	// Validaciones
	if movement.ProductID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El product_id es requerido"})
		return
	}

	if movement.Type != "entrada" && movement.Type != "salida" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El tipo debe ser 'entrada' o 'salida'"})
		return
	}

	if movement.Quantity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "La cantidad debe ser mayor a 0"})
		return
	}

	// Verificar que el producto existe
	var product models.Product
	if err := config.DB.First(&product, movement.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Producto no encontrado"})
		return
	}

	// Validar stock suficiente para salidas
	if movement.Type == "salida" {
		if product.Stock < movement.Quantity {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":           "Stock insuficiente",
				"stock_actual":    product.Stock,
				"cantidad_salida": movement.Quantity,
			})
			return
		}
	}

	// Establecer fecha actual si no se proporcionó
	if movement.MovementDate.IsZero() {
		movement.MovementDate = time.Now()
	}

	// Iniciar transacción para garantizar consistencia
	tx := config.DB.Begin()

	// Crear el movimiento
	if err := tx.Create(&movement).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear movimiento"})
		return
	}

	// Actualizar stock del producto
	if movement.Type == "entrada" {
		product.Stock += movement.Quantity
	} else { // salida
		product.Stock -= movement.Quantity
	}

	if err := tx.Save(&product).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar stock"})
		return
	}

	// Confirmar transacción
	tx.Commit()

	// Cargar relaciones para la respuesta
	config.DB.Preload("Product").Preload("Product.Category").Preload("User").First(&movement, movement.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message":     "Movimiento registrado exitosamente",
		"movement":    movement,
		"nuevo_stock": product.Stock,
	})
}

// GET /api/movements/product/:product_id - Movimientos de un producto específico
func GetMovementsByProduct(c *gin.Context) {
	productID := c.Param("product_id")
	var movements []models.Movement

	if err := config.DB.Preload("User").Where("product_id = ?", productID).Order("movement_date DESC").Find(&movements).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener movimientos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"movements": movements,
		"total":     len(movements),
	})
}

// GET /api/movements/type/:type - Movimientos por tipo (entrada/salida)
func GetMovementsByType(c *gin.Context) {
	movementType := c.Param("type")

	if movementType != "entrada" && movementType != "salida" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tipo inválido. Use 'entrada' o 'salida'"})
		return
	}

	var movements []models.Movement

	if err := config.DB.Preload("Product").Preload("Product.Category").Preload("User").Where("type = ?", movementType).Order("movement_date DESC").Find(&movements).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener movimientos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"movements": movements,
		"total":     len(movements),
		"type":      movementType,
	})
}

// DELETE /api/movements/:id - Eliminar movimiento (solo admin, no revierte stock)
func DeleteMovement(c *gin.Context) {
	id := c.Param("id")
	var movement models.Movement

	if err := config.DB.First(&movement, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movimiento no encontrado"})
		return
	}

	// NOTA: Este delete NO revierte el stock automáticamente
	// Si quieres revertir el stock, deberías hacerlo manualmente
	if err := config.DB.Delete(&movement).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar movimiento"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Movimiento eliminado exitosamente",
		"warning": "El stock del producto NO fue revertido automáticamente",
	})
}
