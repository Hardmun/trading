package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"strings"
	"sync"
	"time"
	"trading/internal/conf"
	"trading/internal/logs"
	"trading/internal/trade"
	"trading/internal/utils"
	"trading/pgk/queries"
)

type KlineData struct {
	data      []interface{}
	tableName string
}

type MessageDataType struct {
	Query       string
	Data        []any
	WriteOption int8
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

	for _, symbol := range trade.Symbols {
		for interval := range trade.Intervals {
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

		_, err = db.Exec(query, params...)
		if err != nil {
			return err
		}

		return tx.Commit()
	}

	_, err = db.Exec(query, params...)
	if err != nil {
		return err
	}

	return nil
}

func FetchSingleData(query string, params ...any) (any, error) {
	row := db.QueryRow(query, params...)

	var resp any
	err := row.Scan(&resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func LastDate(tableName string) int64 {
	minTime := conf.DateStart.UnixMilli()
	query := strings.Replace(queries.QueryLastDay, "&tableName", tableName, 1)

	resultQuery, err := FetchSingleData(query)
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
	MessageChan = make(chan MessageDataType)
	for msg := range MessageChan {
		err := ExecQuery(msg.Query, msg.WriteOption, msg.Data...)
		if err != nil {
			utils.GetErrorMessage().WriteError(err)
		}
	}
}

func CheckTradingData() error {
	var wg sync.WaitGroup
	var chErr = make(chan error, 1)
	var sendErr = func(e error) {
		select {
		case chErr <- e:
		default:
		}
	}

	for _, s := range trade.Symbols {
		for t := range trade.Intervals {
			wg.Add(1)
			intLx := t.UnixMilli()

			func() {
				defer wg.Done()
				tableName := fmt.Sprintf("%s_%s", s, t)
				query := strings.Replace(queries.QueryStartDay, "&tableName", tableName, 1)
				rows, err := db.Query(query)
				if err != nil {
					sendErr(err)
					return
				}
				defer func() {
					err = rows.Close()
					if err != nil {
						log, _ := logs.GetErrorLog()
						log.Write(err)
					}
				}()

				var nextOpenTime int64
				for rows.Next() {
					var openTime int64
					err = rows.Scan(&openTime)
					if err != nil {
						sendErr(err)
					}
					if nextOpenTime != 0 && nextOpenTime != openTime {
						sendErr(errors.New(fmt.Sprintf("miss opentime-%v for %s interval: %s\n"+
							"previous opentime: %v", nextOpenTime, s, t, nextOpenTime-intLx)))
						return
					}
					nextOpenTime = openTime + intLx
				}
			}()
		}
	}

	wg.Wait()
	close(chErr)

	select {
	case err := <-chErr:
		return err
	default:
		return nil
	}
}
