package main

import (
	"gonum.org/v1/gonum/stat"
	"math"
	"os"
	df "trendlines/dataframe"
	"trendlines/visual"
)

func main() {
	read, _ := os.Open("./BTCUSDT86400.csv")
	defer func() {
		_ = read.Close()
	}()

	//TODO: make with CONST
	colsType := []string{
		"string",
		"float64",
		"float64",
		"float64",
		"float64",
	}

	dataFrame := df.ReadCSV(read, df.ColsTypes(colsType))
	length := dataFrame.Len()
	candleCount := 60

	loggedTable := dataFrame.Copy([]int{length - candleCount, length})
	loggedTable.Log(0, 1, 2, 3, 4)

	trendLinesClosePrice(loggedTable.Copy())
}

func fitTrendLinesClosePrice(candles []float64) ([2]float64, [2]float64) {
	length := len(candles)
	x := df.Arange(length, func(t float64, elems ...float64) float64 {
		return t
	})
	a, b := stat.LinearRegression(x, candles, nil, false)

	for i := 0; i < length; i++ {
		x[i] = candles[i] - float64(i)*b + a
	}

	upperPivot := df.Argmax(x...)
	lowerPivot := df.Argmin(x...)

	supportCof := optimizeSlope(true, lowerPivot, b, candles)
	resistCof := optimizeSlope(false, upperPivot, b, candles)

	return supportCof, resistCof
}

func optimizeSlope(support bool, pivot int, initSlope float64, y []float64) [2]float64 {
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
	return [2]float64{bestSlope, -bestSlope*float64(pivot) + y[pivot]}
}

func checkTrendLine(support bool, pivot int, slope float64, y []float64) float64 {
	// compute sum of differences between line and prices,
	// return negative val if invalid

	// Find the intercept of the line going through pivot point with given slope
	length := len(y)
	intercept := -slope*float64(pivot) + y[pivot]
	lineVals := make([]float64, length)
	diffs := make([]float64, length)
	for i := 0; i < length; i++ {
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

func getLinePoints(candles df.DataFrame, sCoff []float64) {

}

func trendLinesClosePrice(candles df.DataFrame) {
	sCoff, rCoff := fitTrendLinesClosePrice(candles.Col(4).([]float64))

	supportLine := df.Arange(candles.Len(), func(t float64, elems ...float64) float64 {
		return t*elems[0] + elems[1]
	}, sCoff[0], sCoff[1])
	resistLine := df.Arange(candles.Len(), func(t float64, elems ...float64) float64 {
		return t*elems[0] + elems[1]
	}, rCoff[0], rCoff[1])

	_, _ = supportLine, resistLine

	var candleVisual = make([]visual.CandleType, 60)
	for r := 0; r < candles.Len(); r++ {
		candleVisual[r] = visual.CandleType{
			Open:  candles.Columns[1].([]float64)[r],
			Close: candles.Columns[4].([]float64)[r],
			High:  candles.Columns[2].([]float64)[r],
			Low:   candles.Columns[3].([]float64)[r],
		}
	}

	itm := visual.Items[[]visual.CandleType]{
		Data:    candleVisual,
		LineSup: supportLine,
		LineRes: resistLine,
	}
	visual.DrawGraph(itm)
}
