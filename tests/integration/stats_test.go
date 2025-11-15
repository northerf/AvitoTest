package integration

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestGetStatistics(t *testing.T) {
	clearDB(t)
	teamBody := `{
        "team_name": "backend",
        "members": [
            {"user_id": "u1", "username": "Alice", "is_active": true},
            {"user_id": "u2", "username": "Bob", "is_active": false}
        ]
    }`
	performRequest(testRouter, "POST", "/team/add", teamBody)
	pr1Body := `{"pull_request_id": "pr1", "pull_request_name": "Feature", "author_id": "u1"}`
	performRequest(testRouter, "POST", "/pullRequest/create", pr1Body)
	mergeBody := `{"pull_request_id": "pr1"}`
	performRequest(testRouter, "POST", "/pullRequest/merge", mergeBody)
	w := performRequest(testRouter, "GET", "/stats/allstats", "")
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		return
	}
	var stats map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &stats)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v. Body: %s", err, w.Body.String())
	}
	if totalUsers, ok := stats["total_users"]; !ok {
		t.Errorf("Response missing 'total_users' field")
	} else if int(totalUsers.(float64)) != 2 {
		t.Errorf("Expected 2 total users, got %v", totalUsers)
	}
	if totalActiveUsers, ok := stats["total_active_users"]; !ok {
		t.Errorf("Response missing 'total_active_users' field")
	} else if int(totalActiveUsers.(float64)) != 1 {
		t.Errorf("Expected 1 active user, got %v", totalActiveUsers)
	}
	if totalPRs, ok := stats["total_prs"]; !ok {
		t.Errorf("Response missing 'total_prs' field")
	} else if int(totalPRs.(float64)) != 1 {
		t.Errorf("Expected 1 total PR, got %v", totalPRs)
	}
	if totalOpenPRs, ok := stats["total_open_prs"]; !ok {
		t.Errorf("Response missing 'total_open_prs' field")
	} else if int(totalOpenPRs.(float64)) != 0 {
		t.Errorf("Expected 0 open PRs, got %v", totalOpenPRs)
	}
	if totalMergedPRs, ok := stats["total_merged_prs"]; !ok {
		t.Errorf("Response missing 'total_merged_prs' field")
	} else if int(totalMergedPRs.(float64)) != 1 {
		t.Errorf("Expected 1 merged PR, got %v", totalMergedPRs)
	}
	if _, ok := stats["top_reviewers"]; !ok {
		t.Errorf("Response missing 'top_reviewers' field")
	}
	if prsWithoutReviewers, ok := stats["prs_without_reviewers"]; !ok {
		t.Errorf("Response missing 'prs_without_reviewers' field")
	} else if int(prsWithoutReviewers.(float64)) != 0 {
		t.Errorf("Expected 0 PRs without reviewers, got %v", prsWithoutReviewers)
	}
}

func TestGetUserStats(t *testing.T) {
	clearDB(t)
	teamBody := `{
        "team_name": "backend",
        "members": [
            {"user_id": "u1", "username": "Alice", "is_active": true},
            {"user_id": "u2", "username": "Bob", "is_active": true},
            {"user_id": "u3", "username": "Charlie", "is_active": true}
        ]
    }`
	performRequest(testRouter, "POST", "/team/add", teamBody)
	pr1Body := `{"pull_request_id": "pr1", "pull_request_name": "Feature", "author_id": "u1"}`
	performRequest(testRouter, "POST", "/pullRequest/create", pr1Body)
	w := performRequest(testRouter, "GET", "/stats/user?user_id=u2", "")
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		return
	}
	var stats map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &stats)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v. Body: %s", err, w.Body.String())
	}
	if userID, ok := stats["user_id"]; !ok {
		t.Errorf("Response missing 'user_id' field")
	} else if userID != "u2" {
		t.Errorf("Expected user_id 'u2', got %v", userID)
	}
	if username, ok := stats["username"]; !ok {
		t.Errorf("Response missing 'username' field")
	} else if username != "Bob" {
		t.Errorf("Expected username 'Bob', got %v", username)
	}
	if reviewsAssigned, ok := stats["reviews_assigned"]; !ok {
		t.Errorf("Response missing 'reviews_assigned' field")
	} else {
		assigned := int(reviewsAssigned.(float64))
		if assigned < 0 {
			t.Errorf("Expected reviews_assigned >= 0, got %d", assigned)
		}
		t.Logf("User u2 has %d reviews assigned", assigned)
	}
	if reviewsCompleted, ok := stats["reviews_completed"]; !ok {
		t.Errorf("Response missing 'reviews_completed' field")
	} else {
		completed := int(reviewsCompleted.(float64))
		if completed < 0 {
			t.Errorf("Expected reviews_completed >= 0, got %d", completed)
		}
	}
	t.Logf("User stats response: %+v", stats)
}
