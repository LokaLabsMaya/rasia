package sqlite

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// DB wraps a SQLite database connection.
type DB struct {
	Client *sql.DB
}

// NewDB opens a SQLite database and creates the secrets schema.
func NewDB(path string) *DB {
	db, err := sql.Open("sqlite3", path+"?_journal_mode=WAL&_foreign_keys=on")
	if err != nil {
		log.Fatalf("sqlite: open %s: %v", path, err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("sqlite: ping %s: %v", path, err)
	}
	if err := migrate(db); err != nil {
		log.Fatalf("sqlite: migrate: %v", err)
	}
	return &DB{Client: db}
}

// Close closes the database connection.
func (d *DB) Close() error {
	return d.Client.Close()
}

// migrate creates the secrets tables if they do not exist.
func migrate(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS namespaces (
			id         TEXT    NOT NULL PRIMARY KEY,
			parent_id  TEXT    DEFAULT NULL,
			name       TEXT    NOT NULL,
			created_at TEXT    NOT NULL DEFAULT (datetime('now')),
			updated_at TEXT    NOT NULL DEFAULT (datetime('now')),
			deleted_at INTEGER NOT NULL DEFAULT 0,
			UNIQUE (parent_id, name, deleted_at)
		);

		CREATE TABLE IF NOT EXISTS secret_files (
			id           TEXT    NOT NULL PRIMARY KEY,
			namespace_id TEXT    NOT NULL,
			name         TEXT    NOT NULL,
			ext          TEXT    NOT NULL,
			created_at   TEXT    NOT NULL DEFAULT (datetime('now')),
			updated_at   TEXT    NOT NULL DEFAULT (datetime('now')),
			deleted_at   INTEGER NOT NULL DEFAULT 0,
			UNIQUE (namespace_id, name, deleted_at),
			FOREIGN KEY (namespace_id) REFERENCES namespaces(id)
		);

		CREATE TABLE IF NOT EXISTS secrets (
			id         TEXT    NOT NULL PRIMARY KEY,
			file_id    TEXT    NOT NULL,
			key_name   TEXT    NOT NULL,
			value_enc  TEXT    NOT NULL,
			created_at TEXT    NOT NULL DEFAULT (datetime('now')),
			updated_at TEXT    NOT NULL DEFAULT (datetime('now')),
			deleted_at INTEGER NOT NULL DEFAULT 0,
			UNIQUE (file_id, key_name, deleted_at),
			FOREIGN KEY (file_id) REFERENCES secret_files(id)
		);

		CREATE TABLE IF NOT EXISTS file_contents (
			id         TEXT    NOT NULL PRIMARY KEY,
			file_id    TEXT    NOT NULL,
			content    TEXT    NOT NULL,
			updated_at TEXT    NOT NULL DEFAULT (datetime('now')),
			deleted_at INTEGER NOT NULL DEFAULT 0,
			UNIQUE (file_id, deleted_at),
			FOREIGN KEY (file_id) REFERENCES secret_files(id)
		);
	`)
	return err
}
