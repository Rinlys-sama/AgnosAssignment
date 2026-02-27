package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/agnos/hospital-middleware/config"
	"github.com/agnos/hospital-middleware/routes"
	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func setupTestRouter(t *testing.T) (*gin.Engine, *sql.DB) {
	cfg := config.LoadConfig()
	db := config.ConnectDB(cfg)
	router := routes.SetupRouter(db, cfg)
	return router, db
}

func cleanupTestData(db *sql.DB) {
	db.Exec("DELETE FROM staff WHERE username LIKE 'test_%'")
}

// === Staff Create Tests ===

func TestCreateStaff_Success(t *testing.T) {
	router, db := setupTestRouter(t)
	defer db.Close()
	defer cleanupTestData(db)

	body := map[string]string{
		"username": "test_user_create",
		"password": "password123",
		"hospital": "hospital_a",
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/staff/create", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d. Body: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	if resp["username"] != "test_user_create" {
		t.Errorf("Expected username 'test_user_create', got '%v'", resp["username"])
	}
	if resp["hospital"] != "hospital_a" {
		t.Errorf("Expected hospital 'hospital_a', got '%v'", resp["hospital"])
	}
}

func TestCreateStaff_DuplicateUsername(t *testing.T) {
	router, db := setupTestRouter(t)
	defer db.Close()
	defer cleanupTestData(db)

	body := map[string]string{
		"username": "test_user_dup",
		"password": "password123",
		"hospital": "hospital_a",
	}
	jsonBody, _ := json.Marshal(body)

	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/staff/create", bytes.NewBuffer(jsonBody))
	req1.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w1, req1)

	w2 := httptest.NewRecorder()
	jsonBody2, _ := json.Marshal(body)
	req2, _ := http.NewRequest("POST", "/staff/create", bytes.NewBuffer(jsonBody2))
	req2.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusConflict {
		t.Errorf("Expected status 409 for duplicate username, got %d. Body: %s", w2.Code, w2.Body.String())
	}
}

func TestCreateStaff_MissingFields(t *testing.T) {
	router, db := setupTestRouter(t)
	defer db.Close()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/staff/create", bytes.NewBuffer([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for missing fields, got %d", w.Code)
	}
}

func TestCreateStaff_InvalidHospital(t *testing.T) {
	router, db := setupTestRouter(t)
	defer db.Close()

	body := map[string]string{
		"username": "test_user_invalid_hospital",
		"password": "password123",
		"hospital": "nonexistent_hospital",
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/staff/create", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid hospital, got %d", w.Code)
	}
}
