package main

import (
	"fmt"
	"trading/data"
)

func main() {
	d := data.DB
	if err := d.Ping(); err != nil {
		fmt.Errorf(err.Error())
	}

	data.DB.Close()

	if err := d.Ping(); err != nil {
		fmt.Println(err.Error())
	}
}
