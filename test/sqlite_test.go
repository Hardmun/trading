package test

import (
	"database/sql"
	"github.com/Hardmun/trading.git/internal/sqlite"
	"testing"
)

func TestGetDb(t *testing.T) {
	var (
		db  *sql.DB
		err error
	)
	t.Run("Get db", func(t *testing.T) {
		db, err = sqlite.GetDb()
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Close db", func(t *testing.T) {
		err = db.Close()
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Get and close db", func(t *testing.T) {
		db, err = sqlite.GetDb()
		if err != nil {
			t.Fatal(err)
		}
		err = db.Close()
		if err != nil {
			t.Fatal(err)
		}
	})

}
