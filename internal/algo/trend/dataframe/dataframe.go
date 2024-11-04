package dataframe

import (
	"encoding/csv"
	"errors"
	"io"
	"math"
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

type DataFrame struct {
	columns [][]any
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

	length := len(records)
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
		df.columns[c] = make([]any, length)
	}
	for r := 0; r < length; r++ {
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

func (df *DataFrame) Copy() DataFrame {
	width := len(df.columns)
	if width == 0 {
		return DataFrame{}
	}
	length := len(df.columns[0])
	if length == 0 {
		return DataFrame{}
	}

	columns := make([][]any, width)
	for c := 0; c < width; c++ {
		columns[c] = make([]any, length)
	}
	for r := 0; r < length; r++ {
		for c := 0; c < width; c++ {
			columns[c][r] = df.columns[c][r]
		}
	}

	return DataFrame{
		columns: columns,
	}
}

func (df *DataFrame) Len() int {
	if len(df.columns) == 0 {
		return 0
	}
	return len(df.columns[0])
}

func (df *DataFrame) Log(columns []int) DataFrame {
	newDf := df.Copy()
	for r := 0; r < len(newDf.columns); r++ {
		for _, c := range columns {
			newDf.columns[c][r] = math.Log(newDf.columns[c][r].(float64))
		}
	}

	return newDf
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
