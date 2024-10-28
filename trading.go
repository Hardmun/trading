package main

import (
	"log"
	"trading/api"
	"trading/db"
)

func main() {
	defer func() {
		if err := db.DB.Close(); err != nil {
			log.Print(err)
		}
	}()

	if err := api.UpdateTables(-1); err != nil {
		log.Fatal(err)
	}

}
