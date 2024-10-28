package main

import (
	"github.com/Hardmun/trading.git/internal/logs"
	"log"
)

// UpdateTables updates the tables based on the provided update option.
//
//	 1  - Updates only non-existing final records.
//	 0  - Updates all records.
//	-1  - Updates only non-existing records for the entire period.
//func updateTradingTables(updateOption int8) error {
//	limiter := config.NewLimiter(time.Second, 50)
//	wGrp := new(sync.WaitGroup)
//	errMsg := config.NewErrorMessage()
//
//	var lastDate int64
//	if updateOption != 1 {
//		lastDate = config.DateStart.UnixMilli()
//	}
//lb:
//	for _, symbol := range config.Symbols {
//		for interval, timeInt := range config.Intervals {
//			currentTime := time.Now().UTC().Truncate(timeInt).UnixMilli()
//			step := int64(timeInt) / int64(time.Millisecond) * int64(config.Step)
//			if updateOption == 1 {
//				lastDate = db.LastDate(fmt.Sprintf("%s_%s", symbol, interval))
//			}
//			for timeStart, timeEnd := max64(currentTime-step, lastDate), currentTime-int64(time.Nanosecond); timeEnd >
//				lastDate; timeStart, timeEnd = max64(timeStart-step, lastDate), timeEnd-step {
//
//				if errMsg.HasError() {
//					break lb
//				}
//
//				limiter.Wait()
//				wGrp.Add(1)
//
//				klParams := klineParams{
//					symbol:    symbol,
//					interval:  interval,
//					timeStart: timeStart,
//					timeEnd:   timeEnd,
//				}
//				go func(params klineParams, wg *sync.WaitGroup, errMessages *config.ErrorMessages) {
//					defer wg.Done()
//					err := requestData(params)
//					if err != nil {
//						errMessages.WriteError(err)
//					}
//				}(klParams, wGrp, errMsg)
//			}
//		}
//	}
//	wGrp.Wait()
//	errMsg.Close()
//	if err := errMsg.GetError(); err != nil {
//		return err
//	}
//	return nil
//}

func main() {
	errLog, err := logs.NewLog("ERROR")
	if err != nil {
		log.Fatal(err)
	}
	defer errLog.Close()

	errLog.Write("Its a first write")

	//defer func() {
	//	if err := sqlite.DB.Close(); err != nil {
	//		log.Print(err)
	//	}
	//}()
	//
	//if err := api.UpdateTables(-1); err != nil {
	//	log.Fatal(err)
	//}
}
