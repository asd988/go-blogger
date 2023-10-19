package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"crypto/sha256"
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
		name TEXT,
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

func StoreFile(name string, data []byte) []byte {
	hasher := sha256.New()
	hasher.Write(data)
	hash := hasher.Sum(nil)

	// Check if the file already exists
	var fileHash []byte
	err := db.QueryRow("SELECT hash FROM file WHERE hash = ?", hash).Scan(&fileHash)
	if err != nil && err != sql.ErrNoRows {
		log.Fatal(err)
	}

	if fileHash == nil {
		_, err := db.Exec("INSERT INTO file(hash, name, data) VALUES(?, ?, ?)", hash, name, data)
		if err != nil {
			log.Fatal(err)
		}
		println("Stored new file")
	}

	return hash
}
