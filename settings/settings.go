package settings

import (
	"embed"
	"time"
)

var (
	Intervals = []string{
		"1h", "4h", "1d",
	}
	Symbols = []string{
		"BTCUSDT",
	}
	Limit     = 500
	DateStart = time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	DateEnd   = time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC)
)

type Limiter struct {
	countLimit int
	count      int
	ticker     *time.Ticker
}

//go:embed autostart.sql
var EmbedFiles embed.FS
var Limits = NewLimiter(2*time.Second, 4)

func NewLimiter(d time.Duration, c int) *Limiter {
	limiter := &Limiter{
		countLimit: c,
		count:      c,
		ticker:     time.NewTicker(d),
	}
	return limiter
}

func (l *Limiter) Wait() {
	select {
	case <-l.ticker.C:
		l.count = l.countLimit
	default:
	}
	if l.count <= 0 {
		<-l.ticker.C
		l.count = l.countLimit
	}
	l.count--
}
