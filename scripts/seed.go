package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func main() {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()
	log.Println("Connected to DB")
	tx, err := db.Begin()
	if err != nil {
		log.Fatal("Failed to begin transaction:", err)
	}
	defer tx.Rollback()
	teamName := "backend"
	_, err = tx.Exec(`INSERT INTO teams (team_name) VALUES ($1) ON CONFLICT (team_name) DO NOTHING`, teamName)
	if err != nil {
		log.Fatal("Failed to create team:", err)
	}
	users := []struct {
		UserID   string
		Username string
		IsActive bool
	}{
		{"user1", "Alice", true},
		{"user2", "Bob", true},
		{"user3", "Charlie", true},
		{"user4", "David", true},
		{"user5", "Eve", false},
	}

	for _, user := range users {
		_, err = tx.Exec(`
            INSERT INTO users (user_id, username, team_name, is_active)
            VALUES ($1, $2, $3, $4)
            ON CONFLICT (user_id) DO UPDATE 
            SET username = EXCLUDED.username,
                team_name = EXCLUDED.team_name,
                is_active = EXCLUDED.is_active
        `, user.UserID, user.Username, teamName, user.IsActive)
		if err != nil {
			log.Fatal("Failed to create user:", err)
		}
		_, err = tx.Exec(`INSERT INTO team_members (team_name, user_id) VALUES ($1, $2)
            ON CONFLICT (team_name, user_id) DO NOTHING`, teamName, user.UserID)
		if err != nil {
			log.Fatal("Failed to add user to team:", err)
		}
	}
	_, err = tx.Exec(`INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, status) VALUES 
            ('pr1', 'Add login feature', 'user1', 'OPEN'),
            ('pr2', 'Fix bug #123', 'user2', 'OPEN'),
            ('pr3', 'Update docs', 'user3', 'MERGED') ON CONFLICT (pull_request_id) DO NOTHING`)
	if err != nil {
		log.Fatal("Failed to create PRs:", err)
	}
	_, err = tx.Exec(`
        INSERT INTO pr_reviewers (pull_request_id, reviewer_id) VALUES 
            ('pr1', 'user2'),
            ('pr1', 'user3'),
            ('pr2', 'user3'),
            ('pr2', 'user4'),
            ('pr3', 'user4')ON CONFLICT (pull_request_id, reviewer_id) DO NOTHING`)
	if err != nil {
		log.Fatal("Failed to assign reviewers:", err)
	}
	if err := tx.Commit(); err != nil {
		log.Fatal("Failed to commit transaction:", err)
	}
}
