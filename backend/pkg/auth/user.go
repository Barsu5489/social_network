package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
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
// const (
// 	sessionName   = "social-network-session"
// 	sessionMaxAge = 24 * time.Hour
// )

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

// auth/user.go

// Authenticate verifies user credentials and returns the user if valid
func (u *UserModel) Authenticate(email, password string) (*models.User, error) {
	// Query to find user by email
	query := `
        SELECT id, email, password_hash, first_name, last_name, 
               nickname, date_of_birth, about_me, avatar_url, 
               is_private, created_at, updated_at 
        FROM users 
        WHERE email = ? AND deleted_at IS NULL
    `

	var user models.User
	var passwordHash string

	err := u.DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&passwordHash,
		&user.FirstName,
		&user.LastName,
		&user.Nickname,
		&user.DateOfBirth,
		&user.AboutMe,
		&user.AvatarURL,
		&user.IsPrivate,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	// Verify password
fmt.Println(password, passwordHash)

	err = CheckPassword(password, passwordHash)
	if err != nil {
		log.Printf("Password comparison failed: %v", err)
		return nil, errors.New("invalid credentials")
	}
	return &user, nil
}
