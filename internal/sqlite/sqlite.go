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
	err := ExecQuery(finalQuery.String(), 0)
	if err != nil {
		return err
	}

	return nil
}
func execQueryConcurrent(query string, wg *sync.WaitGroup, params ...any) {
	defer wg.Done()

	if err := ExecQuery(query, 0, params...); err != nil {
		utils.GetErrorMessage().WriteError(err)
	}
}

// ExecQuery writes a message to the database based on the provided query and write option.
//
// Parameters:
//   - query (string): The SQL query to be executed.
//   - writeOption (int8): Specifies the writing option to the database:
//     0 - Direct writing without preparation.
//     1 - Writing using prepared statements.
//     2 - Writing within a transaction block.
//   - params (...any): Optional parameters for the query.
//
// Returns:
//   - error: Returns an error if the query execution fails, otherwise nil.
func ExecQuery(query string, writeOption int8, params ...any) error {
	var (
		err error
		tx  *sql.Tx
	)
	if writeOption == 1 {
		var prep *sql.Stmt
		prep, err = db.Prepare(query)
		if err != nil {
			return err
		}
		defer func() {
			if err = prep.Close(); err != nil {
				fmt.Println(err.Error())
			}
		}()
		_, err = prep.Exec(params...)

		return nil
	} else if writeOption == 2 {
		tx, err = db.Begin()
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
	}

	_, err = db.Exec(query, params...)
	if err != nil {
		return err
	}

	return nil
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
