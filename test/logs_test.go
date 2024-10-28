package logs

import (
	"github.com/Hardmun/trading.git/internal/logs"
	"sync"
	"testing"
)

func TestNewLog(t *testing.T) {
	var wg sync.WaitGroup
	limiter := make(chan struct{}, 1000000)
	l, err := logs.NewLog("ERROR")
	defer l.Close()
	if err != nil {
		t.Error(err)
	}
	for i := 0; i < 10000; i++ {
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
