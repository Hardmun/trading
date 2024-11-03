package dataframe

import (
	"encoding/csv"
	"errors"
	"io"
)

type loadOptions struct {
	delimiter  rune
	lazyQuotes bool
	comment    rune
	hasHeader  bool
}

type element any

type series struct {
	name     string
	elements element
}

type DataFrame struct {
	columns []series
	ncols   int
	nrows   int
	err     error
}

func (df *DataFrame) LoadRecords(r [][]string, options ...LoadOption) {
	cfg := loadOptions{
		hasHeader: true,
	}
	for _, option := range options {
		option(&cfg)
	}

	//_, _ = headers, records
	//for _, raw := range csvReader[1:] {
	//	for _, col := range raw {
	//		_ = col
	//	}
	//}
}

func HasHeader(has bool) LoadOption {
	return func(cfg *loadOptions) {
		cfg.hasHeader = has
	}
}

type LoadOption func(*loadOptions)

func ReadCSV(r io.Reader, options ...LoadOption) DataFrame {
	dFrame := DataFrame{}
	csvReader := csv.NewReader(r)

	cfg := loadOptions{
		delimiter:  ',',
		lazyQuotes: false,
		comment:    0,
		hasHeader:  true,
	}

	for _, option := range options {
		option(&cfg)
	}

	csvReader.Comma = cfg.delimiter
	csvReader.LazyQuotes = cfg.lazyQuotes
	csvReader.Comment = cfg.comment

	matrix, err := csvReader.ReadAll()
	if err != nil {
		dFrame.err = err
		return dFrame
	}
	if len(matrix) == 0 || len(matrix[0]) <= 1 {
		dFrame.err = errors.New("csv file is empty")
		return dFrame
	}

	dFrame.LoadRecords(matrix, options...)

	return dFrame
}
