package utils

import (
	"fmt"
	"time"
	"trading/data"
	"trading/settings"
)

func UpdateTables() error {
	minTime := time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli()

	for _, symbol := range settings.Symbols {
		for timeString, timeInt := range settings.Intervals {
			currentTime := time.Now().Truncate(time.Hour).UnixMilli()
			_, _ = symbol, timeString
			step := int64(timeInt) / int64(time.Millisecond) * int64(settings.Step)
			_, _ = symbol, timeString
			for timeStart, timeEnd := currentTime-step, currentTime; timeEnd > minTime; timeStart,
				timeEnd = timeStart-step, timeEnd-step {
				settings.Limits.Wait()
				fmt.Printf("%v  %v\n", time.UnixMilli(timeStart).UTC(), time.UnixMilli(timeEnd-int64(time.Nanosecond)).UTC())
			}
		}

	}

	db := data.DB

	_ = db
	return nil
}
