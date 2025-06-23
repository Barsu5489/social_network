package models

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
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatorID   string `json:"creator_id"`
	IsPrivate   bool   `json:"is_private"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

type GroupMember struct {
	ID       string `json:"id"`
	GroupID  string `json:"group_id"`
	UserID   string `json:"user_id"`
	Role     string `json:"role"`
	JoinedAt int64  `json:"joined_at"`
}

type Invitation struct {
	ID         string `json:"id"`
	InviterID  string `json:"inviter_id"`
	InviteeID  string `json:"invitee_id"`
	EntityType string `json:"entity_type"`
	EntityID   string `json:"entity_id"`
	Status     string `json:"status"`
	CreatedAt  int64  `json:"created_at"`
}

type Event struct {
	ID          string `json:"id"`
	GroupID     string `json:"group_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Location    string `json:"location"`
	StartTime   int64  `json:"start_time"`
	EndTime     int64  `json:"end_time"`
	CreatedBy   string `json:"created_by"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

type EventAttendee struct {
	ID        string `json:"id"`
	EventID   string `json:"event_id"`
	UserID    string `json:"user_id"`
	Status    string `json:"status"` // going, maybe, not_going
	CreatedAt int64  `json:"created_at"`
}

// Comment represents a comment on a post
type Comment struct {
	ID        string  `json:"id"`
	PostID    string  `json:"post_id"`
	UserID    string  `json:"user_id"`
	Content   string  `json:"content"`
	ImageURL  *string `json:"image_url,omitempty"`
	CreatedAt int64   `json:"created_at"`
	UpdatedAt int64   `json:"updated_at"`
	DeletedAt *int64  `json:"deleted_at,omitempty"`

	// Additional fields for API responses
	UserNickname string `json:"user_nickname,omitempty"`
	UserAvatar   string `json:"user_avatar,omitempty"`
	LikesCount   int    `json:"likes_count"`
	UserLiked    bool   `json:"user_liked,omitempty"`
}
