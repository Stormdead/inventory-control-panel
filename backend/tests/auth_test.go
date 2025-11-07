package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	SetupTestEnvironment()
	CleanupDatabase()
	m.Run()
}

func TestRegister(t *testing.T) {
	CleanupDatabase()

	testCases := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name: "Registro exitoso",
			payload: map[string]interface{}{
				"username": "testuser",
				"email":    "test@example.com",
				"password": "password123",
				"role":     "employee",
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Contains(t, resp, "token")
				assert.Contains(t, resp, "user")
				user := resp["user"].(map[string]interface{})
				assert.Equal(t, "testuser", user["username"])
				assert.Equal(t, "test@example.com", user["email"])
			},
		},
		{
			name: "Email duplicado",
			payload: map[string]interface{}{
				"username": "testuser2",
				"email":    "test@example.com",
				"password": "password123",
			},
			expectedStatus: http.StatusConflict,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Contains(t, resp, "error")
			},
		},
		{
			name: "Contraseña muy corta",
			payload: map[string]interface{}{
				"username": "testuser3",
				"email":    "test3@example.com",
				"password": "123",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Contains(t, resp, "error")
			},
		},
		{
			name: "Email inválido",
			payload: map[string]interface{}{
				"username": "testuser4",
				"email":    "invalid-email",
				"password": "password123",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Contains(t, resp, "error")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := MakeRequest("POST", "/api/auth/register", tc.payload, "")

			assert.Equal(t, tc.expectedStatus, w.Code)

			var response map[string]interface{}
			err := ParseResponse(w, &response)
			assert.NoError(t, err)

			tc.checkResponse(t, response)
		})
	}
}

func TestLogin(t *testing.T) {
	CleanupDatabase()

	// Primero registrar un usuario
	registerPayload := map[string]interface{}{
		"username": "logintest",
		"email":    "login@example.com",
		"password": "password123",
		"role":     "admin",
	}
	MakeRequest("POST", "/api/auth/register", registerPayload, "")

	testCases := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name: "Login exitoso",
			payload: map[string]interface{}{
				"email":    "login@example.com",
				"password": "password123",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Contains(t, resp, "token")
				assert.Contains(t, resp, "user")

				// Guardar token para tests posteriores
				testToken = resp["token"].(string)
			},
		},
		{
			name: "Contraseña incorrecta",
			payload: map[string]interface{}{
				"email":    "login@example.com",
				"password": "wrongpassword",
			},
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Contains(t, resp, "error")
			},
		},
		{
			name: "Usuario no existe",
			payload: map[string]interface{}{
				"email":    "noexiste@example.com",
				"password": "password123",
			},
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Contains(t, resp, "error")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := MakeRequest("POST", "/api/auth/login", tc.payload, "")

			assert.Equal(t, tc.expectedStatus, w.Code)

			var response map[string]interface{}
			err := ParseResponse(w, &response)
			assert.NoError(t, err)

			tc.checkResponse(t, response)
		})
	}
}

func TestGetProfile(t *testing.T) {
	if testToken == "" {
		t.Skip("No hay token disponible. Ejecuta TestLogin primero")
	}

	t.Run("Obtener perfil con token válido", func(t *testing.T) {
		w := MakeRequest("GET", "/api/auth/profile", nil, testToken)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := ParseResponse(w, &response)
		assert.NoError(t, err)

		assert.Contains(t, response, "user")
		user := response["user"].(map[string]interface{})
		assert.Equal(t, "login@example.com", user["email"])
	})

	t.Run("Obtener perfil sin token", func(t *testing.T) {
		w := MakeRequest("GET", "/api/auth/profile", nil, "")

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Obtener perfil con token inválido", func(t *testing.T) {
		w := MakeRequest("GET", "/api/auth/profile", nil, "token_invalido")

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
