package integration

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestSetUserActive(t *testing.T) {
	clearDB(t)
	teamBody := `{
        "team_name": "backend",
        "members": [{"user_id": "u1", "username": "Alice", "is_active": true}]
    }`
	performRequest(testRouter, "POST", "/team/add", teamBody)
	deactivateBody := `{"user_id": "u1", "is_active": false}`
	w := performRequest(testRouter, "POST", "/users/setIsActive", deactivateBody)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}
	var isActive bool
	err := testDB.Get(&isActive, "SELECT is_active FROM users WHERE user_id = 'u1'")
	if err != nil {
		t.Fatalf("Failed to query user: %v", err)
	}
	if isActive {
		t.Error("Expected user to be inactive")
	}
}

func TestGetUserReviews(t *testing.T) {
	clearDB(t)
	teamBody := `{
        "team_name": "backend",
        "members": [
            {"user_id": "u1", "username": "Alice", "is_active": true},
            {"user_id": "u2", "username": "Bob", "is_active": true}
        ]
    }`
	performRequest(testRouter, "POST", "/team/add", teamBody)
	prBody := `{"pull_request_id": "pr1", "pull_request_name": "Feature", "author_id": "u1"}`
	performRequest(testRouter, "POST", "/pullRequest/create", prBody)
	w := performRequest(testRouter, "GET", "/users/getReview?user_id=u2", "")
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	prs := response["pull_requests"].([]interface{})
	if len(prs) == 0 {
		t.Error("Expected at least 1 PR for reviewer u2")
	}
}
