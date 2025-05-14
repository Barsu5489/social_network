package auth

import (
	"database/sql"
	"time"

	"social-nework/pkg/models"
	"social-nework/pkg/utils"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *sql.DB
}

// Session store
var store = sessions.NewCookieStore([]byte("secret"))

// Session configuration
const (
    sessionName   = "social-network-session"
    sessionMaxAge = 24 * time.Hour
)
func (u *UserModel) Insert(user models.User) error {
	passHash, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), 12)
	if err != nil {
		return err
	}

	stmt := `
	INSERT INTO users (
		id, email, password_hash, first_name, last_name,
		nickname, date_of_birth, about_me, avatar_url,
		is_private, created_at, updated_at, deleted_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`
	now := time.Now().Unix()

	_, err = u.DB.Exec(stmt,
		user.ID,
		user.Email,
		passHash,
		user.FirstName,
		user.LastName,
		user.Nickname,
		user.DateOfBirth,
		user.AboutMe,
		user.AvatarURL,
		utils.BoolToInt(user.IsPrivate),
		now,
		now,
		utils.NilOrNullInt(user.DeletedAt),
	)

	return err
}
