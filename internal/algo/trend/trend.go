package main

import (
	"fmt"
	df "github.com/go-gota/gota/dataframe"
	"os"
	"trend/dataframe"
)

func IsUptrend() bool {
	return true
}

func main() {
	read, _ := os.Open("./BTCUSDT86400.csv")
	defer read.Close()
	//bufio.NewReader(read)

	//filer
	//read, err := os.ReadFile("./BTCUSDT86400.csv")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//io.re

	data := dataframe.ReadCSV(read)
	_ = data

	data1 := df.ReadCSV(read)
	fmt.Println(data1.Col("open"))
}
