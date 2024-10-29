package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
	"trading/internal/api"
	"trading/internal/config"
	"trading/internal/sqlite"
)

type errorMessages struct {
	error   chan error
	isError chan struct{}
}

func (e *errorMessages) writeError(err error) {
	select {
	case e.error <- err:
		e.isError <- struct{}{}
	default:
	}
}

func (e *errorMessages) getError() error {
	select {
	case err := <-e.error:
		return err
	default:
		return nil
	}
}

func (e *errorMessages) hasError() bool {
	select {
	case <-e.isError:
		return true
	default:
		return false
	}
}

func (e *errorMessages) close() {
	close(e.error)
	close(e.isError)
}

func newErrorMessage() *errorMessages {
	return &errorMessages{
		error:   make(chan error, 1),
		isError: make(chan struct{}, 1),
	}
}

type limit struct {
	countLimit int
	count      int
	ticker     *time.Ticker
	Error      chan error
}

func (l *limit) Wait() {
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

func newLimiter(d time.Duration, c int) *limit {
	l := &limit{
		countLimit: c,
		count:      c,
		ticker:     time.NewTicker(d),
		Error:      make(chan error, 1),
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

// UpdateTradingTables UpdateTables updates the tables based on the provided update option.
//
//	 1  - Updates only non-existing final records.
//	 0  - Updates all records.
//	-1  - Updates only non-existing records for the entire period.
func UpdateTradingTables(updateOption int8) {
	var wGrp sync.WaitGroup

	limiter := newLimiter(time.Second*5, 50)
	errMsg := newErrorMessage()

	var lastDate int64
	if updateOption != 1 {
		lastDate = config.DateStart.UnixMilli()
	}

lb:
	for _, symbol := range config.Symbols {
		for interval, timeInt := range config.Intervals {
			currentTime := time.Now().UTC().Truncate(timeInt).UnixMilli()
			step := int64(timeInt) / int64(time.Millisecond) * int64(config.Step)

			if updateOption == 1 {
				lastDate = sqlite.LastDate(fmt.Sprintf("%s_%s", symbol, interval))
			}
			for timeStart, timeEnd := Max64(currentTime-step, lastDate),
				currentTime-int64(time.Nanosecond); timeEnd > lastDate; timeStart, timeEnd =
				Max64(timeStart-step, lastDate), timeEnd-step {

				if errMsg.hasError() {
					break lb
				}

				limiter.Wait()
				wGrp.Add(1)

				klParams := api.KlineParams{
					Symbol:    symbol,
					Interval:  interval,
					TimeStart: timeStart,
					TimeEnd:   timeEnd,
				}
				go func(params api.KlineParams, wg *sync.WaitGroup, eMessage *errorMessages) {
					defer wg.Done()
					err := api.RequestKlineData(params)
					if err != nil {
						eMessage.writeError(err)
					}
				}(klParams, &wGrp, errMsg)
			}
		}
	}
	wGrp.Wait()
	errMsg.close()

	if err := errMsg.getError(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
