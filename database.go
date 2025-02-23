package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Database interface {
	GetLastUpdateID() (int, error)
	SetLastUpdateID(id int) error
}

type sqliteDB struct {
	db *sql.DB
}

func NewSQLiteDatabase(path string) (Database, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS state (key TEXT PRIMARY KEY, value INTEGER)")
	if err != nil {
		return nil, err
	}

	return &sqliteDB{db: db}, nil
}

func (s *sqliteDB) GetLastUpdateID() (int, error) {
	var id int
	err := s.db.QueryRow("SELECT value FROM state WHERE key = 'last_update_id'").Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, err
}

func (s *sqliteDB) SetLastUpdateID(id int) error {
	_, err := s.db.Exec("INSERT OR REPLACE INTO state (key, value) VALUES ('last_update_id', ?)", id)
	return err
}
