package utils

import (
	"fmt"
	"time"
	"trading/data"
	"trading/settings"
)

func writeDataIfNotExists(symbol, interval string, startTime, endTime int64) {

}

func UpdateTables() error {
	minTime := time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli()

	for _, symbol := range settings.Symbols {
		for interval, timeInt := range settings.Intervals {
			currentTime := time.Now().Truncate(timeInt).UnixMilli()
			step := int64(timeInt) / int64(time.Millisecond) * int64(settings.Step)

			for timeStart, timeEnd := currentTime-step, currentTime-int64(time.Nanosecond); timeEnd >
				minTime; timeStart, timeEnd = timeStart-step, timeEnd-step {
				//settings.Limits.Wait()

				writeDataIfNotExists(symbol, interval, timeStart, timeEnd)
				fmt.Printf("%v  %v\n", time.UnixMilli(timeStart).UTC(), time.UnixMilli(timeEnd).UTC())
			}
		}

	}

	db := data.DB

	_ = db
	return nil
}
