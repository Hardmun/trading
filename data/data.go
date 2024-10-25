package data

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"strings"
	"trading/settings"
)

var DB = getDB()

func DBConnect() (*sql.DB, error) {
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

func createNonExistTables(db *sql.DB) error {
	singleQuery, err := settings.EmbedFiles.ReadFile("autostart.sql")
	if err != nil {
		return err
	}

	var finalQuery strings.Builder
	singleQueryString := string(singleQuery)

	for _, symbol := range settings.Symbols {
		for interval := range settings.Intervals {
			finalQuery.WriteString(strings.Replace(singleQueryString, "&table",
				fmt.Sprintf("%s_%s", symbol, interval), 1))
		}
	}

	finalQuery.WriteString("PRAGMA journal_mode= WAL;")
	_, err = db.Exec(finalQuery.String())
	if err != nil {
		return err
	}

	return nil
}

func getDB() *sql.DB {
	db, err := DBConnect()
	if err != nil {
		log.Fatal(err)
	}

	err = createNonExistTables(db)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func executeQuery(query string, params ...any) error {
	prep, err := DB.Prepare(query)
	if err != nil {
		return err
	}
	defer func() {
		if err = prep.Close(); err != nil {
			fmt.Println(err.Error())
		}
	}()
	_, err = prep.Exec(params...)
	if err != nil {
		return err
	}
	return nil
}

func WriteKlineData(data []interface{}) error {
	for _, kl := range data {
		switch klData := kl.(type) {
		case []interface{}:
			if len(klData) == 0 {
				return nil
			}
			query := ""
			err := executeQuery(query, 1, 2, 3)
			if err != nil {
				return err
			}
		default:
			return errors.New("unknown interface{} in func WriteKlineData(data []interface{})")
		}
	}
	return nil
}
