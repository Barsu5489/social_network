package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"social-nework/pkg/auth"
	"social-nework/pkg/db/sqlite"
	"social-nework/pkg/handlers"
	"social-nework/pkg/handlers/groups"
	"social-nework/pkg/models"
	"social-nework/pkg/repository"
	"social-nework/pkg/websocket"
)

func setupChatSystem(db *sql.DB, router *mux.Router) (*websocket.Hub, *repository.ChatRepository, *repository.GroupRepository) {
	// Initialize repositories
	chatRepo := &repository.ChatRepository{DB: db}
	messageRepo := &repository.MessageRepository{DB: db}
	groupRepo := &repository.GroupRepository{DB: db} // Add group repository

	// Initialize WebSocket hub
	hub := websocket.NewHub(db, messageRepo, chatRepo)
	go hub.Run() // Start the hub in a goroutine
	// Initialize HTTP handlers with all required repositories
	chatHandler := handlers.NewChatHandler(chatRepo, messageRepo, groupRepo, hub)

	// Register chat routes
	registerChatRoutes(router, chatHandler)

	// Register WebSocket endpoint
	router.HandleFunc("/ws", websocket.WebSocketAuth(hub, func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWS(hub, w, r)
	})).Methods("GET")

	// Return the instances so they can be used elsewhere
	return hub, chatRepo, groupRepo
}

func registerChatRoutes(router *mux.Router, handler *handlers.ChatHandler) {
	// Chat management routes
	router.HandleFunc("/api/chats", auth.RequireAuth(handler.GetUserChats)).Methods("GET")
	router.HandleFunc("/api/chats/direct", auth.RequireAuth(handler.CreateDirectChat)).Methods("POST")
	router.HandleFunc("/api/chats/group", auth.RequireAuth(handler.CreateGroupChat)).Methods("POST")

	// Message management routes
	router.HandleFunc("/api/chats/{chatId}/messages", auth.RequireAuth(handler.GetChatMessages)).Methods("GET")
	router.HandleFunc("/api/chats/{chatId}/messages", auth.RequireAuth(handler.SendMessage)).Methods("POST")
	router.HandleFunc("/api/chats/{chatId}/participants", auth.RequireAuth(handler.AddParticipant)).Methods("POST")

	// Group chat helper route
	router.HandleFunc("/api/groups/{groupId}/chat", auth.RequireAuth(handler.GetGroupChatForGroup)).Methods("GET")
}

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

	// Initialize router
	router := mux.NewRouter()

	// Setup chat system with all routes and get the required instances
	hub, chatRepo, groupRepo := setupChatSystem(db, router)

	// Now create the GroupHandler with all required dependencies
	groupHandler := groups.NewGroupHandler(db, groupRepo, chatRepo, hub)

	// Handlers
	authHandler := &handlers.AuthHandler{UserModel: userModel}
	followHandler := &handlers.FollowHandler{FollowModel: followModel}

	// Follow routes with middleware RequireAuth
	router.HandleFunc("/follow/{userID}", auth.RequireAuth(followHandler.Follow)).Methods("POST")
	router.HandleFunc("/unfollow/{userID}", auth.RequireAuth(followHandler.Unfollow)).Methods("DELETE")
	router.HandleFunc("/followers", auth.RequireAuth(followHandler.GetFollowers)).Methods("GET")
	router.HandleFunc("/following", auth.RequireAuth(followHandler.GetFollowing)).Methods("GET")

	// Public routes without middleware
	router.HandleFunc("/api/register", authHandler.Register).Methods("POST")
	router.HandleFunc("/api/login", authHandler.Login).Methods("POST")
	router.HandleFunc("/api/logout", authHandler.Logout).Methods("POST")

	// Post routes with middleware
	router.HandleFunc("/post", auth.RequireAuth(handlers.NewPost(db))).Methods("POST")
	router.HandleFunc("/followPosts", auth.RequireAuth(handlers.FollowingPosts(db))).Methods("GET")
	router.HandleFunc("/delPost/{post_id}", auth.RequireAuth(handlers.DeletPost(db))).Methods("DELETE")
	router.HandleFunc("/posts", auth.RequireAuth(handlers.AllPosts(db))).Methods("GET")

	// Comment routes
	router.HandleFunc("/comment/{post_id}", auth.RequireAuth(handlers.NewComment(db))).Methods("POST")
	router.HandleFunc("/comments/{post_id}", auth.RequireAuth(handlers.GetPostComments(db))).Methods("GET")

	// User Profile routes
	router.HandleFunc("/api/profile", auth.RequireAuth(handlers.GetProfile(db))).Methods("GET")
	router.HandleFunc("/api/profile", auth.RequireAuth(handlers.UpdateProfile(db))).Methods("PUT")

	// Like routes
	router.HandleFunc("/posts/{post_id}/like", auth.RequireAuth(handlers.LikePost(db))).Methods(http.MethodPost)
	router.HandleFunc("/posts/{post_id}/like", auth.RequireAuth(handlers.LikePost(db))).Methods(http.MethodDelete)
	router.HandleFunc("/posts/{post_id}/likes", auth.RequireAuth(handlers.GetPostLikes(db))).Methods(http.MethodGet)
	router.HandleFunc("/users/{user_id}/likes", auth.RequireAuth(handlers.GetUserLikedPosts(db))).Methods(http.MethodGet)

	// Group Management Routes
	router.HandleFunc("/api/groups", auth.RequireAuth(groupHandler.CreateGroup)).Methods("POST")
	router.HandleFunc("/api/groups", auth.RequireAuth(groupHandler.GetAllGroups)).Methods("GET")

	// Group Invitation Routes
	router.HandleFunc("/api/groups/invite", auth.RequireAuth(groupHandler.InviteToGroup)).Methods("POST")
	router.HandleFunc("/api/groups/join/{groupId}", auth.RequireAuth(groupHandler.RequestToJoinGroup)).Methods("POST")
	router.HandleFunc("/api/invitations/{id}/respond", auth.RequireAuth(groupHandler.RespondToInvitation)).Methods("PUT")

	// Group Content Routes
	router.HandleFunc("/api/groups/{groupId}/posts", auth.RequireAuth(groupHandler.CreateGroupPost)).Methods("POST")
	router.HandleFunc("/api/groups/{groupId}/posts", auth.RequireAuth(groupHandler.GetGroupPosts)).Methods("GET")

	// Group Event Routes
	router.HandleFunc("/api/groups/{groupId}/events", auth.RequireAuth(groupHandler.CreateEvent)).Methods("POST")
	router.HandleFunc("/api/groups/{groupId}/events", auth.RequireAuth(groupHandler.GetGroupEvents)).Methods("GET")
	router.HandleFunc("/api/events/{eventId}/rsvp", auth.RequireAuth(groupHandler.RSVPEvent)).Methods("POST")

	// CORS MIDDLEWARE
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://127.0.0.1:5500"}, // frontend origin
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// Wrap the router with CORS middleware
	handler := corsHandler.Handler(router)

	// Start server
	if err := http.ListenAndServe(":3000", handler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}