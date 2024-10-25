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

func updateKlineData(params klineParams) error {
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
		err = data.WriteKlineData(val, fmt.Sprintf("%s_%s", params.symbol, params.interval))
		if err != nil {
			return err
		}
	case interface{}:
		return errors.New("unknown interface{}")
	}

	return nil
}

func UpdateTables() error {
	limiter := settings.NewLimiter(time.Second, 50)
	wGrp := new(sync.WaitGroup)
	errMsg := settings.NewErrorMessage()

	minTime := time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli()
lb:
	for _, symbol := range settings.Symbols {
		for interval, timeInt := range settings.Intervals {
			currentTime := time.Now().Truncate(timeInt).UnixMilli()
			step := int64(timeInt) / int64(time.Millisecond) * int64(settings.Step)

			for timeStart, timeEnd := currentTime-step, currentTime-int64(time.Nanosecond); timeEnd >
				minTime; timeStart, timeEnd = timeStart-step, timeEnd-step {

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
					err := updateKlineData(params)
					if err != nil {
						errMessages.WriteError(err)
					}
				}(klParams, wGrp, errMsg)
			}
		}
	}
	wGrp.Wait()
	if err := errMsg.GetError(); err != nil {
		return err
	}
	return nil
}
