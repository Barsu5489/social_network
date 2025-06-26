package repository

import (
	"database/sql"
)

type ChatRepository struct {
	DB *sql.DB
}
