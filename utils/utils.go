package utils

import (
	"fmt"
	"net/http"
	"sync"
	"time"
	"trading/settings"
)

type klineParams struct {
	symbol    string
	interval  string
	timeStart int64
	timeEnd   int64
}

func updateKlineData(params klineParams) error {
	req, err := http.NewRequest(http.MethodGet, "https://api.binance.com/api/v3/klines", nil)
	if err != nil {
		_ = req
		return err

	}
	fmt.Printf("%v  %v\n", time.UnixMilli(params.timeStart).UTC(), time.UnixMilli(params.timeEnd).UTC())

	return nil
}

func UpdateTables() error {
	limiter := settings.NewLimiter(2*time.Second, 2)
	wGrp := new(sync.WaitGroup)
	errMsg := settings.NewErrorMessage()

	minTime := time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli()

lb:
	for _, symbol := range settings.Symbols {
		for interval, timeInt := range settings.Intervals {
			currentTime := time.Now().Truncate(timeInt).UnixMilli()
			step := int64(timeInt) / int64(time.Millisecond) * int64(settings.Step)

			for timeStart, timeEnd := currentTime-step, currentTime-int64(time.Nanosecond); timeEnd >
				minTime; timeStart, timeEnd = timeStart-step, timeEnd-step {

				if errMsg.HasError() {
					fmt.Printf("breaking on: %s\n", symbol)
					break lb
				}

				limiter.Wait()
				wGrp.Add(1)

				klParams := klineParams{
					symbol:    symbol,
					interval:  interval,
					timeStart: timeStart,
					timeEnd:   timeEnd,
				}
				go func(params klineParams, wg *sync.WaitGroup, errMessages *settings.ErrorMessages) {
					defer wg.Done()
					err := updateKlineData(params)
					if err != nil {
						errMessages.WriteError(err)
					}
				}(klParams, wGrp, errMsg)
			}
		}
	}
	wGrp.Wait()

	if err := errMsg.GetError(); err != nil {
		return err
	}
	return nil
}
