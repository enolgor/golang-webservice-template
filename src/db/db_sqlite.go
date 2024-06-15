package db

import (
	"database/sql"
	"errors"

	_ "modernc.org/sqlite"
)

const sqlite_create_users_table string = `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		email TEXT NOT NULL,
		roles TEXT NOT NULL
	)
`

const sqlite_create_user_metadata_table string = `
	CREATE TABLE IF NOT EXISTS user_metadata (
		user_id TEXT PRIMARY KEY,
		example TEXT
	)
`

func NewSQLiteDB(filepath string) (*DB, error) {
	conn, err := sql.Open("sqlite", filepath)
	if err != nil {
		return nil, err
	}
	db := &DB{conn}
	return db, initialize_sqlite(conn)
}

func initialize_sqlite(conn *sql.DB) error {
	var err, errs error
	_, err = conn.Exec(sqlite_create_users_table)
	errs = errors.Join(errs, err)
	_, err = conn.Exec(sqlite_create_user_metadata_table)
	return errs
}
