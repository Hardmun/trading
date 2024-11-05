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

type ColumnType []any

func (c ColumnType) Count() int {
	return len(c)
}

func (c ColumnType) Len() int {
	if len(c) == 0 {
		return 0
	}
	return len(c[0].([]any))
}

func (c ColumnType) Copy(elems ...[2]int) ColumnType {
	width := c.Count()
	cols := make(ColumnType, width)
	if width == 0 {
		return cols
	}

	length := c.Len()
	if len(elems) == 0 {
		for n := 0; n < width; n++ {
			cols[n] = make([]any, length)
			copy(cols[n].([]any), c[n].([]any))
		}
		return cols
	}

	length = Min(length, elems[0][1]-elems[0][0])
	for n := 0; n < width; n++ {
		cols[n] = make([]any, length)
		copy(cols[n].([]any), c[n].([]any)[elems[0][0]:elems[0][0]+length])
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
		if cfg.colTypes != nil {
			switch cfg.colTypes[c] {
			case "float64":
				df.Columns[c] = make([]float64, length)
				continue
			}
		}
		df.Columns[c] = make([]string, length)
	}
	for r := 0; r < length; r++ {
		for c := 0; c < width; c++ {
			if cfg.colTypes != nil {
				switch cfg.colTypes[c] {
				case "float64":
					df.Insert(r, c, records[r][c])
				default:
					df.Insert(r, c, records[r][c])
				}
			} else {
				df.Columns[c].([]any)[r] = records[r][c]
			}
		}
	}
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

func (df *DataFrame) Col(num int) (any, error) {
	if num > df.Columns.Count() {
		return nil, errors.New("column index greater than columns count")
	}
	return df.Columns[num], nil
}

func (df *DataFrame) Insert(row int, col int, val any) {
	switch c := df.Columns[col].(type) {
	case []float64:
		if v, ok := val.(float64); ok {
			c[row] = v
		} else {
			if vs, oks := val.(string); oks {
				var err error
				v, err = strconv.ParseFloat(vs, 64)
				if err != nil {
					df.err = err
					return
				}
				c[row] = v
			} else {
				df.err = errors.New("type not defined")
			}
		}
	case []string:
		if v, ok := val.(string); ok {
			c[row] = v
		} else {
			df.err = errors.New("type not defined")
		}
	default:
		df.err = errors.New("type not defined")
	}
}

func (df *DataFrame) Log(cols ...int) DataFrame {
	colCount := df.Columns.Count()
	newDf := DataFrame{
		Columns: make(ColumnType, colCount),
	}
	for k, c := range cols {
		col, err := df.Col(c)
		if err != nil {
			newDf.err = err
			return newDf
		}
		colFloat64, ok := col.([]float64)
		if !ok {
			newDf.err = errors.New("column type []float64 expected")
			return newDf
		}

		var newSlice []float64
		copy(newSlice, colFloat64)

		newDf.Columns[k] = newSlice
	}
	//newDf := df.Copy()
	//for r := 0; r < newDf.Len(); r++ {
	//	for _, c := range Columns {
	//		newDf.Columns[c].([]any)[r] = math.Log(newDf.Columns[c].([]any)[r].(float64))
	//	}
	//}
	//
	return newDf
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
