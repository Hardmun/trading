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

var (
	//wg = sync.WaitGroup{}
	w = 0
)

func writeDataIfNotExists(symbol, interval string, startTime, endTime int64, wg *sync.WaitGroup) {
	defer wg.Done()

	req, err := http.NewRequest(http.MethodGet, "https://api.binance.com/api/v3/klines", nil)
	if err != nil {
		settings.Limits.Error <- err
		_ = req
	}
	fmt.Println(symbol)
	if w == 2 {
		settings.Limits.Error <- errors.New("LOL")
		return
	}
	w++

	time.Sleep(5 * time.Second)
	fmt.Printf("%v  %v\n", time.UnixMilli(startTime).UTC(), time.UnixMilli(endTime).UTC())

}

func UpdateTables() (err error) {
	minTime := time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli()
	var wg = new(sync.WaitGroup)

lb:
	for _, symbol := range settings.Symbols {
		for interval, timeInt := range settings.Intervals {
			currentTime := time.Now().Truncate(timeInt).UnixMilli()
			step := int64(timeInt) / int64(time.Millisecond) * int64(settings.Step)

			for timeStart, timeEnd := currentTime-step, currentTime-int64(time.Nanosecond); timeEnd >
				minTime; timeStart, timeEnd = timeStart-step, timeEnd-step {
				err = settings.Limits.Wait()
				if err != nil {
					break lb
				}
				wg.Add(1)
				go writeDataIfNotExists(symbol, interval, timeStart, timeEnd, wg)
			}
		}
	}
	for len(settings.Limits.Error) > 0 {
		err = <-settings.Limits.Error
	}
	wg.Wait()
	select {
	case err = <-settings.Limits.Error:
	default:
	}

	db := data.DB

	_ = db
	return
}
