package main

import (
	"log"
	"trading/data"
	"trading/utils"
)

func main() {
	defer func() {
		if err := data.DB.Close(); err != nil {
			log.Print(err)
		}
	}()

	if err := utils.UpdateTables(-1); err != nil {
		log.Fatal(err)
	}

}
