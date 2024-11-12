package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
	"trading/internal/api"
	"trading/internal/conf"
	"trading/internal/logs"
	"trading/internal/mth"
	"trading/internal/sqlite"
	"trading/internal/trade"
	"trading/internal/utils"
)

// UpdateTimeFrameData UpdateTradingTables UpdateTables updates the tables based on the provided update option.
//
//	 1  - Updates only non-existing final records.
//	 0  - Updates all records.
//	-1  - Updates only non-existing records for the entire period.
func UpdateTimeFrameData(updateOption int8, currTime time.Time) {
	var wGrp sync.WaitGroup

	limiter := utils.NewLimiter(time.Second, 50)
	routineLimiter := make(chan struct{}, 100)
	errMsg := utils.GetErrorMessage()

	var lastDate int64
	if updateOption != 1 {
		lastDate = conf.DateStart.UnixMilli()
	}

lb:
	for _, symbol := range trade.Symbols {
		for interval, timeInt := range trade.Intervals {
			currentTime := currTime.Truncate(timeInt).UnixMilli()
			step := int64(timeInt) / int64(time.Millisecond) * conf.Step

			if updateOption == 1 {
				lastDate = sqlite.LastDate(fmt.Sprintf("%s_%s", symbol, interval))
			}
			//TODO:remove after test
			brk := 0
			for timeStart, timeEnd := mth.Max64(currentTime-step, lastDate),
				currentTime-int64(time.Nanosecond); timeEnd > lastDate; timeStart, timeEnd =
				mth.Max64(timeStart-step, lastDate), timeEnd-step {

				brk++
				if brk > 10 {
					break
				}

				if errMsg.HasError() {
					break lb
				}

				routineLimiter <- struct{}{}
				limiter.Wait()
				wGrp.Add(1)

				klParams := api.KlineParams{
					Symbol:    symbol,
					Interval:  interval.Str(),
					TimeStart: timeStart,
					TimeEnd:   timeEnd,
				}
				go func(params api.KlineParams, wg *sync.WaitGroup, eMessage *utils.ErrorMessages) {
					defer func() {
						<-routineLimiter
						wg.Done()
					}()
					err := api.RequestKlineData(params)
					if err != nil {
						eMessage.WriteError(err)
					}
				}(klParams, &wGrp, errMsg)
			}
		}
	}
	wGrp.Wait()
	errMsg.Close()
	close(routineLimiter)

	if err := errMsg.GetError(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	//1. Log initialization
	errLog, err := logs.GetErrorLog()
	if err != nil {
		log.Fatal(err)
	}
	defer errLog.Close()

	//2. DB initialization
	db, errDb := sqlite.GetDb()
	if errDb != nil {
		errLog.Fatal(errDb)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			errLog.Write(err)
		}
	}()

	//3. DB tables updates
	err = sqlite.UpdateDatabaseTables()
	if err != nil {
		errLog.Fatal()
	}

	//4. Background DB query receiver
	go sqlite.BackgroundDBWriter()

	currentTime := time.Now().UTC()
	//5. Uploading new trading data
	UpdateTimeFrameData(-1, currentTime)
}
