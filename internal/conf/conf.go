package conf

import "time"

const (
	Step     int64 = 500
	KlineURL       = "https://api.binance.com/api/v3/klines"
)

var DateStart = time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC)
