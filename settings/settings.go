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
	Error      chan error
}

//go:embed autostart.sql
var EmbedFiles embed.FS
var Limits = NewLimiter(2*time.Second, 4)

func NewLimiter(d time.Duration, c int) *Limiter {
	limiter := &Limiter{
		countLimit: c,
		count:      c,
		ticker:     time.NewTicker(d),
		Error:      make(chan error, 1),
	}
	return limiter
}

func (l *Limiter) Wait() error {
	select {
	case err := <-l.Error:
		return err
	default:
	}

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

	return nil
}
