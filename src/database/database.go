package database

import (
	"database/sql"
	"go-blogger/src/utils"
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

func CreateSnapshot(blogId string, pageFileHash []byte, otherFileHashes [][]byte) []byte {
	time := time.Now()
	id := utils.GenerateRandomBytes(8)

	_, err := db.Exec("INSERT INTO snapshot(snapshot_id, page_file, creation_date, blog_id) VALUES(?, ?, ?, ?)", id, pageFileHash, time, blogId)
	if err != nil {
		log.Fatal(err)
	}

	for _, hash := range append(otherFileHashes, pageFileHash) {
		_, err := db.Exec("INSERT INTO snapshot_file(snapshot_id, file_id) VALUES(?, ?)", id, hash)
		if err != nil {
			log.Fatal(err)
		}
	}

	return id
}

func CreateBlog(title string, pageFileHash []byte, otherFileHashes [][]byte) string {
	time := time.Now()
	id := utils.GenerateRandomString(6)

	_, err := db.Exec("INSERT INTO blogs(id, title, publish_date) VALUES(?, ?, ?)", id, title, time)
	if err != nil {
		log.Fatal(err)
	}

	snapshotId := CreateSnapshot(id, pageFileHash, otherFileHashes)

	_, err = db.Exec("UPDATE blogs SET snapshot_id = ? WHERE id = ?", snapshotId, id)
	if err != nil {
		log.Fatal(err)
	}

	return id
}

func GetBlogTitle(id string) string {
	var title string
	err := db.QueryRow("SELECT title FROM blogs WHERE id = ?", id).Scan(&title)
	if err != nil {
		if err == sql.ErrNoRows {
			return ""
		}
		log.Fatal(err)
	}

	return title
}

func GetFileContent(hash []byte) []byte {
	var data []byte
	err := db.QueryRow("SELECT data FROM file WHERE hash = ?", hash).Scan(&data)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Fatal(err)
	}

	return data
}

func GetSnapshotContent(id []byte) []byte {
	var pageFileHash []byte
	err := db.QueryRow("SELECT page_file FROM snapshot WHERE snapshot_id = ?", id).Scan(&pageFileHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Fatal(err)
	}

	return GetFileContent(pageFileHash)
}

func GetBlogSnapshotId(id string) []byte {
	var snapshotId []byte
	err := db.QueryRow("SELECT snapshot_id FROM blogs WHERE id = ?", id).Scan(&snapshotId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Fatal(err)
	}

	return snapshotId
}

func GetBlogContent(id string) []byte {
	snapshotId := GetBlogSnapshotId(id)

	return GetSnapshotContent(snapshotId)
}

type FileSummary struct {
	Name string
	Hash []byte
}

func GetSnapshotFiles(id []byte) []FileSummary {
	var files []FileSummary
	rows, err := db.Query("SELECT file.name, file.hash FROM file INNER JOIN snapshot_file ON file.hash = snapshot_file.file_id WHERE snapshot_file.snapshot_id = ?", id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		var hash []byte
		err := rows.Scan(&name, &hash)
		if err != nil {
			log.Fatal(err)
		}

		files = append(files, FileSummary{name, hash})
	}

	return files
}

type Blog struct {
	Title       string
	Id          string
	PublishDate time.Time
}

func GetBlogs(index int, amount int) []Blog {
	var blogs []Blog

	rows, err := db.Query("SELECT title, id, publish_date FROM blogs ORDER BY publish_date DESC LIMIT ? OFFSET ?", amount, index)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var title string
		var id string
		var publishDate time.Time
		err := rows.Scan(&title, &id, &publishDate)
		if err != nil {
			log.Fatal(err)
		}

		blogs = append(blogs, Blog{title, id, publishDate})
	}

	return blogs
}

func GetFile(hash []byte) (string, []byte) {
	var name string
	var data []byte
	err := db.QueryRow("SELECT name, data FROM file WHERE hash = ?", hash).Scan(&name, &data)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		log.Fatal(err)
	}

	return name, data
}
