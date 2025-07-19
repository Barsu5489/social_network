package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"social-nework/pkg/models"
	"social-nework/pkg/websocket"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type FollowHandler struct {
	FollowModel       *models.FollowModel
	NotificationModel *models.NotificationModel
	Hub               *websocket.Hub
}

func (h *FollowHandler) Follow(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	followerID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	followedID := vars["userID"]
	if followedID == "" {
		http.Error(w, "User ID required", http.StatusBadRequest)
		return
	}

	if followerID == followedID {
		http.Error(w, "Cannot follow yourself", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := h.FollowModel.Follow(ctx, followerID, followedID)
	if err != nil {
		if err.Error() == "cannot follow yourself" || err.Error() == "follow already exists" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("Failed to follow user: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Create a notification for the user being followed
	notification := models.Notification{
		ID:          uuid.New().String(),
		UserID:      followedID,
		Type:        "new_follower",
		ReferenceID: followerID,
		IsRead:      false,
		CreatedAt:   time.Now(),
	}

	_, err = h.NotificationModel.Insert(ctx, notification)
	if err != nil {
		log.Printf("Failed to create notification: %v", err)
	} else {
		// Send real-time notification if hub is available
		if h.Hub != nil {
			h.Hub.SendNotification(followedID, notification, map[string]interface{}{
				"follower_id": followerID,
				"action":      "followed",
			})
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully followed user"))
}

func (h *FollowHandler) Unfollow(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	followedID := vars["userID"]
	if followedID == "" {
		http.Error(w, "User ID required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := h.FollowModel.Unfollow(ctx, userID, followedID)
	if err != nil {
		if err.Error() == "follow does not exist" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("Failed to unfollow user: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Unfollowed successfully"})
}

func (h *FollowHandler) GetFollowers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value("user_id").(string)
	fmt.Println(userID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	fmt.Println(ctx)
	defer cancel()

	followers, err := h.FollowModel.GetFollowers(ctx, userID)
	fmt.Println(followers)
	if err != nil {
		log.Printf("Failed to get followers: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(followers)
}

func (h *FollowHandler) GetFollowing(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	following, err := h.FollowModel.GetFollowing(ctx, userID)
	fmt.Println(following)
	if err != nil {
		log.Printf("Failed to get following: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(following)
}

// CheckFollowStatus checks if the current user is following a specific user and vice versa
func (h *FollowHandler) CheckFollowStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	targetUserID := r.URL.Query().Get("targetUserId")
	if targetUserID == "" {
		http.Error(w, "targetUserId parameter is required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	isFollowing, isFollowedBy, err := h.FollowModel.CheckFollowStatus(ctx, userID, targetUserID)
	if err != nil {
		log.Printf("Failed to check follow status: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{
		"isFollowing":  isFollowing,
		"isFollowedBy": isFollowedBy,
	})
}
