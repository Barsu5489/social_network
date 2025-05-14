package sqlite

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)


func NewDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "pkg/db/sqlite/" + dataSourceName)
	if err != nil {
		log.Fatal(err)
	}


	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://pkg/db/migrations/sqlite",
		"sqlite3", driver,
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	log.Println("Migration applied successfully")
	return db, nil
}
