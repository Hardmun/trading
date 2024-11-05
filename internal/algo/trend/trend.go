package main

import (
	"fmt"
	"gonum.org/v1/gonum/stat"
	"os"
	"sync"
	"time"
	df "trend/dataframe"
)

var flowCandle chan [4]float64

//func IsUptrend() bool {
//	return true
//}

func main() {
	read, _ := os.Open("./BTCUSDT86400.csv")
	defer func() {
		_ = read.Close()
	}()

	colsType := []string{
		"string",
		"float64",
		"float64",
		"float64",
		"float64",
	}

	dataFrame := df.ReadCSV(read, df.ColsTypes(colsType))
	dfLogged := dataFrame.Log(2, 3, 4)

	length := dfLogged.Len()
	candleCount := 30

	candles := dfLogged.Copy(0, candleCount)
	flowCandle = make(chan [4]float64, candleCount)
	go defineTrend(candles)

	var wg sync.WaitGroup
	for i := candleCount; i < length; i++ {
		wg.Add(1)
		currentData := dfLogged.Columns.Copy(i, i+1)
		flowCandle <- [4]float64{float64(i), currentData[0].([]float64)[0],
			currentData[1].([]float64)[0], currentData[2].([]float64)[0]}
		time.Sleep(time.Second * 2)
	}
	//wg.Wait()

	fmt.Println("DONE")
}

func defineTrend(candles df.DataFrame) {
	length := candles.Len()
	supportSlope := make([]float64, length)
	//resistSlope := make([]float64, length)

	//candlesHeight := candles.Col(0)
	//candlesLow := candles.Col(1)
	//candlesClose := candles.Col(2)

	support, resist := trendLineHighLow(candles.Col(0).([]float64),
		candles.Col(1).([]float64),
		candles.Col(2).([]float64))

	supportSlope[length-1] = support
	supportSlope[length-1] = resist

	_, _ = support, resist

	for flow := range flowCandle {
		_ = flow
		fmt.Println(flow)

		//supportSlope[int(flow[0])] =
	}
}

func arange(length int) []float64 {
	slc := make([]float64, length)
	for i := range slc {
		slc[i] = float64(i)
	}

	return slc
}

func polyfit(x, y []float64, degree int) []float64 {

	return nil
}

func trendLineHighLow(high, low, close []float64) (float64, float64) {
	x := arange(len(close))
	a, b := stat.LinearRegression(x, close, nil, false)
	for i := range x {
		x[i] = float64(i)*b + a
	}
	for i := range high {
		high[i] = high[i] - x[i]
	}
	for i := range low {
		low[i] = low[i] - x[i]
	}
	upperPivot := df.Argmax(high...)
	lowerPivot := df.Argmin(low...)

	//support := optimizeSlope(true, lowerPivot, b, low)

	fmt.Println(a, b, x, upperPivot, lowerPivot)
	return 0, 0
}

func optimizeSlope(support bool, pivot int, intiSlope float64, y []float64) {
	//FINAL
}
