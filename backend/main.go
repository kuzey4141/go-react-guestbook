package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func openDB(connStr string) (*sql.DB, error) {
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

// enableCORS adds the headers required for the React app on localhost:3000.
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	connStr := "postgres://kuzey:yourpassword@localhost:5432/guestbook_db?sslmode=disable"
	db, err := openDB(connStr)
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}
	fmt.Println("Connected to database successfully.")

	store := NewMessageStore(db)

	app := &Application{
		store: store,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/messages", app.handleMessages)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: enableCORS(mux),
	}

	fmt.Println("Starting API server on http://localhost:8080...")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
