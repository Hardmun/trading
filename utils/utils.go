package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
	"trading/data"
	"trading/settings"
)

type klineParams struct {
	symbol    string
	interval  string
	timeStart int64
	timeEnd   int64
}

func getDataFromServer(params klineParams) error {
	var (
		req     *http.Request
		reqResp *http.Response
		bReader []byte
		err     error
	)

	values := url.Values{}
	values.Add("symbol", params.symbol)
	values.Add("interval", params.interval)
	values.Add("startTime", strconv.FormatInt(params.timeStart, 10))
	values.Add("endTime", strconv.FormatInt(params.timeEnd, 10))

	baseURL := fmt.Sprintf("%s?%s", "https://api.binance.com/api/v3/klines", values.Encode())

	req, err = http.NewRequest(http.MethodGet, baseURL, nil)
	if err != nil {
		return err
	}

	client := http.Client{}
	reqResp, err = client.Do(req)
	if err != nil {
		return err
	}

	bReader, err = io.ReadAll(reqResp.Body)
	if err != nil {
		return err
	}

	var resp any
	err = json.Unmarshal(bReader, &resp)
	if err != nil {
		return err
	}
	switch val := resp.(type) {
	case map[string]interface{}:
		var (
			code float64
			msg  string
			ok   bool
		)
		if _, ok = val["code"]; ok {
			code = val["code"].(float64)
		}
		if _, ok = val["msg"]; ok {
			msg = val["msg"].(string)
		}
		return errors.New(fmt.Sprintf("code: %v\nmsg: %s\n", code, msg))
	case []interface{}:

		//err = data.WriteKlineData(val, fmt.Sprintf("%s_%s", params.symbol, params.interval))
		//if err != nil {
		//	return err
		//}
	case interface{}:
		return errors.New("unknown interface{}")
	}

	return nil
}

func max64(a, b int64) int64 {
	if b > a {
		return b
	}
	return a
}

// UpdateTables updates the tables based on the provided update option.
//
//	 1  - Updates only non-existing final records.
//	 0  - Updates all records.
//	-1  - Updates only non-existing records for the entire period.
func UpdateTables(updateOption int8) error {
	limiter := settings.NewLimiter(time.Second, 50)
	wGrp := new(sync.WaitGroup)
	errMsg := settings.NewErrorMessage()

	var lastDate int64
	if updateOption != 1 {
		lastDate = settings.DateStart.UnixMilli()
	}
lb:
	for _, symbol := range settings.Symbols {
		for interval, timeInt := range settings.Intervals {
			currentTime := time.Now().UTC().Truncate(timeInt).UnixMilli()
			step := int64(timeInt) / int64(time.Millisecond) * int64(settings.Step)
			if updateOption == 1 {
				lastDate = data.LastDate(fmt.Sprintf("%s_%s", symbol, interval))
			}
			for timeStart, timeEnd := max64(currentTime-step, lastDate), currentTime-int64(time.Nanosecond); timeEnd >
				lastDate; timeStart, timeEnd = max64(timeStart-step, lastDate), timeEnd-step {

				if errMsg.HasError() {
					break lb
				}

				limiter.Wait()
				wGrp.Add(1)

				klParams := klineParams{
					symbol:    symbol,
					interval:  interval,
					timeStart: timeStart,
					timeEnd:   timeEnd,
				}
				go func(params klineParams, wg *sync.WaitGroup, errMessages *settings.ErrorMessages) {
					defer wg.Done()
					err := getDataFromServer(params)
					if err != nil {
						errMessages.WriteError(err)
					}
				}(klParams, wGrp, errMsg)
			}
		}
	}
	wGrp.Wait()
	errMsg.Close()
	if err := errMsg.GetError(); err != nil {
		return err
	}
	return nil
}
