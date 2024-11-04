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

type ColumnType [][]any

func (c ColumnType) Count() int {
	return len(c)
}

func (c ColumnType) Len() int {
	if len(c) == 0 {
		return 0
	}
	return len(c[0])
}

func (c ColumnType) Copy(elems ...[2]int) ColumnType {
	width := c.Count()
	cols := make(ColumnType, width)
	if width == 0 {
		return cols
	}

	length := len(c[0])
	if len(elems) == 0 {
		for n := 0; n < width; n++ {
			cols[n] = make([]any, length)
			copy(cols[n], c[n])
		}
		return cols
	}

	length = Min(length, elems[0][1]-elems[0][0])
	for n := 0; n < width; n++ {
		cols[n] = make([]any, length)
		copy(cols[n], c[n][elems[0][0]:elems[0][0]+length])
	}
	return cols
}

type DataFrame struct {
	Columns ColumnType
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
	df.Columns = make(ColumnType, width)
	for c := 0; c < width; c++ {
		df.Columns[c] = make([]any, length)
	}
	for r := 0; r < length; r++ {
		for c := 0; c < width; c++ {
			if cfg.colTypes != nil {
				var err error
				switch cfg.colTypes[c] {
				case "int64":
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
	width := df.Columns.Count()
	if width == 0 {
		return DataFrame{}
	}
	length := df.Len()
	if length == 0 {
		return DataFrame{}
	}

	cols := df.Columns.Copy(elems...)
	//cols := make(ColumnType, width)
	//if len(elems) == 0 {
	//cols = df.Columns.Copy()
	return DataFrame{
		Columns: cols,
	}
	//}
	//length = Min(length, elems[0][1]-elems[0][0])
	//
	//for c := 0; c < width; c++ {
	//	cols[c] = make([]any, length)
	//	//copy(cols[c], df.Columns[c][elems[0][0]:elems[0][0]+length])
	//}

	//return DataFrame{
	//	Columns: cols,
	//}
}

func (df *DataFrame) Len() int {
	return df.Columns.Len()
}

func (df *DataFrame) Log(Columns []int) DataFrame {
	newDf := df.Copy()
	for r := 0; r < newDf.Len(); r++ {
		for _, c := range Columns {
			newDf.Columns[c][r] = math.Log(newDf.Columns[c][r].(float64))
		}
	}

	return newDf
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
