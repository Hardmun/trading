package config

import (
	"time"
)

var (
	Intervals = map[string]time.Duration{
		//"1m": time.Minute,
		//"15m": time.Minute * 15,
		"1h": time.Hour,
		//"1d": time.Hour * 24,
	}
	Symbols = []string{
		"BTCUSDT",
	}
	Step      = 500
	DateStart = time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC)
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

type ErrorMessages struct {
	error   chan error
	IsError chan struct{}
}

func (e *ErrorMessages) WriteError(err error) {
	select {
	case e.error <- err:
		e.IsError <- struct{}{}
	default:
	}
}

func (e *ErrorMessages) GetError() error {
	select {
	case err := <-e.error:
		return err
	default:
		return nil
	}
}

func (e *ErrorMessages) HasError() bool {
	select {
	case <-e.IsError:
		return true
	default:
		return false
	}
}

func (e *ErrorMessages) Close() {
	close(e.error)
	close(e.IsError)
}

func NewErrorMessage() *ErrorMessages {
	return &ErrorMessages{
		error:   make(chan error, 1),
		IsError: make(chan struct{}, 1),
	}
}