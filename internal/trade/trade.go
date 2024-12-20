package trade

import "time"

type TimeFrame string

const (
	Time_1m  TimeFrame = "1m"
	Time_5m  TimeFrame = "5m"
	Time_15m TimeFrame = "15m"
	Time_30m TimeFrame = "30m"
	Time_1h  TimeFrame = "1h"
	Time_2h  TimeFrame = "2h"
	Time_4h  TimeFrame = "4h"
	Time_1d  TimeFrame = "1d"
)

func (t TimeFrame) Str() string {
	return string(t)
}

func (t TimeFrame) UnixMilli() int64 {
	return int64(Intervals[t] / time.Millisecond)
}

var Intervals = map[TimeFrame]time.Duration{
	Time_1m: time.Minute,
	//Time_5m: time.Minute * 5,
	//Time_15m: time.Minute * 15,
	//Time_1h: time.Hour,
	//"1d": time.Hour * 24,
}

var (
	Symbols = []string{
		"BTCUSDT",
	}
)
