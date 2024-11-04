package main

import (
	"fmt"
	"os"
	df "trend/dataframe"
)

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
	dfLogged := dataFrame.Log([]int{1, 2, 3, 4})

	//length := dataFrame.Len()
	//colCount := dataFrame.Columns.Count()
	candleCount := 30
	//
	//supportSlope := make([]float64, length)
	//resistSlope := make([]float64, length)
	//
	candles := dfLogged.Copy([2]int{0, candleCount})
	//for i := candleCount; i < candles.Len(); i++ {
	//	newColumn := candles.Columns.Copy()
	//
	//}

	fmt.Println(candles)
}
