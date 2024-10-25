package main

import (
	"fmt"
	"time"
)

func main() {
	//date to linux
	utcTame := time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC)
	linuxTime := utcTame.UnixMilli()
	fmt.Printf("utc time:\nutc: %v\nlinux: %v\n", utcTame, linuxTime)

	var lnxTime int64 = 1709302500000
	uTime := time.UnixMilli(lnxTime).UTC()
	fmt.Printf("linux time:\nutc: %v\nlinux: %v\n", uTime, lnxTime)

}
