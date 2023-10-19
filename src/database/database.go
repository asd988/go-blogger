package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitDB() {
	var err error
	db, err = sql.Open("sqlite3", "blog.db")
	if err != nil {
		log.Fatal(err)
	}

	// Create the "blogs" table if it doesn't exist
	createTables := `
    CREATE TABLE IF NOT EXISTS blogs (
		id TEXT PRIMARY KEY,
		title TEXT,
		publish_date DATETIME,
		snapshot_id TEXT,
		FOREIGN KEY (snapshot_id) REFERENCES snapshot(snapshot_id)
	);
	CREATE TABLE IF NOT EXISTS snapshot (
		snapshot_id TEXT PRIMARY KEY,
		page_file BLOB,
		creation_date DATETIME,
		blog_id TEXT,
		FOREIGN KEY (blog_id) REFERENCES blogs(id),
		FOREIGN KEY (page_file) REFERENCES file(hash)
	);
	CREATE TABLE IF NOT EXISTS file (
		hash BLOB PRIMARY KEY,
		extension TEXT,
		data BLOB
	);	
	CREATE TABLE IF NOT EXISTS snapshot_file (
		snapshot_id TEXT,
		file_id BLOB,
		FOREIGN KEY (snapshot_id) REFERENCES snapshot(snapshot_id),
		FOREIGN KEY (file_id) REFERENCES file(hash)
	);	
    `

	_, err = db.Exec(createTables)
	if err != nil {
		log.Fatal(err)
	}
}
