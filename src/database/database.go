package database

import (
	"database/sql"
	"go-blogger/src/genrandom"
	"log"
	"time"

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
		snapshot_id BLOB,
		FOREIGN KEY (snapshot_id) REFERENCES snapshot(snapshot_id)
	);
	CREATE TABLE IF NOT EXISTS snapshot (
		snapshot_id BLOB PRIMARY KEY,
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
		snapshot_id BLOB,
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

func GetFileByName(name string) []byte {
	var data []byte
	err := db.QueryRow("SELECT data FROM file WHERE name = ? LIMIT 1", name).Scan(&data)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Fatal(err)
	}

	return data
}

func CreateSnapshot(blog_id string, page_file_hash []byte, other_files_hash [][]byte) []byte {
	time := time.Now()
	id := genrandom.GenerateRandomBytes(8)

	_, err := db.Exec("INSERT INTO snapshot(snapshot_id, page_file, creation_date, blog_id) VALUES(?, ?, ?, ?)", id, page_file_hash, time, blog_id)
	if err != nil {
		log.Fatal(err)
	}

	for _, hash := range append(other_files_hash, page_file_hash) {
		_, err := db.Exec("INSERT INTO snapshot_file(snapshot_id, file_id) VALUES(?, ?)", id, hash)
		if err != nil {
			log.Fatal(err)
		}
	}

	return id
}

func CreateBlog(title string, page_file_hash []byte, other_files_hash [][]byte) string {
	time := time.Now()
	id := genrandom.GenerateRandomString(6)

	_, err := db.Exec("INSERT INTO blogs(id, title, publish_date) VALUES(?, ?, ?)", id, title, time)
	if err != nil {
		log.Fatal(err)
	}

	snapshot_id := CreateSnapshot(id, page_file_hash, other_files_hash)

	_, err = db.Exec("UPDATE blogs SET snapshot_id = ? WHERE id = ?", snapshot_id, id)
	if err != nil {
		log.Fatal(err)
	}

	return id
}
