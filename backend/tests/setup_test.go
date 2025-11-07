package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/Stormdead/inventory-control-panel/backend/config"
	"github.com/Stormdead/inventory-control-panel/backend/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var router *gin.Engine
var testToken string

// SetupTestEnvironment configura el entorno de testing
func SetupTestEnvironment() {
	// Cargar variables de entorno
	godotenv.Load("../.env")

	// Usar base de datos de test (puedes crear una separada)
	//os.Setenv("DB_NAME", "inventory_db_test")

	// Conectar a la base de datos
	config.ConnectDB()

	// Configurar Gin en modo test
	gin.SetMode(gin.TestMode)
	router = gin.Default()

	// Configurar rutas
	routes.SetupRoutes(router)
}

// MakeRequest es un helper para hacer peticiones HTTP
func MakeRequest(method, url string, body interface{}, token string) *httptest.ResponseRecorder {
	var reqBody []byte
	var err error

	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			panic(err)
		}
	}

	req, _ := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	return w
}

// ParseResponse es un helper para parsear respuestas JSON
func ParseResponse(w *httptest.ResponseRecorder, target interface{}) error {
	return json.Unmarshal(w.Body.Bytes(), target)
}

// CleanupDatabase limpia la base de datos después de los tests
func CleanupDatabase() {
	config.DB.Exec("DELETE FROM movements")
	config.DB.Exec("DELETE FROM products")
	config.DB.Exec("DELETE FROM categories")
	config.DB.Exec("DELETE FROM users")
}
