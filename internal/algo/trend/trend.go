package main

import (
	"fmt"
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
	dfLogged := dataFrame.Log([]int{2, 3, 4})

	length := dfLogged.Len()
	candleCount := 30

	candles := dfLogged.Copy([2]int{0, candleCount})
	flowCandle = make(chan [4]float64, candleCount)
	go defineTrend(candles)

	var wg sync.WaitGroup
	for i := candleCount; i < length; i++ {
		wg.Add(1)
		currentData := dfLogged.Columns.Copy([2]int{i, i + 1})
		flowCandle <- [4]float64{float64(i), currentData[2][0].(float64),
			currentData[3][0].(float64), currentData[4][0].(float64)}
		time.Sleep(time.Second * 2)
	}
	//wg.Wait()

	fmt.Println("DONE")
}

func defineTrend(candles df.DataFrame) {
	length := candles.Len()
	supportSlope := make([]float64, length)
	//resistSlope := make([]float64, length)

	candlesHeight := candles.Col(2)
	candlesLow := candles.Col(3)
	candlesClose := candles.Col(4)

	support, resist := trendLineHighLow(candles.Col(2),
		candles.Col(3),
		candles.Col(4))

	supportSlope[length-1] = support
	supportSlope[length-1] = resist

	_, _ = support, resist

	for flow := range flowCandle {
		_ = flow
		//fmt.Println(flow)

		//supportSlope[int(flow[0])] =
	}
	_, _, _ = candlesHeight, candlesLow, candlesClose
}

func arange(length int) []float64 {
	slc := make([]float64, length)
	for i := range slc {
		slc[i] = float64(i)
	}

	return slc
}

func polyfit(x, y []float64, degree int) []float64 {
	//n := len(x)
	//if n != len(y) {
	//	panic("x and y slices must have the same length")
	//}
	//
	//// Create the Vandermonde matrix
	//vandermonde := mat.NewDense(n, degree+1, nil)
	//for i := 0; i < n; i++ {
	//	for j := 0; j <= degree; j++ {
	//		vandermonde.Set(i, j, floats.Pow(x[i], float64(j)))
	//	}
	//}
	//
	//// Perform least squares fitting
	//var coeff mat.VecDense
	//stat.Regression(&coeff, vandermonde, y, nil)
	//
	//return coeff.RawVector().Data

	return []float64{}
}

func trendLineHighLow(high, low, close []any) (float64, float64) {
	x := arange(len(close))
	//coefs := polyfit(x, close, 1)
	//
	//println(coefs)

	fmt.Println(x)
	return 0, 0
}
