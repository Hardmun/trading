package utils

import "trading/data"

func UpdateTables() error {
	db := data.DB

	_ = db
	return nil
}
