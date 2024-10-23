package data

import (
	"database/sql"
	"embed"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

//go:embed autostart.sql
var embedFiles embed.FS
var DB = getDB()

func DBConnection() (*sql.DB, error) {
	dbPath := "./data/sqlite.db"
	if _, err := os.Stat("./data"); os.IsNotExist(err) {
		if err = os.MkdirAll("./data", 0755); err != nil {
			return nil, err
		}
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, err
}

func getDB() *sql.DB {
	db, err := DBConnection()
	if err != nil {
		log.Fatal(err)
	}

	createTableSQL, errTab := embedFiles.ReadFile("autostart.sql")
	if errTab != nil {
		log.Fatal(errTab)
	}

	_, err = db.Exec(string(createTableSQL))
	if err != nil {
		log.Fatal(err)
	}

	return db
}
