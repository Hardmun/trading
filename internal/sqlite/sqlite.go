package sqlite

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"strings"
	"sync"
	"time"
	"trading/internal/config"
	"trading/internal/utils"
	"trading/pgk/queries"
)

type KlineData struct {
	data      []interface{}
	tableName string
}

type MessageDataType struct {
	Query string
	Data  []any
	Wg    *sync.WaitGroup
}

var db *sql.DB
var MessageChan chan MessageDataType

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

func UpdateDatabaseTables() error {
	var finalQuery strings.Builder
	singleQueryString := queries.CreateTables

	for _, symbol := range config.Symbols {
		for interval := range config.Intervals {
			finalQuery.WriteString(strings.Replace(singleQueryString, "&table",
				fmt.Sprintf("%s_%s", symbol, interval), 1))
		}
	}

	finalQuery.WriteString("PRAGMA journal_mode= WAL;")
	_, err := db.Exec(finalQuery.String())
	//err := ExecQuery(finalQuery.String())
	if err != nil {
		return err
	}

	return nil
}
func execQueryConcurrent(query string, wg *sync.WaitGroup, params ...any) {
	defer wg.Done()

	if err := ExecQuery(query, params...); err != nil {
		utils.GetErrorMessage().WriteError(err)
	}
}

func ExecQuery(query string, params ...any) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	_, err = db.Exec(query, params...)
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
		}
	}()
	defer func() {
		if tx != nil {
			_ = tx.Rollback()
		}
	}()

	return tx.Commit()

	//prep, err := db.Prepare(query)
	//if err != nil {
	//	return err
	//}
	//defer func() {
	//	if err = prep.Close(); err != nil {
	//		fmt.Println(err.Error())
	//	}
	//}()
	//_, err = prep.Exec(params...)
	//if err != nil {
	//	return err
	//}
	//return nil
}

func fetchData(query string, params ...any) (any, error) {
	row := db.QueryRow(query, params...)

	var resp any
	err := row.Scan(&resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func LastDate(tableName string) int64 {
	minTime := config.DateStart.UnixMilli()
	query := strings.Replace(queries.QueryLastDay, "&tableName", tableName, 1)

	resultQuery, err := fetchData(query)
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

func BackgroundDBWriter() {
	MessageChan = make(chan MessageDataType, 50)
	for msg := range MessageChan {
		go execQueryConcurrent(msg.Query, msg.Wg, msg.Data...)
	}
}
