package utils

import (
	"os"
	"path/filepath"
	"time"
)

type Limiter struct {
	countLimit int
	count      int
	ticker     *time.Ticker
	Error      chan error
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

func NewLimiter(d time.Duration, c int) *Limiter {
	limiter := &Limiter{
		countLimit: c,
		count:      c,
		ticker:     time.NewTicker(d),
		Error:      make(chan error, 1),
	}
	return limiter
}

func DirPath(path ...string) (string, error) {
	pathDir := filepath.Join(path...)
	if info, errDir := os.Stat(pathDir); errDir != nil || !info.IsDir() {
		if errDir = os.Mkdir(pathDir, os.ModePerm); errDir != nil {
			return "", errDir
		}
	}
	return pathDir, nil
}

func Max64(a, b int64) int64 {
	if b > a {
		return b
	}
	return a
}
