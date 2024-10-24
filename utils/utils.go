package utils

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
	"trading/data"
	"trading/settings"
)

type klineParams struct {
	symbol      string
	interval    string
	timeStart   int64
	timeEnd     int64
	wg          *sync.WaitGroup
	limiter     *settings.Limiter
	errMessages *settings.ErrorMessages
}

func writeDataIfNotExists(params klineParams) {
	defer params.wg.Done()

	req, err := http.NewRequest(http.MethodGet, "https://api.binance.com/api/v3/klines", nil)
	if err != nil {
		params.errMessages.WriteError(err)
		_ = req
	}

	fmt.Println(params.symbol)
	if params.timeStart == 1643328000000 {
		params.errMessages.WriteError(errors.New("LOL"))
		return
	}

	time.Sleep(5 * time.Second)
	fmt.Printf("%v  %v\n", time.UnixMilli(params.timeStart).UTC(), time.UnixMilli(params.timeEnd).UTC())

}

func UpdateTables() error {
	limiter := settings.NewLimiter(2*time.Second, 2)
	wg := new(sync.WaitGroup)
	errMessage := settings.NewErrorMessage()

	minTime := time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli()

lb:
	for _, symbol := range settings.Symbols {
		for interval, timeInt := range settings.Intervals {
			currentTime := time.Now().Truncate(timeInt).UnixMilli()
			step := int64(timeInt) / int64(time.Millisecond) * int64(settings.Step)

			for timeStart, timeEnd := currentTime-step, currentTime-int64(time.Nanosecond); timeEnd >
				minTime; timeStart, timeEnd = timeStart-step, timeEnd-step {

				if errMessage.HasError() {
					fmt.Printf("breaking on: %s\n", symbol)
					break lb
				}

				limiter.Wait()
				wg.Add(1)
				params := klineParams{
					symbol:      symbol,
					interval:    interval,
					timeStart:   timeStart,
					timeEnd:     timeEnd,
					wg:          wg,
					limiter:     limiter,
					errMessages: errMessage,
				}
				go writeDataIfNotExists(params)
			}
		}
	}
	wg.Wait()

	if err := errMessage.GetError(); err != nil {
		return err
	}

	db := data.DB

	_ = db
	return nil
}
