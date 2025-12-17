package database

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rhydori/logs"
)

func OpenSQLite(path string) *sql.DB {
	logs.Info("Starting SQLite...")

	_, err := os.Stat(path)
	exists := err == nil
	if exists == false {
		logs.Warnf("%s not found", path)
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			logs.Errorf("OpenSQLite: Error creating dir: %v", err)
			return nil
		}
		logs.Infof("Database directory created: %s", dir)
	}

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		logs.Fatalf("Error opening the Database: %v", err)
		return nil
	}
	_, _ = db.Exec(`PRAGMA foreign_keys = ON`)

	if err := migrate(db); err != nil {
		logs.Fatal(err)
	}
	return db
}
