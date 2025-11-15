package integration

import (
	"Avito/internal/handler"
	"Avito/internal/repository"
	"Avito/internal/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"io"
	"log"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var (
	testDB     *sqlx.DB
	testRouter *gin.Engine
)

func TestMain(m *testing.M) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))
	var err error
	testDB, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to test DB: %v", err)
	}
	defer testDB.Close()
	repo := repository.NewRepository(testDB)
	services := service.NewService(repo)
	handlers := handler.NewHandler(services)
	testRouter = handlers.InitRoutes()
	code := m.Run()
	os.Exit(code)
}

func clearDB(t *testing.T) {
	_, err := testDB.Exec(`TRUNCATE TABLE pr_reviewers CASCADE; TRUNCATE TABLE pull_requests CASCADE;
        TRUNCATE TABLE team_members CASCADE; TRUNCATE TABLE users CASCADE; TRUNCATE TABLE teams CASCADE;`)
	if err != nil {
		t.Fatalf("Failed to clear DB: %v", err)
	}
}

func performRequest(r *gin.Engine, method, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
		req.Body = io.NopCloser(strings.NewReader(body))
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
