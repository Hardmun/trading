package main

import (
	"fmt"
	"os"
	df "trend/dataframe"
)

func IsUptrend() bool {
	return true
}

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
	dfLogged := dataFrame.Log([]int{1, 2, 3, 4})

	length := dataFrame.Len()
	//colCount := dataFrame.Columns.Count()
	lookback := 30
	//
	supportSlope := make([]float64, length)
	resistSlope := make([]float64, length)
	//
	candles := dfLogged.Copy([2]int{0, lookback})

	fmt.Println(dfLogged, supportSlope, resistSlope, candles)
}
