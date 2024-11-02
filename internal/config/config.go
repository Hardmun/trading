package config

import (
	"time"
)

const (
	Step     int64 = 500
	KlineURL       = "https://api.binance.com/api/v3/klines"
)

var (
	Symbols = []string{
		"BTCUSDT",
	}
	Intervals = map[string]time.Duration{
		"1m": time.Minute,
		//"15m": time.Minute * 15,
		//"1h": time.Hour,
		//"1d": time.Hour * 24,
	}

	DateStart = time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC)
)
