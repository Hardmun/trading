package test

import (
	"strings"
	"sync"
	"testing"
	"time"
	"trading/internal/config"
	"trading/internal/sqlite"
	"trading/internal/utils"
	"trading/pgk/queries"
)

var errMessage *utils.ErrorMessages
var wg sync.WaitGroup
var Record = []any{
	0,
	"49089.49000000",
	"50166.00000000",
	"49013.00000000",
	"50060.01000000",
	"3915.22028000",
	1640282399999,
	"194417790.45947460",
	111929,
	"2255.59685000",
	"112059089.48321300",
}

func BatchWriting(step int, startDate int64) {
	defer wg.Done()
	queryText := strings.Replace(queries.InsertTradingData, "&tableName",
		"BTCUSDT_1h", 1)
	for i := 0; i <= step-1; i++ {
		NewRecord := Record
		NewRecord[0] = startDate
		startDate += int64(time.Second / time.Millisecond)

		err := sqlite.ExecQuery(queryText, NewRecord...)
		if err != nil {
			errMessage.WriteError(err)
		}
	}
}

func PrepareBatchWriting() error {
	startDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli()
	step := config.Step
	for t := 0; t < 2000; t, startDate = t+1, startDate+int64(step)*int64(time.Second/time.Millisecond) {
		if errMessage.HasError() {
			break
		}
		wg.Add(1)
		go BatchWriting(500, startDate)
	}

	wg.Wait()

	err := errMessage.GetError()
	if err != nil {
		return err
	}

	return nil
}

func Test(t *testing.T) {
	errMessage = utils.GetErrorMessage()
	t.Run("DB connection and table updating", func(t *testing.T) {
		var err error
		_, err = sqlite.GetDb()
		if err != nil {
			t.Fatal(err)
		}
		err = sqlite.UpdateDatabaseTables()
		if err != nil {
			t.Fatal(err)
		}
	})
	t.Run("Writing messages to database", func(t *testing.T) {
		var err error
		err = PrepareBatchWriting()
		if err != nil {
			t.Fatal(err)
		}
	})
}
