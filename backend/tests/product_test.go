package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testProductID uint

func TestCreateProduct(t *testing.T) {
	// Primero crear una categoría
	categoryPayload := map[string]interface{}{
		"name":        "Test Category",
		"description": "For product tests",
	}
	w := MakeRequest("POST", "/api/categories", categoryPayload, testToken)
	var catResp map[string]interface{}
	ParseResponse(w, &catResp)
	categoryID := uint(catResp["category"].(map[string]interface{})["id"].(float64))

	t.Run("Crear producto exitosamente", func(t *testing.T) {
		payload := map[string]interface{}{
			"name":        "Producto Test",
			"description": "Descripción de prueba",
			"category_id": categoryID,
			"price":       99.99,
			"stock":       10,
		}

		w := MakeRequest("POST", "/api/products", payload, testToken)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		ParseResponse(w, &response)

		assert.Contains(t, response, "product")
		product := response["product"].(map[string]interface{})
		testProductID = uint(product["id"].(float64))
		assert.Equal(t, "Producto Test", product["name"])
	})

	t.Run("Crear producto sin precio", func(t *testing.T) {
		payload := map[string]interface{}{
			"name": "Producto sin precio",
		}

		w := MakeRequest("POST", "/api/products", payload, testToken)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGetProducts(t *testing.T) {
	t.Run("Listar productos", func(t *testing.T) {
		w := MakeRequest("GET", "/api/products", nil, testToken)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		ParseResponse(w, &response)

		assert.Contains(t, response, "products")
	})
}

func TestGetLowStockProducts(t *testing.T) {
	t.Run("Obtener productos con stock bajo", func(t *testing.T) {
		w := MakeRequest("GET", "/api/products/low-stock", nil, testToken)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		ParseResponse(w, &response)

		assert.Contains(t, response, "products")
	})
}
