package dataframe

import (
	"encoding/csv"
	"io"
)

type element any

type series struct {
	name     string
	elements element
}

type DataFrame struct {
	columns []series
	ncols   int
	nrows   int
}

func ReadCSV(r io.Reader) DataFrame {
	csvReader := csv.NewReader(r)
	_ = csvReader

	return DataFrame{}
}
