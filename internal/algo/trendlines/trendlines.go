package main

import (
	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/plot/vg"
	"math"
	"os"
	df "trendlines/dataframe"
	"trendlines/visual"
)

var grph = visual.NewPlot()

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
	candleCount := 30

	loggedTable := dataFrame.Copy([]int{length - candleCount, length})
	//loggedTable.Log(0, 1, 2, 3, 4)

	trendLinesClosePrice(loggedTable.Copy())
}

func optimizeSlope(support bool, pivot int, initSlope float64, y []float64, startIndex int) [2]float64 {
	// Amount to change slope by multiply by optStep
	slopeUnit := (df.Max(y...) - df.Min(y...)) / float64(len(y))
	//Optimization variables
	var optStep float64 = 1
	var minStep = 0.0001
	currStep := optStep //current step
	//Initiate at the slope of the line of best fit
	bestSlope := initSlope
	bestErr := checkTrendLine(support, pivot, initSlope, y, startIndex)
	if bestErr < 0 {
		//TODO: write err logs
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
			testErr := checkTrendLine(support, pivot, slopeChange, y, startIndex)
			derivative = testErr - bestErr

			//# If increasing by a small amount fails,
			//# try decreasing by a small amount
			if testErr < 0 {
				slopeChange = bestSlope - slopeUnit*minStep
				testErr = checkTrendLine(support, pivot, slopeChange, y, startIndex)
				derivative = bestErr - testErr
			}
			if testErr < 0 { // Derivative failed, give up
				//TODO: write err logs
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

		testErr := checkTrendLine(support, pivot, testSlope, y, startIndex)
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
	//bestSlope := float64(pivot) - y[pivot]

	return [2]float64{bestSlope, -bestSlope*float64(pivot) + y[pivot]}
}

func checkTrendLine(support bool, pivot int, slope float64, y []float64, startIndex int) float64 {
	// compute sum of differences between line and prices,
	// return negative val if invalid

	// Find the intercept of the line going through pivot point with given slope
	length := len(y)
	intercept := -slope*float64(pivot) + y[pivot]
	//lineVals := make([]float64, length)
	diffs := make([]float64, length)
	for i := 0; i < length; i++ {
		//lineVals[i] =
		diffs[i] = slope*float64(i) + intercept - y[i]
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

func fitTrendLinesClosePrice(dataResist, dataSupport []float64) ([]float64, []float64) {
	length := len(dataResist)
	x := df.Arange(length, func(t float64, elems ...float64) float64 {
		return t
	})
	b, m := stat.LinearRegression(x, dataResist, nil, false)
	bS, mS := stat.LinearRegression(x, dataSupport, nil, false)

	xS := make([]float64, length)
	copy(xS, x)
	for i := 0; i < length; i++ {
		x[i] = dataResist[i] - float64(i)*m + b
	}

	for i := 0; i < length; i++ {
		xS[i] = dataSupport[i] - float64(i)*mS + bS
	}

	startIndex := df.Argmax(dataResist...)
	upperPivot := df.Argmax(x...)
	//lowerPivot := df.Argmin(x...)
	//
	//upperPivotS := df.Argmax(xS...)
	lowerPivotS := df.Argmin(xS...)

	supportCof := optimizeSlope(true, lowerPivotS, mS, dataSupport, startIndex)
	resistCof := optimizeSlope(false, upperPivot, m, dataResist, startIndex)

	supportLine := df.Arange(length, func(t float64, elems ...float64) float64 {
		return t*supportCof[0] + supportCof[1]
	})
	resistLine := df.Arange(length, func(t float64, elems ...float64) float64 {
		return t*resistCof[0] + resistCof[1]
	})

	midle := df.Arange(length, func(t float64, elems ...float64) float64 {
		return t*m + b
	})
	grph.Lines(midle)

	return supportLine, resistLine
}

func trendLinesClosePrice(candles df.DataFrame) {
	realColResist := make([]float64, candles.Len())
	for i := 0; i < candles.Len(); i++ {
		realColResist[i] = df.Max(candles.Columns[1].([]float64)[i], candles.Columns[4].([]float64)[i])
	}

	realColSupport := make([]float64, candles.Len())
	for i := 0; i < candles.Len(); i++ {
		realColSupport[i] = df.Min(candles.Columns[1].([]float64)[i], candles.Columns[4].([]float64)[i])
	}

	supportLine, resistLine := fitTrendLinesClosePrice(realColResist, realColSupport)
	_, _ = supportLine, resistLine

	grph.DataFrame(candles, 1, 4, 2, 3)
	grph.Lines(resistLine)
	//graph.Lines(supportLine)
	grph.Lines(realColResist)
	_ = grph.Save(10*vg.Inch, 6*vg.Inch, "candles.png")
}
