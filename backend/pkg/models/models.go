package models

import "time"

type User struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"password"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Nickname     string `json:"nickname"`
	DateOfBirth  string `json:"date_of_birth"`
	AboutMe      string `json:"about_me"`
	AvatarURL    string `json:"avatar_url"`
	IsPrivate    bool   `json:"is_private"`
	CreatedAt    int64  `json:"created_at"`
	UpdatedAt    int64  `json:"updated_at"`
	DeletedAt    *int64 `json:"deleted_at"`
}

type Follow struct {
	ID         string `json:"id"`
	FollowerID string `json:"follower_id"`
	FollowedID string `json:"followed_id"`
	Status     string `json:"status"`
	CreatedAt  int64  `json:"created_at"`
	DeletedAt  *int64 `json:"deleted_at,omitempty"`
}

type Post struct {
	ID        string  `json:"id"`
	UserID    string  `json:"user_id"`
	GroupID   *string `json:"group_id"` // Nullable
	Content   string  `json:"content"`
	Privacy   string  `json:"privacy"`
	CreatedAt int64   `json:"created_at"`

	UpdatedAt  int64  `json:"updated_at"`
	DeletedAt  *int64 `json:"deleted_at"` // Nullable
	LikesCount int    `json:"likes_count"`
	UserLiked  bool   `json:"user_liked,omitempty"`
	AuthorID     int       `json:"author_id" db:"author_id"`
    Title        string    `json:"title" db:"title"`
    Author       *User     `json:"author,omitempty"`
    CommentCount int       `json:"comment_count,omitempty"`
}

// Like represents a like on a post or comment
type Like struct {
	ID           string `json:"id"`
	UserID       string `json:"user_id"`
	LikeableType string `json:"likeable_type"` // "post" or "comment"
	LikeableID   string `json:"likeable_id"`
	CreatedAt    int64  `json:"created_at"`
	DeletedAt    *int64 `json:"deleted_at,omitempty"`
}

type Group struct {
	ID          int       `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	CreatorID   int       `json:"creator_id" db:"creator_id"`
	Creator     *User     `json:"creator,omitempty"`
	MemberCount int       `json:"member_count,omitempty"`
	IsJoined    bool      `json:"is_joined,omitempty"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type GroupMember struct {
    ID       int       `json:"id" db:"id"`
    GroupID  int       `json:"group_id" db:"group_id"`
    UserID   int       `json:"user_id" db:"user_id"`
    User     *User     `json:"user,omitempty"`
    JoinedAt time.Time `json:"joined_at" db:"joined_at"`
}

type InvitationType string

const (
    InvitationTypeInvitation  InvitationType = "invitation"
    InvitationTypeJoinRequest InvitationType = "join_request"
)

type InvitationStatus string

const (
    InvitationStatusPending  InvitationStatus = "pending"
    InvitationStatusAccepted InvitationStatus = "accepted"
    InvitationStatusDeclined InvitationStatus = "declined"
)

type Invitation struct {
    ID         int              `json:"id" db:"id"`
    FromUserID int              `json:"from_user_id" db:"from_user_id"`
    ToUserID   int              `json:"to_user_id" db:"to_user_id"`
    GroupID    int              `json:"group_id" db:"group_id"`
    Type       InvitationType   `json:"type" db:"type"`
    Status     InvitationStatus `json:"status" db:"status"`
    FromUser   *User            `json:"from_user,omitempty"`
    ToUser     *User            `json:"to_user,omitempty"`
    Group      *Group           `json:"group,omitempty"`
    CreatedAt  time.Time        `json:"created_at" db:"created_at"`
    UpdatedAt  time.Time        `json:"updated_at" db:"updated_at"`
}

type Comment struct {
    ID        int       `json:"id" db:"id"`
    PostID    int       `json:"post_id" db:"post_id"`
    AuthorID  int       `json:"author_id" db:"author_id"`
    Content   string    `json:"content" db:"content"`
    Author    *User     `json:"author,omitempty"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Event struct {
    ID          int            `json:"id" db:"id"`
    GroupID     int            `json:"group_id" db:"group_id"`
    CreatorID   int            `json:"creator_id" db:"creator_id"`
    Title       string         `json:"title" db:"title"`
    Description string         `json:"description" db:"description"`
    EventDate   time.Time      `json:"event_date" db:"event_date"`
    Creator     *User          `json:"creator,omitempty"`
    Responses   []EventResponse `json:"responses,omitempty"`
    UserResponse *EventResponse `json:"user_response,omitempty"`
    CreatedAt   time.Time      `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at" db:"updated_at"`
}

type EventResponseType string

const (
    EventResponseGoing    EventResponseType = "going"
    EventResponseNotGoing EventResponseType = "not_going"
)

type EventResponse struct {
    ID        int               `json:"id" db:"id"`
    EventID   int               `json:"event_id" db:"event_id"`
    UserID    int               `json:"user_id" db:"user_id"`
    Response  EventResponseType `json:"response" db:"response"`
    User      *User             `json:"user,omitempty"`
    CreatedAt time.Time         `json:"created_at" db:"created_at"`
    UpdatedAt time.Time         `json:"updated_at" db:"updated_at"`
}

// Request/Response structs
type CreateGroupRequest struct {
    Title       string `json:"title" validate:"required,min=1,max=100"`
    Description string `json:"description" validate:"max=1000"`
}

type InviteUserRequest struct {
    UserID int `json:"user_id" validate:"required"`
}

type JoinRequestRequest struct {
    GroupID int `json:"group_id" validate:"required"`
}

type HandleInvitationRequest struct {
    Action string `json:"action" validate:"required,oneof=accept decline"`
}

type CreatePostRequest struct {
    Title   string `json:"title" validate:"required,min=1,max=200"`
    Content string `json:"content" validate:"required,min=1"`
}

type CreateCommentRequest struct {
    Content string `json:"content" validate:"required,min=1"`
}

type CreateEventRequest struct {
    Title       string    `json:"title" validate:"required,min=1,max=200"`
    Description string    `json:"description" validate:"max=1000"`
    EventDate   time.Time `json:"event_date" validate:"required"`
}

type EventResponseRequest struct {
    Response EventResponseType `json:"response" validate:"required,oneof=going not_going"`
}