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

	//lookback := 30
	dataFrame := df.ReadCSV(read, df.ColsTypes(colsType))
	fmt.Println(dataFrame)
}
