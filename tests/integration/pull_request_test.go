package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func TestCreatePullRequest(t *testing.T) {
	clearDB(t)
	teamBody := `{
        "team_name": "backend",
        "members": [
            {"user_id": "u1", "username": "Alice", "is_active": true},
            {"user_id": "u2", "username": "Bob", "is_active": true}
        ]
    }`
	performRequest(testRouter, "POST", "/team/add", teamBody)
	prBody := `{"pull_request_id": "pr1", "pull_request_name": "Feature A", "author_id": "u1"}`
	w := performRequest(testRouter, "POST", "/pullRequest/create", prBody)
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d. Body: %s", w.Code, w.Body.String())
	}
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	pr := response["pr"].(map[string]interface{})
	if pr["pull_request_id"] != "pr1" {
		t.Errorf("Expected PR ID 'pr1', got %v", pr["pull_request_id"])
	}
	reviewers := pr["assigned_reviewers"].([]interface{})
	if len(reviewers) == 0 {
		t.Error("Expected at least 1 reviewer assigned")
	}
}

func TestMergePullRequest(t *testing.T) {
	clearDB(t)
	teamBody := `{
        "team_name": "backend",
        "members": [{"user_id": "u1", "username": "Alice", "is_active": true}]
    }`
	performRequest(testRouter, "POST", "/team/add", teamBody)
	prBody := `{"pull_request_id": "pr1", "pull_request_name": "Feature", "author_id": "u1"}`
	performRequest(testRouter, "POST", "/pullRequest/create", prBody)
	mergeBody := `{"pull_request_id": "pr1"}`
	w := performRequest(testRouter, "POST", "/pullRequest/merge", mergeBody)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}
	var status string
	err := testDB.Get(&status, "SELECT status FROM pull_requests WHERE pull_request_id = 'pr1'")
	if err != nil {
		t.Fatalf("Failed to query PR: %v", err)
	}
	if status != "MERGED" {
		t.Errorf("Expected status 'MERGED', got %s", status)
	}
}

func TestReassignReviewer(t *testing.T) {
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
	prBody := `{"pull_request_id": "pr1", "pull_request_name": "Feature", "author_id": "u1"}`
	createResp := performRequest(testRouter, "POST", "/pullRequest/create", prBody)
	var createResponse map[string]interface{}
	json.Unmarshal(createResp.Body.Bytes(), &createResponse)
	pr := createResponse["pr"].(map[string]interface{})
	reviewers := pr["assigned_reviewers"].([]interface{})
	if len(reviewers) == 0 {
		t.Fatal("No reviewers assigned to PR")
	}
	oldReviewer := reviewers[0].(string)
	newReviewer := ""
	candidates := []string{"u2", "u3", "u4"}
	for _, c := range candidates {
		if c == oldReviewer {
			continue
		}
		isCurrentReviewer := false
		for _, r := range reviewers {
			if r.(string) == c {
				isCurrentReviewer = true
				break
			}
		}
		if !isCurrentReviewer {
			newReviewer = c
			break
		}
	}
	if newReviewer == "" {
		t.Fatal("Could not find new reviewer candidate")
	}
	reassignBody := fmt.Sprintf(`{"pull_request_id": "pr1", "old_reviewer_id": "%s", "new_reviewer_id": "%s"}`, oldReviewer, newReviewer)
	w := performRequest(testRouter, "POST", "/pullRequest/reassign", reassignBody)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}
	var exists bool
	err := testDB.Get(&exists, "SELECT EXISTS(SELECT 1 FROM pr_reviewers WHERE pull_request_id = 'pr1' AND reviewer_id = $1)", newReviewer)
	if err != nil {
		t.Fatalf("Failed to query reviewers: %v", err)
	}
	if !exists {
		t.Errorf("Expected %s to be a reviewer", newReviewer)
	}
}
