package dataframe

import (
	"encoding/csv"
	"errors"
	"io"
	"strconv"
)

type LoadOption func(*loadOptions)

type loadOptions struct {
	delimiter  rune
	lazyQuotes bool
	comment    rune
	hasHeader  bool
	colTypes   []string
}

//type element any
//
//type series struct {
//	name     string
//	elements element
//}

type DataFrame struct {
	columns [][]any
	ncols   int
	nrows   int
	err     error
}

func (df *DataFrame) LoadRecords(records [][]string, options ...LoadOption) {
	cfg := loadOptions{
		hasHeader: true,
		colTypes:  nil,
	}
	for _, option := range options {
		option(&cfg)
	}

	if cfg.hasHeader {
		records = records[1:]
	}

	height := len(records)
	width := len(records[0])

	if cfg.colTypes != nil {
		if len(cfg.colTypes) != width {
			df.err = errors.New("mismatch between the number of column types " +
				"provided and the actual number of columns")
		}
	}

	// Determining column width first; rows are greater than columns
	df.columns = make([][]any, width)
	for c := 0; c < width; c++ {
		df.columns[c] = make([]any, height)
	}
	for r := 0; r < height; r++ {
		for c := 0; c < width; c++ {
			if cfg.colTypes != nil {
				var err error
				switch cfg.colTypes[c] {
				case "int":
					df.columns[c][r], err = strconv.ParseInt(records[r][c], 10, 64)
					if err != nil {
						df.err = err
						return
					}
				case "float64":
					df.columns[c][r], err = strconv.ParseFloat(records[r][c], 64)
					if err != nil {
						df.err = err
						return
					}
				default:
					df.columns[c][r] = records[r][c]
				}
			} else {
				df.columns[c][r] = records[r][c]
			}
		}
	}
}

func HasHeader(has bool) LoadOption {
	return func(cfg *loadOptions) {
		cfg.hasHeader = has
	}
}

func ColsTypes(ct []string) LoadOption {
	return func(cfg *loadOptions) {
		cfg.colTypes = ct
	}
}

func ReadCSV(r io.Reader, options ...LoadOption) DataFrame {
	dFrame := DataFrame{}
	csvReader := csv.NewReader(r)

	cfg := loadOptions{
		delimiter:  ',',
		lazyQuotes: false,
		comment:    0,
	}

	for _, option := range options {
		option(&cfg)
	}

	csvReader.Comma = cfg.delimiter
	csvReader.LazyQuotes = cfg.lazyQuotes
	csvReader.Comment = cfg.comment

	records, err := csvReader.ReadAll()
	if err != nil {
		dFrame.err = err
		return dFrame
	}
	if len(records) == 0 || len(records[0]) <= 1 {
		dFrame.err = errors.New("csv file is empty")
		return dFrame
	}

	dFrame.LoadRecords(records, options...)

	return dFrame
}
