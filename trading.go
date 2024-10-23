package main

import (
	"fmt"
	"log"
	"sync"
	"trading/data"
	"trading/settings"
	"trading/utils"
)

func printMe(i int) {
	defer wg.Done()
	settings.Limits.Wait()
	fmt.Println(i)
}

var wg sync.WaitGroup

func main() {
	defer func() {
		if err := data.DB.Close(); err != nil {
			log.Print(err)
		}
	}()

	if err := utils.UpdateTables(); err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go printMe(i)
	}
	wg.Wait()
}
