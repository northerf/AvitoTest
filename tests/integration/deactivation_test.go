package integration

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestDeactivateUsersAndReassign(t *testing.T) {
	clearDB(t)
	teamBody := `{
        "team_name": "backend",
        "members": [
            {"user_id": "u1", "username": "Alice", "is_active": true},
            {"user_id": "u2", "username": "Bob", "is_active": true},
            {"user_id": "u3", "username": "Charlie", "is_active": true},
            {"user_id": "u4", "username": "David", "is_active": true}
        ]
    }`
	performRequest(testRouter, "POST", "/team/add", teamBody)
	pr1Body := `{"pull_request_id": "pr1", "pull_request_name": "Feature A", "author_id": "u1"}`
	performRequest(testRouter, "POST", "/pullRequest/create", pr1Body)
	pr2Body := `{"pull_request_id": "pr2", "pull_request_name": "Feature B", "author_id": "u2"}`
	performRequest(testRouter, "POST", "/pullRequest/create", pr2Body)
	deactivateBody := `{"team_name": "backend", "user_ids": ["u1", "u2"]}`
	w := performRequest(testRouter, "POST", "/team/deactivate", deactivateBody)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	deactivatedCount := int(response["deactivated_count"].(float64))
	if deactivatedCount != 2 {
		t.Errorf("Expected 2 deactivated users, got %d", deactivatedCount)
	}
	var count int
	err := testDB.Get(&count, `SELECT COUNT(*) FROM pr_reviewers prr
        INNER JOIN pull_requests pr ON prr.pull_request_id = pr.pull_request_id
        WHERE pr.status = 'OPEN' AND prr.reviewer_id IN ('u1', 'u2')`)
	if err != nil {
		t.Fatalf("Failed to query: %v", err)
	}
	if count != 0 {
		t.Errorf("Expected 0 deactivated reviewers, got %d", count)
	}
}
