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

type Columns [][]any

func (c Columns) Count() int {
	return len(c)
}

type DataFrame struct {
	Columns Columns
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
				"provided and the actual number of Columns")
		}
	}

	// Determining column width first; rows are greater than Columns
	df.Columns = make([][]any, width)
	for c := 0; c < width; c++ {
		df.Columns[c] = make([]any, length)
	}
	for r := 0; r < length; r++ {
		for c := 0; c < width; c++ {
			if cfg.colTypes != nil {
				var err error
				switch cfg.colTypes[c] {
				case "int":
					df.Columns[c][r], err = strconv.ParseInt(records[r][c], 10, 64)
					if err != nil {
						df.err = err
						return
					}
				case "float64":
					df.Columns[c][r], err = strconv.ParseFloat(records[r][c], 64)
					if err != nil {
						df.err = err
						return
					}
				default:
					df.Columns[c][r] = records[r][c]
				}
			} else {
				df.Columns[c][r] = records[r][c]
			}
		}
	}
}

func Min[nm ~int | ~float64](numbers ...nm) nm {
	length := len(numbers)
	if length == 0 {
		return 0
	}
	if length == 1 {
		return numbers[0]
	}
	minNum := numbers[0]
	for n := 1; n < length; n++ {
		if numbers[n] < minNum {
			minNum = numbers[n]
		}
	}
	return minNum
}

func (df *DataFrame) Copy(elems ...[2]int) DataFrame {
	width := len(df.Columns)
	if width == 0 {
		return DataFrame{}
	}
	startRow := 0
	length := len(df.Columns[0])
	if len(elems) != 0 && length != 0 {
		length = Min(length, elems[0][1]-elems[0][0])
		startRow = elems[0][0]
	}

	if length == 0 {
		return DataFrame{}
	}

	cols := make([][]any, width)
	for c := 0; c < width; c++ {
		cols[c] = make([]any, length)
	}
	for r := startRow; r < length+startRow; r++ {
		for c := 0; c < width; c++ {
			cols[c][r] = df.Columns[c][r]
		}
	}

	return DataFrame{
		Columns: cols,
	}
}

func (df *DataFrame) Len() int {
	if len(df.Columns) == 0 {
		return 0
	}
	return len(df.Columns[0])
}

func (df *DataFrame) Log(Columns []int) DataFrame {
	newDf := df.Copy()
	for r := 0; r < len(newDf.Columns); r++ {
		for _, c := range Columns {
			newDf.Columns[c][r] = math.Log(newDf.Columns[c][r].(float64))
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
