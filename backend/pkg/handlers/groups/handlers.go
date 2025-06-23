package groups

import "database/sql"

type GroupHandler struct {
  db *sql.DB
}

func NewGroupHandler(db *sql.DB) *GroupHandler {
  return &GroupHandler{db: db}
}
