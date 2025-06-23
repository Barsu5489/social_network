package groups

import (
	"encoding/json"
	"net/http"
	"social-nework/pkg/models"

	"github.com/gorilla/mux"
)

// Get group events
func (gh *GroupHandler) GetGroupEvents(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupID := vars["groupId"]
	userID := r.URL.Query().Get("user_id")

	// Check if user is member of the group
	memberQuery := `SELECT id FROM group_members WHERE group_id = ? AND user_id = ? AND deleted_at IS NULL`
	var memberExists string
	err := gh.db.QueryRow(memberQuery, groupID, userID).Scan(&memberExists)
	if err != nil {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	query := `SELECT id, group_id, title, description, location, start_time, end_time, created_by, created_at, updated_at 
			  FROM events WHERE group_id = ? AND deleted_at IS NULL ORDER BY start_time ASC`

	rows, err := gh.db.Query(query, groupID)
	if err != nil {
		http.Error(w, "Failed to fetch events", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var event models.Event
		err := rows.Scan(&event.ID, &event.GroupID, &event.Title, &event.Description,
			&event.Location, &event.StartTime, &event.EndTime, &event.CreatedBy,
			&event.CreatedAt, &event.UpdatedAt)
		if err != nil {
			continue
		}
		events = append(events, event)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}
