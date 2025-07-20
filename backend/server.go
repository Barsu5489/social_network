package main

import (
	"database/sql"
	"log"
	"net/http"
	"strings"

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

func setupChatSystem(db *sql.DB, router *mux.Router, notificationModel *models.NotificationModel) (*websocket.Hub, *repository.ChatRepository, *repository.GroupRepository) {
	// Initialize repositories
	chatRepo := &repository.ChatRepository{DB: db}
	messageRepo := &repository.MessageRepository{DB: db}
	groupRepo := &repository.GroupRepository{DB: db} // Add group repository

	// Initialize WebSocket hub
	hub := websocket.NewHub(db, messageRepo, chatRepo)
	hub.SetNotificationModel(notificationModel) // Set the notification model
	go hub.Run() // Start the hub in a goroutine
	// Initialize HTTP handlers with all required repositories
	chatHandler := handlers.NewChatHandler(chatRepo, messageRepo, groupRepo, hub, notificationModel)

	// Register chat routes
	registerChatRoutes(router, chatHandler)

	// WebSocket endpoint with authentication
	router.HandleFunc("/ws", auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("WebSocket connection attempt from %s", r.RemoteAddr)
		websocket.ServeWS(hub, w, r)
	})).Methods("GET")

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

func registerGroupRoutes(router *mux.Router, handler *groups.GroupHandler) {
	// Group management routes
	router.HandleFunc("/api/groups", auth.RequireAuth(handler.GetAllGroups)).Methods("GET")
	router.HandleFunc("/api/groups", auth.RequireAuth(handler.CreateGroup)).Methods("POST")
	router.HandleFunc("/api/groups/{groupId}/join", auth.RequireAuth(handler.JoinGroup)).Methods("POST")
	router.HandleFunc("/api/groups/{groupId}/leave", auth.RequireAuth(handler.LeaveGroup)).Methods("POST")
	router.HandleFunc("/api/groups/join/{groupId}", auth.RequireAuth(handler.RequestToJoinGroup)).Methods("POST")

	// Group content routes
	router.HandleFunc("/api/groups/{groupId}/posts", auth.RequireAuth(handler.GetGroupPosts)).Methods("GET")
	router.HandleFunc("/api/groups/{groupId}/posts", auth.RequireAuth(handler.CreateGroupPost)).Methods("POST")
	router.HandleFunc("/api/groups/{groupId}/events", auth.RequireAuth(handler.GetGroupEvents)).Methods("GET")
	router.HandleFunc("/api/groups/{groupId}/events", auth.RequireAuth(handler.CreateEvent)).Methods("POST")
	router.HandleFunc("/api/groups/{groupId}/events/{eventId}/rsvp", auth.RequireAuth(handler.RSVPEvent)).Methods("POST")

	// Group invitation routes
	router.HandleFunc("/api/groups/invite", auth.RequireAuth(handler.InviteToGroup)).Methods("POST")
	router.HandleFunc("/api/invitations/{id}/respond", auth.RequireAuth(handler.RespondToInvitation)).Methods("POST")

	// Group chat route
	router.HandleFunc("/api/groups/{groupId}/chat", auth.RequireAuth(handler.GetGroupChat)).Methods("GET")
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
	notificationModel := &models.NotificationModel{DB: db}

	// Initialize router
	router := mux.NewRouter()

	// Setup chat system with all routes and get the required instances
	hub, chatRepo, groupRepo := setupChatSystem(db, router, notificationModel)

	// Handlers with hub for real-time notifications
	authHandler := &handlers.AuthHandler{UserModel: userModel}
	followHandler := &handlers.FollowHandler{
		FollowModel:       followModel,
		NotificationModel: notificationModel,
		Hub:               hub,
		DB:                db,
	}
	notificationHandler := handlers.NewNotificationHandler(notificationModel)

	// Now create the GroupHandler with all required dependencies
	groupHandler := groups.NewGroupHandler(db, groupRepo, chatRepo, hub, notificationModel)

	// Auth routes
	router.HandleFunc("/api/register", authHandler.Register).Methods("POST")
	router.HandleFunc("/api/login", authHandler.Login).Methods("POST")
	router.HandleFunc("/api/logout", authHandler.Logout).Methods("POST")
	router.HandleFunc("/api/profile", auth.RequireAuth(handlers.GetProfile(db))).Methods("GET")
	router.HandleFunc("/api/profile", auth.RequireAuth(handlers.UpdateProfile(db))).Methods("PUT")

	// Follow routes
	router.HandleFunc("/api/users/{userID}/follow", auth.RequireAuth(followHandler.Follow)).Methods("POST")
	router.HandleFunc("/api/users/{userID}/unfollow", auth.RequireAuth(followHandler.Unfollow)).Methods("DELETE")
	router.HandleFunc("/api/follow/check", auth.RequireAuth(followHandler.CheckFollowStatus)).Methods("GET")

	// Notification routes
	router.HandleFunc("/api/notifications", auth.RequireAuth(notificationHandler.GetNotifications)).Methods("GET")
	router.HandleFunc("/api/notifications/mark-read", auth.RequireAuth(notificationHandler.MarkNotificationAsRead)).Methods("POST")

	// Group routes
	registerGroupRoutes(router, groupHandler)

	// Posts routes
	router.HandleFunc("/api/posts", auth.RequireAuth(handlers.AllPosts(db))).Methods("GET")
	router.HandleFunc("/api/posts", auth.RequireAuth(handlers.NewPost(db))).Methods("POST")
	router.HandleFunc("/api/posts/{post_id}", auth.RequireAuth(handlers.GetSinglePost(db))).Methods("GET")
	router.HandleFunc("/api/posts/{post_id}", auth.RequireAuth(handlers.DeletPost(db))).Methods("DELETE")

	// Comment routes - temporarily remove auth from GET to debug
	router.HandleFunc("/comments/{postId}", handlers.GetPostComments(db)).Methods("GET")
	router.HandleFunc("/comment/{postId}", auth.RequireAuth(handlers.CreateComment(db, notificationModel, hub))).Methods("POST")

	// Add debugging middleware for comment routes
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.URL.Path, "comment") {
				log.Printf("DEBUG: Comment route accessed - Method: %s, Path: %s", r.Method, r.URL.Path)
			}
			next.ServeHTTP(w, r)
		})
	})

	// Like routes
	router.HandleFunc("/api/posts/{post_id}/like", auth.RequireAuth(handlers.LikePost(db, notificationModel, hub))).Methods("POST")
	router.HandleFunc("/api/posts/{post_id}/like", auth.RequireAuth(handlers.LikePost(db, notificationModel, hub))).Methods("DELETE")
	router.HandleFunc("/api/posts/{post_id}/likes", auth.RequireAuth(handlers.GetPostLikes(db))).Methods("GET")
	router.HandleFunc("/api/comments/{comment_id}/like", auth.RequireAuth(handlers.LikeComment(db, notificationModel, hub))).Methods("POST")
	router.HandleFunc("/api/comments/{comment_id}/like", auth.RequireAuth(handlers.LikeComment(db, notificationModel, hub))).Methods("DELETE")
	router.HandleFunc("/api/comments/{comment_id}/likes", auth.RequireAuth(handlers.GetCommentLikes(db))).Methods("GET")

	// Enable CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		Debug:            true, // Add debug logging
	})

	handler := c.Handler(router)

	log.Println("Server starting on :3000")
	log.Fatal(http.ListenAndServe(":3000", handler))
}
