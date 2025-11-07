package controllers

import (
	"net/http"

	"github.com/Stormdead/inventory-control-panel/backend/config"
	"github.com/Stormdead/inventory-control-panel/backend/models"
	"github.com/gin-gonic/gin"
)

// GET /api/categories - listar todas las categorías
func GetCategories(c *gin.Context) {
	var categories []models.Category

	if err := config.DB.Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener categorias"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"categories": categories})
}

// GET /api/categories/:id - Obtener una categoría por ID
func GetCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.Category

	if err := config.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Categoría no encontrada"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"category": category,
	})
}

// POST /api/categories - crear una nueva categoría(solo admin)
func CreateCategory(c *gin.Context) {
	var category models.Category

	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Validar si el nombre esta vacio
	if category.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El nombre de la categoria no puede estar vacio"})
		return
	}

	//Crear categoria
	if err := config.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear categoria"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Categoria creada exitosamente",
		"category": category,
	})
}

// PUT /api/categories/:id - actualizar una categoría(solo admin)
func UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.Category

	//Verificar si la categoria existe
	if err := config.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Categoria no encontrada"})
		return
	}

	//Obtener los datos de actualizacion
	var updateData models.Category
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Actualizar campos
	if updateData.Name != "" {
		category.Name = updateData.Name
	}
	category.Description = updateData.Description

	if err := config.DB.Save(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar la categoria"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Categoria actualizada exitosamente",
		"category": category,
	})
}

// DELETE /api/categories/:id - eliminar una categoría(solo admin)
func DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.Category

	//Verificar si la categoria existe
	if err := config.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Categoria no encontrada"})
		return
	}

	//Eliminar categoria
	if err := config.DB.Delete(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar la categoria"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Categoria eliminada exitosamente"})
}
