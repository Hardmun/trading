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
	_ = dataFrame
	//dfLogged := dataFrame.Log([]int{1, 2, 3, 4})
	//
	//candleCount := 30

	//supportSlope := make([]float64, length)
	//resistSlope := make([]float64, length)
	//
	//candles := dfLogged.Copy([2]int{0, candleCount})
	//fmt.Println(candles)

	cols := make(df.ColumnType, 1)
	cols[0][0] = make([]any, 5)
	//for i := 0; i < 5; i++ {
	//
	//	//cols[0][i] = []any{1, 2, 3}
	//}
	fmt.Println(cols)
}
