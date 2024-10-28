package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"strings"
	"time"
	"trading/internal/config"
	"trading/pgk/queries"
)

type KlineData struct {
	data      []interface{}
	tableName string
}

var db *sql.DB

func dbConnection() (*sql.DB, error) {
	dbPath := "./db/sqlite.db"
	if _, err := os.Stat("./db"); os.IsNotExist(err) {
		if err = os.MkdirAll("./db", 0755); err != nil {
			return nil, err
		}
	}

	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	return db, err
}

func createNonExistTables(db *sql.DB) error {
	var finalQuery strings.Builder
	singleQueryString := queries.Queries[0]

	for _, symbol := range config.Symbols {
		for interval := range config.Intervals {
			finalQuery.WriteString(strings.Replace(singleQueryString, "&table",
				fmt.Sprintf("%s_%s", symbol, interval), 1))
		}
	}

	finalQuery.WriteString("PRAGMA journal_mode= WAL;")
	_, err := db.Exec(finalQuery.String())
	if err != nil {
		return err
	}

	return nil
}

func GetDb() (*sql.DB, error) {
	var err error
	if db == nil {
		db, err = dbConnection()
		if err != nil {
			return nil, err
		}
	}

	if err = db.Ping(); err != nil {
		db, err = dbConnection()
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}

func executeQuery(query string, params ...any) error {
	prep, err := db.Prepare(query)
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

func getQuery(query string, params ...any) (any, error) {
	row := db.QueryRow(query, params...)

	var resp any
	err := row.Scan(&resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func WriteKlineData(data []interface{}, tableName string) error {
	for _, kl := range data {
		switch klData := kl.(type) {
		case []interface{}:
			if len(klData) == 0 {
				return nil
			}

			query := strings.Replace(queries.Queries[1], "&tableName", tableName, 1)
			switch dataSlice := kl.(type) {
			case []interface{}:
				err := executeQuery(query, dataSlice[:11]...)
				if err != nil {
					return err
				}
			default:
				return errors.New("unknown interface{} in func WriteKlineData(data []interface{}, tableName string)")
			}
		default:
			return errors.New("unknown interface{} in func WriteKlineData(data []interface{})")
		}
	}
	return nil
}

func LastDate(tableName string) int64 {
	minTime := config.DateStart.UnixMilli()
	query := strings.Replace(queries.Queries[2], "&tableName", tableName, 1)

	resultQuery, err := getQuery(query)
	if err == nil {
		switch t := resultQuery.(type) {
		case int64:
			return t + int64(time.Nanosecond)
		default:
			return minTime
		}
	}

	return minTime
}
