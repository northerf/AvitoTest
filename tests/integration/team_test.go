package integration

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestCreateTeam(t *testing.T) {
	clearDB(t)
	body := `{
        "team_name": "backend",
        "members": [
            {"user_id": "u1", "username": "Alice", "is_active": true},
            {"user_id": "u2", "username": "Bob", "is_active": true}
        ]
    }`
	w := performRequest(testRouter, "POST", "/team/add", body)
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d. Body: %s", w.Code, w.Body.String())
	}
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	team := response["team"].(map[string]interface{})
	if team["team_name"] != "backend" {
		t.Errorf("Expected team_name 'backend', got %v", team["team_name"])
	}
}

func TestGetTeam(t *testing.T) {
	clearDB(t)
	createBody := `{
        "team_name": "backend",
        "members": [{"user_id": "u1", "username": "Alice", "is_active": true}]
    }`
	performRequest(testRouter, "POST", "/team/add", createBody)
	w := performRequest(testRouter, "GET", "/team/get?team_name=backend", "")
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	if response["team_name"] != "backend" {
		t.Errorf("Expected team_name 'backend', got %v", response["team_name"])
	}
}
