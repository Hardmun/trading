package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"trading/internal/conf"
	"trading/internal/sqlite"
	"trading/pgk/queries"
)

type KlineParams struct {
	Symbol    string
	Interval  string
	TimeStart int64
	TimeEnd   int64
}

func RequestKlineData(params KlineParams) error {
	var (
		req     *http.Request
		reqResp *http.Response
		bReader []byte
		err     error
	)

	values := url.Values{}
	values.Add("symbol", params.Symbol)
	values.Add("interval", params.Interval)
	values.Add("startTime", strconv.FormatInt(params.TimeStart, 10))
	values.Add("endTime", strconv.FormatInt(params.TimeEnd, 10))

	baseURL := fmt.Sprintf("%s?%s", conf.KlineURL, values.Encode())

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
		query := strings.Replace(queries.InsertTradingData, "&tableName",
			fmt.Sprintf("%s_%s", params.Symbol, params.Interval), 1)

		for _, kl := range val {
			switch klData := kl.(type) {
			case []interface{}:
				if len(klData) == 0 {
					return nil
				}
				switch dataSlice := kl.(type) {
				case []interface{}:
					var dbParams = sqlite.MessageDataType{
						Query:       query,
						Data:        dataSlice[:11],
						WriteOption: 0,
					}
					//wg.Add(1)
					sqlite.MessageChan <- dbParams
				default:
					return errors.New("unknown interface{} in func RequestKlineData(params KlineParams)")
				}
			default:
				return errors.New("unknown interface{} in func RequestKlineData(params KlineParams)")
			}
		}
	case interface{}:
		return errors.New("unknown interface{}")
	}

	return nil
}
