package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"social-nework/pkg/auth"
	"social-nework/pkg/db/sqlite"
	"social-nework/pkg/handlers"
	"social-nework/pkg/models"
)

func main() {
	// Initialize SQLite database
	db, err := sqlite.NewDB("social_network.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Models
	userModel := &auth.UserModel{DB: db}
	followModel := &models.FollowModel{DB: db}
	// postModel := &models.PostModel{DB: db}

	// Handlers
	authHandler := &handlers.AuthHandler{UserModel: userModel}
	followHandler := &handlers.FollowHandler{FollowModel: followModel}
	// postHandler := &handlers.PostHandler{Post: postModel}

	//  Initialize router
	router := mux.NewRouter()

	// follow routes with middleware RequireAuth
	router.HandleFunc("/follow/{userID}", auth.RequireAuth(followHandler.Follow)).Methods("POST")
	router.HandleFunc("/unfollow/{userID}", auth.RequireAuth(followHandler.Unfollow)).Methods("DELETE")
	router.HandleFunc("/followers", auth.RequireAuth(followHandler.GetFollowers)).Methods("GET")
	router.HandleFunc("/following", auth.RequireAuth(followHandler.GetFollowing)).Methods("GET")

	// Public routes without middleware
	router.HandleFunc("/api/register", authHandler.Register).Methods("POST")
	router.HandleFunc("/api/login", authHandler.Login).Methods("POST")
	router.HandleFunc("/api/logout", authHandler.Logout).Methods("POST")

	// post routes with middleware
	router.HandleFunc("/post", auth.RequireAuth(handlers.NewPost(db))).Methods("POST")
	router.HandleFunc("/followPosts", auth.RequireAuth(handlers.FollowingPosts(db))).Methods("GET")
	router.HandleFunc("/delPost/{post_id}", auth.RequireAuth(handlers.DeletPost(db))).Methods("DELETE")
	router.HandleFunc("/posts", auth.RequireAuth(handlers.AllPosts(db))).Methods("GET")


	// User Profile routes
	router.HandleFunc("/api/profile", auth.RequireAuth(handlers.GetProfile(db))).Methods("GET")
	router.HandleFunc("/api/profile", auth.RequireAuth(handlers.UpdateProfile(db))).Methods("PUT")

// todo - fix likes models and handlers
	// Like a post
	router.HandleFunc("/posts/{post_id}/like",auth.RequireAuth( handlers.LikePost(db))).Methods(http.MethodPost)

	// Unlike a post
	router.HandleFunc("/posts/{post_id}/like", auth.RequireAuth(handlers.LikePost(db))).Methods(http.MethodDelete)

	// Get likes for a post
	router.HandleFunc("/posts/{post_id}/likes", auth.RequireAuth(handlers.GetPostLikes(db))).Methods(http.MethodGet)

	// Get posts liked by a user
	router.HandleFunc("/users/{user_id}/likes", auth.RequireAuth(handlers.GetUserLikedPosts(db))).Methods(http.MethodGet)

	// Optional: To get liked posts by currently logged-in user
	// router.HandleFunc("/me/likes", auth.RequireAuth(handlers.GetUserLikedPosts(db))).Methods(http.MethodGet)


	// Start server
	http.ListenAndServe(":3000", router)

	// Test route
	router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Server is running"))
	})

	// Start server with router
	log.Println("Server starting on :3000...")
	if err := http.ListenAndServe(":3000", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
