package main

import (
	"log"
	"trading/internal/logs"
	"trading/internal/sqlite"
	"trading/internal/utils"
)

func main() {
	errLog, err := logs.GetErrorLog()
	if err != nil {
		log.Fatal(err)
	}
	defer errLog.Close()

	db, errDb := sqlite.GetDb()
	if errDb != nil {
		errLog.Fatal(errDb)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			errLog.Write(err)
		}
	}()

	utils.UpdateTradingTables(-1)
}
