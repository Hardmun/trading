package main

import (
	"fmt"
	"github.com/go-gota/gota/dataframe"
	"os"
)

var data dataframe.DataFrame

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

	data = dataframe.ReadCSV(read)
	fmt.Println(data)
}
