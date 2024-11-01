package test

import (
	"errors"
	"fmt"
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
var queryText = strings.Replace(queries.InsertTradingData, "&tableName",
	"BTCUSDT_1h", 1)

const loopTestNumber int64 = 500

// Batch writing
func BatchWriting(step int64, startDate int64) {
	defer wg.Done()
	for i := int64(0); i < step; i++ {
		NewRecord := Record
		NewRecord[0] = startDate
		startDate += int64(time.Second / time.Millisecond)

		err := sqlite.ExecQuery(queryText, 0, NewRecord...)
		if err != nil {
			errMessage.WriteError(err)
		}
	}
}

func PrepareBatchWriting() error {
	startDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli()
	step := config.Step
	for t := 0; t < 1; t, startDate = t+1, startDate+step*int64(time.Second/time.Millisecond) {
		if errMessage.HasError() {
			break
		}
		wg.Add(1)
		go BatchWriting(step, startDate)
	}

	wg.Wait()

	err := errMessage.GetError()
	if err != nil {
		return err
	}

	return nil
}

func TestBatchWriting(t *testing.T) {
	errMessage = utils.GetErrorMessage()
	t.Run("DB connection and table updating", func(t *testing.T) {
		var err error
		_, err = sqlite.GetDb()
		if err != nil {
			t.Error(err)
		}
		err = sqlite.UpdateDatabaseTables()
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("Writing messages to database", func(t *testing.T) {
		var err error
		err = PrepareBatchWriting()
		if err != nil {
			t.Error(err)
		}
	})
}

func BenchmarkBatchWriting(b *testing.B) {
	startDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli()
	step := config.Step
	_, err := sqlite.GetDb()
	if err != nil {
		b.Error(err)
	}
	err = sqlite.UpdateDatabaseTables()
	if err != nil {
		b.Error(err)
	}

	for i := 0; i < b.N; i++ {
		wg.Add(1)
		startDate += startDate + step*int64(time.Second/time.Millisecond)
		BatchWriting(step, startDate)
		wg.Wait()
	}
}

// Sync writing
func TestSyncWriting(t *testing.T) {
	t.Run("DB connection and table updating", func(t *testing.T) {
		var err error
		_, err = sqlite.GetDb()
		if err != nil {
			t.Error(err)
		}
		err = sqlite.UpdateDatabaseTables()
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("Writing messages to database", func(t *testing.T) {
		startDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli()

		for i := 0; i < 1000000; i++ {
			NewRecord := Record
			NewRecord[0] = startDate
			startDate += int64(time.Second / time.Millisecond)

			err := sqlite.ExecQuery(queryText, 0, NewRecord...)
			if err != nil {
				errMessage.WriteError(err)
			}
		}
	})

}

// groped by [loopTestNumber]slices
func GetGroupedRecords() [][][]any {
	var grp [][][]any
	sec1 := int64(time.Second / time.Millisecond)
	step := config.Step

	startDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	startTime := startDate.UnixMilli()
	endTime := startTime + loopTestNumber*step*sec1

	var arrayStep = make([][]any, step)

	for s := int64(0); startTime <= endTime; startTime, s = startTime+sec1, s+1 {
		if s >= step {
			s = 0
			grp = append(grp, arrayStep)
			arrayStep = make([][]any, step)
		}
		newStartTime := startTime
		NewRecord := append([]any{}, Record...)
		NewRecord[0] = newStartTime
		arrayStep[s] = NewRecord
	}
	return grp
}

func apiEmulation(v [][]any) error {

	for _, r := range v {
		newMessage := sqlite.MessageDataType{
			Query:       queryText,
			Data:        r,
			WriteOption: 0,
		}

		sqlite.MessageChan <- newMessage
		//err := sqlite.ExecQuery(queryText, 0, r...)
		//if err != nil {
		//	return err
		//}
	}
	return nil
}

// Testing using pocket transfer
// TODO: SQLite doesnt support more than 5000 concurrent connections
func TestGroupWriting(t *testing.T) {
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

		err = sqlite.ExecQuery("delete from BTCUSDT_1h", 0)
		if err != nil {
			t.Fatal(err)
		}
	})

	groupedRecords := GetGroupedRecords()
	go sqlite.BackgroundDBWriter()

	t.Run("Writing messages to database", func(t *testing.T) {
		for _, v := range groupedRecords {
			if errMessage.HasError() {
				break
			}
			func(v [][]any) {
				err := apiEmulation(v)
				if err != nil {
					errMessage.WriteError(err)
				}
			}(v)
		}
		err := errMessage.GetError()
		if err != nil {
			t.Error(err)
		}

		var data any
		data, err = sqlite.FetchData("SELECT count() from  BTCUSDT_1h")
		if err != nil {
			t.Error(err)
		}

		if n := config.Step*loopTestNumber - data.(int64); n > 0 {
			t.Error(errors.New(fmt.Sprintf("Actuall numbers rows less than expected on %v", n)))
		}

	})
}
