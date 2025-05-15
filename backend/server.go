package main

import (
	"log"
	"net/http"
	"social-nework/pkg/auth"

	"social-nework/pkg/db/sqlite"
	"social-nework/pkg/handlers"
)

func main() {
	// Initialize SQLite database
	db, err := sqlite.NewDB("social_network.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize the user model with database connection
	userModel := &auth.UserModel{DB: db}
authHandler := &handlers.AuthHandler{UserModel: userModel}
	http.HandleFunc("/api/register", authHandler.Register)
	http.HandleFunc("/api/login", authHandler.Login)
	http.HandleFunc("/api/logout", authHandler.Logout)
	// Set up a basic HTTP server
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Server is running"))
	})

	// Start the server
	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
