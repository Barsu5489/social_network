package groups

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"social-nework/pkg/models"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Create event
func (gh *GroupHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	groupID := vars["groupId"]

	var event models.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	event.CreatedBy = userID

	// Check if user is member of the group
	memberQuery := `SELECT id FROM group_members WHERE group_id = ? AND user_id = ? AND deleted_at IS NULL`
	var memberExists string
	err := gh.db.QueryRow(memberQuery, groupID, event.CreatedBy).Scan(&memberExists)
	if err != nil {
		http.Error(w, "You must be a member to create events", http.StatusForbidden)
		return
	}

	event.ID = uuid.New().String()
	event.GroupID = groupID
	event.CreatedAt = time.Now().Unix()
	event.UpdatedAt = time.Now().Unix()

	query := `INSERT INTO events (id, group_id, title, description, location, start_time, end_time, created_by, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err = gh.db.Exec(query, event.ID, event.GroupID, event.Title, event.Description,
		event.Location, event.StartTime, event.EndTime, event.CreatedBy, event.CreatedAt, event.UpdatedAt)
	if err != nil {
		http.Error(w, "Failed to create event", http.StatusInternalServerError)
		return
	}

	// Notify all group members
	membersQuery := `SELECT user_id FROM group_members WHERE group_id = ? AND user_id != ? AND deleted_at IS NULL`
	rows, err := gh.db.Query(membersQuery, groupID, event.CreatedBy)
	if err == nil {
		defer rows.Close()
		memberCount := 0
		for rows.Next() {
			var memberUserID string
			if rows.Scan(&memberUserID) == nil {
				memberCount++
				notification := models.Notification{
					UserID:      memberUserID,
					Type:        "event_created",
					ReferenceID: event.ID,
					IsRead:      false,
					CreatedAt:   time.Now(),
				}

				log.Printf("DEBUG: Creating event notification for member %s - EventID: %s, GroupID: %s",
					memberUserID, event.ID, groupID)

				_, err = gh.NotificationModel.Insert(r.Context(), notification)
				if err != nil {
					log.Printf("ERROR: Failed to create event notification for user %s: %v", memberUserID, err)
				} else {
					log.Printf("SUCCESS: Event notification created for member %s", memberUserID)

					// Send real-time notification
					if gh.h != nil {
						// Get creator info
						var creatorNickname, creatorAvatar string
						gh.db.QueryRow("SELECT nickname, avatar_url FROM users WHERE id = ?", event.CreatedBy).Scan(&creatorNickname, &creatorAvatar)
						
						gh.h.SendNotification(memberUserID, notification, map[string]interface{}{
							"event_id":       event.ID,
							"group_id":       groupID,
							"creator_id":     event.CreatedBy,
							"actor_nickname": creatorNickname,
							"actor_avatar":   creatorAvatar,
						})
					}
				}
			}
		}
		log.Printf("DEBUG: Event notifications sent to %d group members", memberCount)
	} else {
		log.Printf("ERROR: Failed to get group members for event notification: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}
