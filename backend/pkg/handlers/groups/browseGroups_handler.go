package groups

import (
    "encoding/json"
    "net/http"
    "social-nework/pkg/models"
)

// Browse all public groups
func (gh *GroupHandler) BrowseGroups(w http.ResponseWriter, r *http.Request) {
    userID, ok := r.Context().Value("user_id").(string)
    if !ok || userID == "" {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    query := `
        SELECT g.id, g.name, g.description, g.creator_id, g.is_private, g.created_at,
               COUNT(gm.id) as member_count
        FROM groups g
        LEFT JOIN group_members gm ON g.id = gm.group_id AND gm.deleted_at IS NULL
        WHERE g.deleted_at IS NULL AND g.is_private = false
        GROUP BY g.id
        ORDER BY g.created_at DESC`

    rows, err := gh.db.Query(query)
    if err != nil {
        http.Error(w, "Failed to fetch groups", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var groups []models.Group
    for rows.Next() {
        var group models.Group
        err := rows.Scan(&group.ID, &group.Name, &group.Description, 
            &group.CreatorID, &group.IsPrivate, &group.CreatedAt, &group.MemberCount)
        if err != nil {
            continue
        }
        groups = append(groups, group)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(groups)
}