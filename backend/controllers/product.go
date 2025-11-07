package controllers

import (
	"net/http"

	"github.com/Stormdead/inventory-control-panel/backend/config"
	"github.com/Stormdead/inventory-control-panel/backend/models"
	"github.com/gin-gonic/gin"
)

// GET /api/products - Listar todos los productos
func GetProducts(c *gin.Context) {
	var products []models.Product

	// Incluir la relación con Category
	if err := config.DB.Preload("Category").Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener productos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": products,
		"total":    len(products),
	})
}

// GET /api/products/:id - Obtener un producto por ID
func GetProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	if err := config.DB.Preload("Category").First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Producto no encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"product": product,
	})
}

// POST /api/products - Crear nuevo producto
func CreateProduct(c *gin.Context) {
	var product models.Product

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validaciones
	if product.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El nombre es requerido"})
		return
	}

	if product.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El precio debe ser mayor a 0"})
		return
	}

	// Verificar que la categoría existe (si se proporcionó)
	if product.CategoryID != nil {
		var category models.Category
		if err := config.DB.First(&category, *product.CategoryID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "La categoría especificada no existe"})
			return
		}
	}

	// Crear producto
	if err := config.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear producto"})
		return
	}

	// Cargar la categoría para la respuesta
	config.DB.Preload("Category").First(&product, product.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Producto creado exitosamente",
		"product": product,
	})
}

// PUT /api/products/:id - Actualizar producto
func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	// Verificar que el producto existe
	if err := config.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Producto no encontrado"})
		return
	}

	// Obtener datos de actualización
	var updateData models.Product
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Actualizar campos
	if updateData.Name != "" {
		product.Name = updateData.Name
	}
	if updateData.Description != "" {
		product.Description = updateData.Description
	}
	if updateData.Price > 0 {
		product.Price = updateData.Price
	}
	if updateData.Stock >= 0 {
		product.Stock = updateData.Stock
	}
	if updateData.ImageURL != "" {
		product.ImageURL = updateData.ImageURL
	}
	if updateData.CategoryID != nil {
		// Verificar que la categoría existe
		var category models.Category
		if err := config.DB.First(&category, *updateData.CategoryID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "La categoría especificada no existe"})
			return
		}
		product.CategoryID = updateData.CategoryID
	}

	if err := config.DB.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar producto"})
		return
	}

	// Cargar la categoría para la respuesta
	config.DB.Preload("Category").First(&product, product.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Producto actualizado exitosamente",
		"product": product,
	})
}

// DELETE /api/products/:id - Eliminar producto
func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	// Verificar que el producto existe
	if err := config.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Producto no encontrado"})
		return
	}

	// Eliminar producto (soft delete por el DeletedAt en el modelo)
	if err := config.DB.Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar producto"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Producto eliminado exitosamente",
	})
}

// GET /api/products/low-stock - Productos con stock bajo (menos de 10 unidades)
func GetLowStockProducts(c *gin.Context) {
	var products []models.Product

	if err := config.DB.Preload("Category").Where("stock < ?", 10).Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener productos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": products,
		"total":    len(products),
	})
}

// GET /api/products/category/:category_id - Productos por categoría
func GetProductsByCategory(c *gin.Context) {
	categoryID := c.Param("category_id")
	var products []models.Product

	if err := config.DB.Preload("Category").Where("category_id = ?", categoryID).Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener productos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": products,
		"total":    len(products),
	})
}
