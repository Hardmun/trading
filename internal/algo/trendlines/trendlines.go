package main

import (
	"gonum.org/v1/gonum/stat"
	"math"
	"os"
	df "trendlines/dataframe"
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
	length := dataFrame.Len()
	candleCount := 30
	dfMonth := dataFrame.Copy(length-candleCount, length)
	candles := dfMonth.Log(4)

	trendLinesClosePrice(candles)
	//flowCandle = make(chan [4]float64, candleCount)

	//defineTrendLines(candles)

	//var wg sync.WaitGroup
	//for i := candleCount; i < length; i++ {
	//	wg.Add(1)
	//	currentData := dfLogged.Columns.Copy(i, i+1)
	//	flowCandle <- [4]float64{float64(i), currentData[0].([]float64)[0],
	//		currentData[1].([]float64)[0], currentData[2].([]float64)[0]}
	//	time.Sleep(time.Second * 2)
	//}
	//wg.Wait()

	//fmt.Println("DONE")
}

func fitTrendLinesClosePrice(candles []float64) (float64, float64) {
	x := df.Arange(len(candles))
	a, b := stat.LinearRegression(x, candles, nil, false)
	for i := range x {
		x[i] = float64(i)*b + a
	}

	upperPivot := df.Argmax(candles...)
	lowerPivot := df.Argmin(candles...)

	support1, support2 := optimizeSlope(true, lowerPivot, b, candles)
	resist1, resist2 := optimizeSlope(false, upperPivot, b, candles)
	_, _ = resist1, resist2
	return support1, support2
}

func optimizeSlope(support bool, pivot int, initSlope float64, y []float64) (float64, float64) {
	// Amount to change slope by multiply by optStep
	slopeUnit := (df.Max(y...) - df.Min(y...)) / float64(len(y))
	//Optimization variables
	var optStep float64 = 1
	var minStep = 0.0001
	currStep := optStep //current step
	//Initiate at the slope of the line of best fit
	bestSlope := initSlope
	bestErr := checkTrendLine(support, pivot, initSlope, y)
	if bestErr < 0 {
		//TODO: write logs instead panic
		panic("bestErr must to be positive")
	}

	getDerivative := true
	var derivative float64
	for currStep > minStep {
		if getDerivative {
			// Numerical differentiation, increase slope by very small amount
			// to see if error increases/decreases.
			// Gives us the direction to change slope.
			slopeChange := bestSlope + slopeUnit*minStep
			testErr := checkTrendLine(support, pivot, slopeChange, y)
			derivative = testErr - bestErr

			//# If increasing by a small amount fails,
			//# try decreasing by a small amount
			if testErr < 0 {
				slopeChange = bestSlope - slopeUnit*minStep
				testErr = checkTrendLine(support, pivot, slopeChange, y)
				derivative = bestErr - testErr
			}
			if testErr < 0 { // Derivative failed, give up
				//TODO: make logs
				panic("Derivative failed. Check your data. ")
			}
			getDerivative = false
		}
		var testSlope float64
		if derivative > 0 { // Increasing slope increased error
			testSlope = bestSlope - slopeUnit*currStep
		} else { // Increasing slope decreased error
			testSlope = bestSlope + slopeUnit*currStep
		}

		testErr := checkTrendLine(support, pivot, testSlope, y)
		if testErr < 0 || testErr >= bestErr {
			// slope failed/didn't reduce error
			currStep *= 0.5 // Reduce step size
		} else {
			bestErr = testErr
			bestSlope = testSlope
			getDerivative = true // Recompute derivative
		}

	}
	// Optimize done, return best slope and intercept
	return bestSlope, -bestSlope*float64(pivot) + y[pivot]
}

func checkTrendLine(support bool, pivot int, slope float64, y []float64) float64 {
	// compute sum of differences between line and prices,
	// return negative val if invalid

	// Find the intercept of the line going through pivot point with given slope
	intercept := -slope*float64(pivot) + y[pivot]
	lineVals := make([]float64, len(y))
	diffs := make([]float64, len(y))
	for i := range lineVals {
		lineVals[i] = slope*float64(i) + intercept
		diffs[i] = lineVals[i] - y[i]
	}

	//Check to see if the line is valid, return -1 if it is not valid.
	if support && df.Max(diffs...) > 1e-5 {
		return -1.0
	}
	if !support && df.Min(diffs...) < -1e-5 {
		return -1.0
	}

	// Squared sum of diffs between data and line
	var calcErr float64
	for _, v := range diffs {
		calcErr += math.Pow(v, 2)
	}

	return calcErr
}

func trendLinesClosePrice(candles df.DataFrame) {
	supCoff, resCoff := fitTrendLinesClosePrice(candles.Col(0).([]float64))

	println(supCoff, resCoff)
}

//support_coefs_c, resist_coefs_c = fit_trendlines_single(candles['close'])

//func defineTrendLines(candles df.DataFrame) {
//	length := candles.Len()
//	supportSlope := make([]float64, length)
//	resistSlope := make([]float64, length)
//
//	support, resist := trendLineHighLow(
//		candles.Col(0).([]float64),
//		candles.Col(1).([]float64),
//		candles.Col(2).([]float64))
//
//	supportSlope[length-1] = support
//	resistSlope[length-1] = resist
//
//	_, _ = support, resist
//
//	for flow := range flowCandle {
//		_ = flow
//		fmt.Println(flow)
//
//		//supportSlope[int(flow[0])] =
//	}
//}

//func trendLineHighLow(high, low, close []float64) (float64, float64) {
//	x := df.Arange(len(close))
//	a, b := stat.LinearRegression(x, close, nil, false)
//	for i := range x {
//		x[i] = float64(i)*b + a
//	}
//	for i := range high {
//		high[i] = high[i] - x[i]
//	}
//	for i := range low {
//		low[i] = low[i] - x[i]
//	}
//	upperPivot := df.Argmax(high...)
//	lowerPivot := df.Argmin(low...)
//
//	//support := optimizeSlope(true, lowerPivot, b, low)
//
//	fmt.Println(a, b, x, upperPivot, lowerPivot)
//	return 0, 0
//}
