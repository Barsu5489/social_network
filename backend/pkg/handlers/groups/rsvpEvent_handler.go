package groups

import (
	"encoding/json"
	"net/http"
	"social-nework/pkg/models"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// RSVP to event
func (gh *GroupHandler) RSVPEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventID := vars["eventId"]

	var attendee models.EventAttendee
	if err := json.NewDecoder(r.Body).Decode(&attendee); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Check if user is member of the group that owns the event
	memberQuery := `SELECT gm.id FROM group_members gm 
					INNER JOIN events e ON gm.group_id = e.group_id 
					WHERE e.id = ? AND gm.user_id = ? AND gm.deleted_at IS NULL`
	var memberExists string
	err := gh.db.QueryRow(memberQuery, eventID, attendee.UserID).Scan(&memberExists)
	if err != nil {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	attendee.ID = uuid.New().String()
	attendee.EventID = eventID
	attendee.CreatedAt = time.Now().Unix()

	// Use UPSERT to handle updating existing RSVPs
	query := `INSERT INTO event_attendees (id, event_id, user_id, status, created_at) 
			  VALUES (?, ?, ?, ?, ?) 
			  ON CONFLICT(event_id, user_id) 
			  DO UPDATE SET status = excluded.status`

	_, err = gh.db.Exec(query, attendee.ID, attendee.EventID, attendee.UserID,
		attendee.Status, attendee.CreatedAt)
	if err != nil {
		http.Error(w, "Failed to RSVP", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(attendee)
}
