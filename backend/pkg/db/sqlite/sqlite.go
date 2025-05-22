package sqlite

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

func NewDB(dataSourceName string) (*sql.DB, error) {
	dbPath := "pkg/db/sqlite/" + dataSourceName
	log.Println("DB path:", dbPath)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://pkg/db/migrations/sqlite", // directory of your migration .sql files
		"sqlite3", driver,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Force re-run migration (optional if DB is freshly deleted)
	if err := m.Up(); err != nil && err.Error() != "no change" {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migration applied successfully")

	// DEBUG: Check if 'users' table exists
	rows, err := db.Query(`SELECT name FROM sqlite_master WHERE type='table'`)
	if err != nil {
		log.Fatalf("Failed to list tables: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var table string
		rows.Scan(&table)
		log.Println("Found table:", table)
	}

	return db, nil
}
