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
var DB = dbConnection()

func dbConnection() *sql.DB {
	dbPath := "./data/sqlite.db"
	if _, err := os.Stat("./data"); os.IsNotExist(err) {
		if err = os.MkdirAll("./data", 0755); err != nil {
			log.Fatal(err)
		}
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
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

func GetDB() *sql.DB {
	if err := DB.Ping(); err != nil {
		DB = dbConnection()
	}
	return DB
}
