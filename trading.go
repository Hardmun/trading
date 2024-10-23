package main

import (
	"log"
	"trading/utils"
)

func main() {
	if err := utils.CheckLocalTables(); err != nil {
		log.Fatal(err)
	}
}
