package test

import (
	"sync"
	"testing"
	"trading/internal/logs"
)

func TestNewLog(t *testing.T) {
	var wg sync.WaitGroup
	limiter := make(chan struct{}, 1000000)
	l, err := logs.GetErrorLog()
	defer l.Close()
	if err != nil {
		t.Error(err)
	}
	for i := 0; i < 10000000; i++ {
		limiter <- struct{}{}
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			defer func() {
				<-limiter
			}()
			l.Write(i)
		}(&wg)
	}
	wg.Wait()
}
