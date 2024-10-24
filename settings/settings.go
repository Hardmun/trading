package settings

import (
	"embed"
	"time"
)

var (
	Intervals = map[string]time.Duration{
		//"1m": time.Minute,
		//"1h": time.Hour,
		"1d": time.Hour * 24,
	}
	Symbols = []string{
		"BTCUSDT",
	}
	Step = 500
)

type Limiter struct {
	countLimit int
	count      int
	ticker     *time.Ticker
}

//go:embed autostart.sql
var EmbedFiles embed.FS
var Limits = NewLimiter(2*time.Second, 2)

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
