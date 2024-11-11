package test

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
	"trading/internal/conf"
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

const loopTestNumber int64 = 5000

// Testing using pocket transfer
func TestGroupWriting(t *testing.T) {
	errMessage = utils.GetErrorMessage()
	limiter := utils.NewLimiter(time.Second, 50)
	prepare(t)

	groupedRecords := GetGroupedRecords()
	go sqlite.BackgroundDBWriter()

	t.Run("Writing messages to database", func(t *testing.T) {
		wg = sync.WaitGroup{}
		lim := make(chan struct{}, 50)
		for _, v := range groupedRecords {
			if errMessage.HasError() {
				break
			}
			limiter.Wait()
			lim <- struct{}{}
			wg.Add(1)
			go func(v [][]any, group *sync.WaitGroup) {
				defer func() {
					<-lim
					group.Done()
				}()

				err := apiEmulation(v)
				if err != nil {
					errMessage.WriteError(err)
				}
			}(v, &wg)
		}
		err := errMessage.GetError()
		if err != nil {
			t.Error(err)
		}

		wg.Wait()
		errMessage.Close()
		close(sqlite.MessageChan)

		var data any
		data, err = sqlite.FetchData("SELECT count() from  BTCUSDT_1h")
		if err != nil {
			t.Error(err)
		}

		if n := conf.Step*loopTestNumber - data.(int64); n > 0 {
			t.Error(errors.New(fmt.Sprintf("Actuall numbers rows less than expected on %v", n)))
		}

	})
}

func BenchmarkSingleRowWriting(b *testing.B) {
	startDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli()
	step := conf.Step
	_, err := sqlite.GetDb()
	if err != nil {
		b.Error(err)
	}

	prepare(b)

	for i := 0; i < b.N; i++ {
		wg.Add(1)
		startDate += startDate + step*int64(time.Second/time.Millisecond)
		SingleRowWriting(step, startDate)
		wg.Wait()
	}
}

// groped by [loopTestNumber]slices
func GetGroupedRecords() [][][]any {
	var grp [][][]any
	sec1 := int64(time.Second / time.Millisecond)
	step := conf.Step

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
	}
	return nil
}

// Batch writing
func SingleRowWriting(step int64, startDate int64) {
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

func prepare(t any) {
	var exeFunc = func() error {
		var err error
		conf.Intervals = map[conf.TimeFrame]time.Duration{
			"1h": time.Hour,
		}

		_, err = sqlite.GetDb()
		if err != nil {
			return err
		}
		err = sqlite.UpdateDatabaseTables()
		if err != nil {
			return err
		}

		err = sqlite.ExecQuery("delete from BTCUSDT_1h", 0)
		if err != nil {
			return err
		}

		return nil
	}

	switch v := t.(type) {
	case *testing.T:
		v.Run("DB connection and table updating", func(t *testing.T) {
			err := exeFunc()
			if err != nil {
				t.Fatal(err)
			}
		})
	case *testing.B:
		v.Run("DB connection and table updating", func(b *testing.B) {
			err := exeFunc()
			if err != nil {
				b.Fatal(err)
			}
		})
	}
}
