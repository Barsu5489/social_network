package groups

import (
	"encoding/json"
	"net/http"
	"social-nework/pkg/models"
	"github.com/gorilla/mux"
)

// Get group events with RSVP details
func (gh *GroupHandler) GetGroupEvents(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	groupID := vars["groupId"]

	// Check if user is member of the group
	memberQuery := `SELECT id FROM group_members WHERE group_id = ? AND user_id = ? AND deleted_at IS NULL`
	var memberExists string
	err := gh.db.QueryRow(memberQuery, groupID, userID).Scan(&memberExists)
	if err != nil {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	// Get events with attendee count
	eventQuery := `
		SELECT 
			e.id, 
			e.group_id, 
			e.title, 
			e.description, 
			e.location, 
			e.start_time, 
			e.end_time, 
			e.created_by, 
			e.created_at, 
			e.updated_at,
			COALESCE(COUNT(ea.id), 0) as attendee_count
		FROM events e
		LEFT JOIN event_attendees ea ON e.id = ea.event_id AND ea.deleted_at IS NULL
		WHERE e.group_id = ? AND e.deleted_at IS NULL 
		GROUP BY e.id
		ORDER BY e.start_time ASC`

	rows, err := gh.db.Query(eventQuery, groupID)
	if err != nil {
		http.Error(w, "Failed to fetch events", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var events []models.Event
	var eventIDs []string

	// First, collect all events and their IDs
	for rows.Next() {
		var event models.Event
		err := rows.Scan(
			&event.ID, 
			&event.GroupID, 
			&event.Title, 
			&event.Description,
			&event.Location, 
			&event.StartTime, 
			&event.EndTime, 
			&event.CreatedBy,
			&event.CreatedAt, 
			&event.UpdatedAt,
			&event.AttendeeCount,
		)
		if err != nil {
			continue
		}
		events = append(events, event)
		eventIDs = append(eventIDs, event.ID)
	}

	// If we have events, get all RSVP details for these events
	if len(eventIDs) > 0 {
		// Create placeholders for IN clause
		placeholders := ""
		args := []interface{}{}
		for i, eventID := range eventIDs {
			if i > 0 {
				placeholders += ","
			}
			placeholders += "?"
			args = append(args, eventID)
		}

		// Get attendee details for all events
		attendeeQuery := `
			SELECT 
				ea.id,
				ea.event_id,
				ea.user_id,
				COALESCE(u.first_name || ' ' || u.last_name, u.nickname, 'Unknown User') as user_name,
				ea.status,
				ea.created_at as joined_at
			FROM event_attendees ea
			JOIN users u ON ea.user_id = u.id
			WHERE ea.event_id IN (` + placeholders + `) AND ea.deleted_at IS NULL
			ORDER BY ea.created_at ASC`

		attendeeRows, err := gh.db.Query(attendeeQuery, args...)
		if err != nil {
			http.Error(w, "Failed to fetch attendee details", http.StatusInternalServerError)
			return
		}
		defer attendeeRows.Close()

		// Group attendees by event ID
		eventAttendees := make(map[string][]models.EventAttendee)
		for attendeeRows.Next() {
			var attendee models.EventAttendee
			var eventID string
			err := attendeeRows.Scan(
				&attendee.ID,
				&eventID,
				&attendee.UserID,
				&attendee.UserName,
				&attendee.Status,
				&attendee.JoinedAt,
			)
			if err != nil {
				continue
			}
			eventAttendees[eventID] = append(eventAttendees[eventID], attendee)
		}

		// Attach attendee details to events
		for i := range events {
			if attendees, exists := eventAttendees[events[i].ID]; exists {
				events[i].Attendees = attendees
			} else {
				events[i].Attendees = []models.EventAttendee{} // Empty slice instead of nil
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}