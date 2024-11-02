package utils

import (
	"os"
	"path/filepath"
	"sync"
	"time"
)

type ErrorMessages struct {
	error   chan error
	isError chan struct{}
}

var (
	errorMessage     *ErrorMessages
	onceErrorMessage sync.Once
)

func (e *ErrorMessages) WriteError(err error) {
	select {
	case e.error <- err:
		e.isError <- struct{}{}
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
	case <-e.isError:
		return true
	default:
		return false
	}
}

func (e *ErrorMessages) Close() {
	close(e.error)
	close(e.isError)
}

func NewErrorMessage() *ErrorMessages {
	return &ErrorMessages{
		error:   make(chan error, 1),
		isError: make(chan struct{}, 1),
	}
}

func GetErrorMessage() *ErrorMessages {
	onceErrorMessage.Do(func() {
		errorMessage = NewErrorMessage()
	})

	return errorMessage
}

type Limit struct {
	countLimit int
	count      int
	ticker     *time.Ticker
	limiter    chan struct{}
}

func (l *Limit) Wait() {
	l.limiter <- struct{}{}
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

func (l *Limit) Done() {
	<-l.limiter
}

func NewLimiter(d time.Duration, c int) Limit {
	l := Limit{
		countLimit: c,
		count:      c,
		ticker:     time.NewTicker(d),
		limiter:    make(chan struct{}, c),
	}
	return l
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
