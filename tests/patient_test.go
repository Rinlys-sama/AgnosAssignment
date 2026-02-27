package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// === Staff Login Tests ===

func TestLogin_Success(t *testing.T) {
	router, db := setupTestRouter(t)
	defer db.Close()
	defer cleanupTestData(db)

	createBody := map[string]string{
		"username": "test_login_user",
		"password": "password123",
		"hospital": "hospital_a",
	}
	createJSON, _ := json.Marshal(createBody)
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/staff/create", bytes.NewBuffer(createJSON))
	req1.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w1, req1)

	loginBody := map[string]string{
		"username": "test_login_user",
		"password": "password123",
		"hospital": "hospital_a",
	}
	loginJSON, _ := json.Marshal(loginBody)
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/staff/login", bytes.NewBuffer(loginJSON))
	req2.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", w2.Code, w2.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w2.Body.Bytes(), &resp)
	if resp["token"] == nil || resp["token"] == "" {
		t.Error("Expected a JWT token in response")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	router, db := setupTestRouter(t)
	defer db.Close()
	defer cleanupTestData(db)

	createBody := map[string]string{
		"username": "test_wrong_pw",
		"password": "password123",
		"hospital": "hospital_a",
	}
	createJSON, _ := json.Marshal(createBody)
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/staff/create", bytes.NewBuffer(createJSON))
	req1.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w1, req1)

	loginBody := map[string]string{
		"username": "test_wrong_pw",
		"password": "wrong_password",
		"hospital": "hospital_a",
	}
	loginJSON, _ := json.Marshal(loginBody)
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/staff/login", bytes.NewBuffer(loginJSON))
	req2.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401 for wrong password, got %d", w2.Code)
	}
}

func TestLogin_NonexistentUser(t *testing.T) {
	router, db := setupTestRouter(t)
	defer db.Close()

	loginBody := map[string]string{
		"username": "nonexistent_user",
		"password": "password123",
		"hospital": "hospital_a",
	}
	loginJSON, _ := json.Marshal(loginBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/staff/login", bytes.NewBuffer(loginJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401 for nonexistent user, got %d", w.Code)
	}
}

// === Patient Search Tests ===

func TestPatientSearch_WithoutAuth(t *testing.T) {
	router, db := setupTestRouter(t)
	defer db.Close()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/patient/search", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401 without auth, got %d", w.Code)
	}
}

func TestPatientSearch_WithInvalidToken(t *testing.T) {
	router, db := setupTestRouter(t)
	defer db.Close()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/patient/search", nil)
	req.Header.Set("Authorization", "Bearer invalid-token-here")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401 with invalid token, got %d", w.Code)
	}
}

func TestPatientSearch_Success(t *testing.T) {
	router, db := setupTestRouter(t)
	defer db.Close()
	defer cleanupTestData(db)

	createBody := map[string]string{
		"username": "test_search_user",
		"password": "password123",
		"hospital": "hospital_a",
	}
	createJSON, _ := json.Marshal(createBody)
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/staff/create", bytes.NewBuffer(createJSON))
	req1.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w1, req1)

	loginJSON, _ := json.Marshal(createBody)
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/staff/login", bytes.NewBuffer(loginJSON))
	req2.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w2, req2)

	var loginResp map[string]interface{}
	json.Unmarshal(w2.Body.Bytes(), &loginResp)
	token := loginResp["token"].(string)

	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("GET", "/patient/search?national_id=1234567890123", nil)
	req3.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w3, req3)

	if w3.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", w3.Code, w3.Body.String())
	}

	var searchResp map[string]interface{}
	json.Unmarshal(w3.Body.Bytes(), &searchResp)
	patients, ok := searchResp["patients"].([]interface{})
	if !ok {
		t.Error("Expected 'patients' array in response")
	}

	if len(patients) == 0 {
		t.Error("Expected at least one patient result for national_id=1234567890123")
	}
}

func TestPatientSearch_HospitalIsolation(t *testing.T) {
	router, db := setupTestRouter(t)
	defer db.Close()
	defer cleanupTestData(db)

	createBody := map[string]string{
		"username": "test_isolation_user",
		"password": "password123",
		"hospital": "hospital_b",
	}
	createJSON, _ := json.Marshal(createBody)
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/staff/create", bytes.NewBuffer(createJSON))
	req1.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w1, req1)

	loginJSON, _ := json.Marshal(createBody)
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/staff/login", bytes.NewBuffer(loginJSON))
	req2.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w2, req2)

	var loginResp map[string]interface{}
	json.Unmarshal(w2.Body.Bytes(), &loginResp)
	token := loginResp["token"].(string)

	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("GET", "/patient/search?national_id=1234567890123", nil)
	req3.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w3, req3)

	if w3.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w3.Code)
	}

	var searchResp map[string]interface{}
	json.Unmarshal(w3.Body.Bytes(), &searchResp)
	patients := searchResp["patients"].([]interface{})

	if len(patients) != 0 {
		t.Errorf("Hospital isolation failed! Staff from hospital_b should NOT see hospital_a patients. Got %d results", len(patients))
	}
}
