package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testCategoryID uint

func TestCreateCategory(t *testing.T) {
	if testToken == "" {
		// Crear usuario admin y obtener token
		registerPayload := map[string]interface{}{
			"username": "admin",
			"email":    "admin@test.com",
			"password": "admin123",
			"role":     "admin",
		}
		w := MakeRequest("POST", "/api/auth/register", registerPayload, "")
		var response map[string]interface{}
		ParseResponse(w, &response)
		testToken = response["token"].(string)
	}

	testCases := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name: "Crear categoría exitosamente",
			payload: map[string]interface{}{
				"name":        "Electrónica Test",
				"description": "Categoría de prueba",
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Contains(t, resp, "category")
				category := resp["category"].(map[string]interface{})
				assert.Equal(t, "Electrónica Test", category["name"])

				// Guardar ID para otros tests
				testCategoryID = uint(category["id"].(float64))
			},
		},
		{
			name: "Crear categoría sin nombre",
			payload: map[string]interface{}{
				"description": "Sin nombre",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Contains(t, resp, "error")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := MakeRequest("POST", "/api/categories", tc.payload, testToken)

			assert.Equal(t, tc.expectedStatus, w.Code)

			var response map[string]interface{}
			err := ParseResponse(w, &response)
			assert.NoError(t, err)

			tc.checkResponse(t, response)
		})
	}
}

func TestGetCategories(t *testing.T) {
	t.Run("Listar categorías", func(t *testing.T) {
		w := MakeRequest("GET", "/api/categories", nil, testToken)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := ParseResponse(w, &response)
		assert.NoError(t, err)

		assert.Contains(t, response, "categories")
		categories := response["categories"].([]interface{})
		assert.Greater(t, len(categories), 0)
	})
}

func TestGetCategory(t *testing.T) {
	if testCategoryID == 0 {
		t.Skip("No hay ID de categoría disponible")
	}

	t.Run("Obtener categoría por ID", func(t *testing.T) {
		url := fmt.Sprintf("/api/categories/%d", testCategoryID)
		w := MakeRequest("GET", url, nil, testToken)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := ParseResponse(w, &response)
		assert.NoError(t, err)

		assert.Contains(t, response, "category")
		category := response["category"].(map[string]interface{})
		assert.Equal(t, "Electrónica Test", category["name"])
	})

	t.Run("Obtener categoría inexistente", func(t *testing.T) {
		w := MakeRequest("GET", "/api/categories/99999", nil, testToken)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestUpdateCategory(t *testing.T) {
	if testCategoryID == 0 {
		t.Skip("No hay ID de categoría disponible")
	}

	t.Run("Actualizar categoría", func(t *testing.T) {
		url := fmt.Sprintf("/api/categories/%d", testCategoryID)
		payload := map[string]interface{}{
			"name":        "Electrónica Actualizada",
			"description": "Descripción actualizada",
		}

		w := MakeRequest("PUT", url, payload, testToken)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := ParseResponse(w, &response)
		assert.NoError(t, err)

		category := response["category"].(map[string]interface{})
		assert.Equal(t, "Electrónica Actualizada", category["name"])
	})
}

func TestDeleteCategory(t *testing.T) {
	if testCategoryID == 0 {
		t.Skip("No hay ID de categoría disponible")
	}

	t.Run("Eliminar categoría", func(t *testing.T) {
		url := fmt.Sprintf("/api/categories/%d", testCategoryID)
		w := MakeRequest("DELETE", url, nil, testToken)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := ParseResponse(w, &response)
		assert.NoError(t, err)

		assert.Contains(t, response, "message")
	})
}
