package handlers

import (
	"encoding/json"
	"net/http"
	"social-nework/pkg/models"
	"strconv"

	"github.com/gorilla/mux"
)


type Handler struct {
	db       *db.DB
	validate *validator.Validate
}

func (h *Handler) respondWithError(w http.ResponseWriter, code int, message string) {
	h.respondWithJSON(w, code, map[string]string{"error": message})
}

func (h *Handler) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (h *Handler) getCurrentUserID(r *http.Request) (int, bool) {
	// This would typically extract user ID from JWT token
	// For demo purposes, we'll get it from a header
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		return 0, false
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return 0, false
	}

	return userID, true
}

// Group handlers
func (h *Handler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.getCurrentUserID(r)
	if !ok {
		h.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req models.CreateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	group := &models.Group{
		Title:       req.Title,
		Description: req.Description,
		CreatorID:   userID,
	}

	if err := h.db.CreateGroup(group); err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to create group")
		return
	}

	// Add creator as a member
	if err := h.db.AddGroupMember(group.ID, userID); err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to add creator to group")
		return
	}

	h.respondWithJSON(w, http.StatusCreated, group)
}

func (h *Handler) GetAllGroups(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.getCurrentUserID(r)
	if !ok {
		h.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	limit := 20
	offset := 0

	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	groups, err := h.db.GetAllGroups(userID, limit, offset)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to get groups")
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"groups": groups,
		"limit":  limit,
		"offset": offset,
	})
}

func (h *Handler) GetGroup(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.getCurrentUserID(r)
	if !ok {
		h.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	groupID, err := strconv.Atoi(vars["groupId"])
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid group ID")
		return
	}

	group, err := h.db.GetGroupByID(groupID, userID)
	if err != nil {
		h.respondWithError(w, http.StatusNotFound, "Group not found")
		return
	}

	h.respondWithJSON(w, http.StatusOK, group)
}

// Invitation handlers
func (h *Handler) InviteUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.getCurrentUserID(r)
	if !ok {
		h.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	groupID, err := strconv.Atoi(vars["groupId"])
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid group ID")
		return
	}

	// Check if user is a member of the group
	isMember, err := h.db.IsGroupMember(groupID, userID)
	if err != nil || !isMember {
		h.respondWithError(w, http.StatusForbidden, "You must be a member to invite others")
		return
	}

	var req models.InviteUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Check if target user is already a member
	isAlreadyMember, err := h.db.IsGroupMember(groupID, req.UserID)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Database error")
		return
	}
	if isAlreadyMember {
		h.respondWithError(w, http.StatusBadRequest, "User is already a member")
		return
	}

	invitation := &models.Invitation{
		FromUserID: userID,
		ToUserID:   req.UserID,
		GroupID:    groupID,
		Type:       models.InvitationTypeInvitation,
		Status:     models.InvitationStatusPending,
	}

	if err := h.db.CreateInvitation(invitation); err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to create invitation")
		return
	}

	h.respondWithJSON(w, http.StatusCreated, invitation)
}

func (h *Handler) RequestToJoin(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.getCurrentUserID(r)
	if !ok {
		h.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	groupID, err := strconv.Atoi(vars["groupId"])
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid group ID")
		return
	}

	// Check if user is already a member
	isMember, err := h.db.IsGroupMember(groupID, userID)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Database error")
		return
	}
	if isMember {
		h.respondWithError(w, http.StatusBadRequest, "You are already a member")
		return
	}

	// Get group to find creator
	group, err := h.db.GetGroupByID(groupID, userID)
	if err != nil {
		h.respondWithError(w, http.StatusNotFound, "Group not found")
		return
	}

	invitation := &models.Invitation{
		FromUserID: userID,
		ToUserID:   group.CreatorID,
		GroupID:    groupID,
		Type:       models.InvitationTypeJoinRequest,
		Status:     models.InvitationStatusPending,
	}

	if err := h.db.CreateInvitation(invitation); err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to create join request")
		return
	}

	h.respondWithJSON(w, http.StatusCreated, invitation)
}

func (h *Handler) GetUserInvitations(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.getCurrentUserID(r)
	if !ok {
		h.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	invitations, err := h.db.GetUserInvitations(userID)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to get invitations")
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"invitations": invitations,
	})
}

func (h *Handler) GetGroupJoinRequests(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.getCurrentUserID(r)
	if !ok {
		h.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	groupID, err := strconv.Atoi(vars["groupId"])
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid group ID")
		return
	}

	// Check if user is the group creator
	group, err := h.db.GetGroupByID(groupID, userID)
	if err != nil {
		h.respondWithError(w, http.StatusNotFound, "Group not found")
		return
	}

	if group.CreatorID != userID {
		h.respondWithError(w, http.StatusForbidden, "Only group creator can view join requests")
		return
	}

	requests, err := h.db.GetGroupJoinRequests(groupID)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to get join requests")
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"requests": requests,
	})
}

func (h *Handler) HandleInvitation(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.getCurrentUserID(r)
	if !ok {
		h.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	invitationID, err := strconv.Atoi(vars["invitationId"])
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid invitation ID")
		return
	}

	var req models.HandleInvitationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Get invitation
	invitation, err := h.db.GetInvitationByID(invitationID)
	if err != nil {
		h.respondWithError(w, http.StatusNotFound, "Invitation not found")
		return
	}

	// Check authorization based on invitation type
	var canHandle bool
	if invitation.Type == models.InvitationTypeInvitation {
		canHandle = invitation.ToUserID == userID
	} else if invitation.Type == models.InvitationTypeJoinRequest {
		canHandle = invitation.ToUserID == userID
	}

	if !canHandle {
		h.respondWithError(w, http.StatusForbidden, "Not authorized to handle this invitation")
		return
	}

	// Update status
	var newStatus models.InvitationStatus
	if req.Action == "accept" {
		newStatus = models.InvitationStatusAccepted
	} else {
		newStatus = models.InvitationStatusDeclined
	}

	if err := h.db.UpdateInvitationStatus(invitationID, newStatus); err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to update invitation")
		return
	}

	// If accepted, add user to group
	if newStatus == models.InvitationStatusAccepted {
		var targetUserID int
		if invitation.Type == models.InvitationTypeInvitation {
			targetUserID = invitation.ToUserID
		} else {
			targetUserID = invitation.FromUserID
		}

		if err := h.db.AddGroupMember(invitation.GroupID, targetUserID); err != nil {
			h.respondWithError(w, http.StatusInternalServerError, "Failed to add user to group")
			return
		}
	}

	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Invitation " + req.Action + "ed successfully",
	})
}

// Post handlers
func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.getCurrentUserID(r)
	if !ok {
		h.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	groupID, err := strconv.Atoi(vars["groupId"])
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid group ID")
		return
	}

	// Check if user is a member
	isMember, err := h.db.IsGroupMember(groupID, userID)
	if err != nil || !isMember {
		h.respondWithError(w, http.StatusForbidden, "You must be a member to post")
		return
	}

	var req models.CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	post := &models.Post{
		GroupID:  groupID,
		AuthorID: userID,
		Title:    req.Title,
		Content:  req.Content,
	}

	if err := h.db.CreatePost(post); err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to create post")
		return
	}

	h.respondWithJSON(w, http.StatusCreated, post)
}

func (h *Handler) GetGroupPosts(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.getCurrentUserID(r)
	if !ok {
		h.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	groupID, err := strconv.Atoi(vars["groupId"])
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid group ID")
		return
	}

	// Check if user is a member
	isMember, err := h.db.IsGroupMember(groupID, userID)
	if err != nil || !isMember {
		h.respondWithError(w, http.StatusForbidden, "You must be a member to view posts")
		return
	}

	limit := 20
	offset := 0

	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	posts, err := h.db.GetGroupPosts(groupID, limit, offset)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to get posts")
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"posts":  posts,
		"limit":  limit,
		"offset": offset,
	})
}

// Comment handlers
func (h *Handler) CreateComment(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.getCurrentUserID(r)
	if !ok {
		h.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	postID, err := strconv.Atoi(vars["postId"])
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	// TODO: Verify user is member of the group that owns this post

	var req models.CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	comment := &models.Comment{
		PostID:   postID,
		AuthorID: userID,
		Content:  req.Content,
	}

	if err := h.db.CreateComment(comment); err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to create comment")
		return
	}

	h.respondWithJSON(w, http.StatusCreated, comment)
}

func (h *Handler) GetPostComments(w http.ResponseWriter, r *http.Request) {
	_, ok := h.getCurrentUserID(r)
	if !ok {
		h.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	postID, err := strconv.Atoi(vars["postId"])
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	// TODO: Verify user is member of the group that owns this post

	comments, err := h.db.GetPostComments(postID)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to get comments")
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"comments": comments,
	})
}

// Event handlers
func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.getCurrentUserID(r)
	if !ok {
		h.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	groupID, err := strconv.Atoi(vars["groupId"])
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid group ID")
		return
	}

	// Check if user is a member
	isMember, err := h.db.IsGroupMember(groupID, userID)
	if err != nil || !isMember {
		h.respondWithError(w, http.StatusForbidden, "You must be a member to create events")
		return
	}

	var req models.CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	event := &models.Event{
		GroupID:     groupID,
		CreatorID:   userID,
		Title:       req.Title,
		Description: req.Description,
		EventDate:   req.EventDate,
	}

	if err := h.db.CreateEvent(event); err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to create event")
		return
	}

	h.respondWithJSON(w, http.StatusCreated, event)
}

func (h *Handler) GetGroupEvents(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.getCurrentUserID(r)
	if !ok {
		h.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	groupID, err := strconv.Atoi(vars["groupId"])
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid group ID")
		return
	}

	// Check if user is a member
	isMember, err := h.db.IsGroupMember(groupID, userID)
	if err != nil || !isMember {
		h.respondWithError(w, http.StatusForbidden, "You must be a member to view events")
		return
	}

	events, err := h.db.GetGroupEvents(groupID, userID)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to get events")
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"events": events,
	})
}

func (h *Handler) RespondToEvent(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.getCurrentUserID(r)
	if !ok {
		h.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	eventID, err := strconv.Atoi(vars["eventId"])
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid event ID")
		return
	}

	var req models.EventResponseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// TODO: Verify user is member of the group that owns this event

	if err := h.db.CreateOrUpdateEventResponse(eventID, userID, req.Response); err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to record response")
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Response recorded successfully",
	})
}

func (h *Handler) GetEventResponses(w http.ResponseWriter, r *http.Request) {
	_, ok := h.getCurrentUserID(r)
	if !ok {
		h.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	eventID, err := strconv.Atoi(vars["eventId"])
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid event ID")
		return
	}

	// TODO: Verify user is member of the group that owns this event

	responses, err := h.db.GetEventResponses(eventID)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to get responses")
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"responses": responses,
	})
}
