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

type ColumnType []any

func (c ColumnType) Count() int {
	return len(c)
}

func (c ColumnType) Len() int {
	if len(c) == 0 {
		return 0
	}
	switch v := c[0].(type) {
	case []float64:
		return len(v)
	case []string:
		return len(v)
	}

	return 0
}

func (c ColumnType) Copy(elems ...[]int) ColumnType {
	var slctColumns []int
	var cols ColumnType
	width := c.Count()

	if l := len(elems); (l > 1 && len(elems[1]) > width) || (l > 2) {
		return cols
	} else if l == 2 {
		slctColumns = elems[1]
		width = len(slctColumns)
	}

	if width == 0 {
		return cols
	}
	cols = make(ColumnType, width)
	length := c.Len()
	if len(elems) == 0 {
		for n := 0; n < width; n++ {
			switch v := c[n].(type) {
			case []float64:
				cols[n] = make([]float64, length)
				copy(cols[n].([]float64), v)
			case []string:
				cols[n] = make([]string, length)
				copy(cols[n].([]string), v)
			}
		}
		return cols
	}

	if elems[0] != nil {
		length = Min(length, elems[0][1]-elems[0][0])
	}

	setRowsColsFunc := func(i, n int) {
		switch v := c[n].(type) {
		case []float64:
			cols[i] = make([]float64, length)
			if elems[0] != nil {
				copy(cols[i].([]float64), v[elems[0][0]:elems[0][0]+length])
			} else {
				copy(cols[i].([]float64), v)
			}
		case []string:
			cols[i] = make([]string, length)
			if elems[0] != nil {
				copy(cols[i].([]string), v[elems[0][0]:elems[0][0]+length])
			} else {
				copy(cols[i].([]string), v)
			}
		}
	}

	if len(slctColumns) > 0 {
		for i, n := range slctColumns {
			setRowsColsFunc(i, n)
		}
	} else {
		for n := 0; n < width; n++ {
			setRowsColsFunc(n, n)
		}
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

// Copy creates a DataFrame copy with specific indexes rows and columns
//
// Parameters (params ...int):
//   - if len(params) = 0 	- copy the whole DataFrame
//   - if len(params) = 1 	- only rows to select, rows : [start, end] - all columns
//   - if len(params) > 1 	- rows : [start, end], columns [1,3,5] || rows : nil, columns [1,3,5] - all rows
func (df *DataFrame) Copy(elems ...[]int) DataFrame {
	width := df.Columns.Count()
	newDf := DataFrame{}
	if width == 0 {
		newDf.err = errors.New("table is empty")
		return newDf
	}
	length := df.Len()
	if length == 0 {
		newDf.err = errors.New("table is empty")
		return newDf
	}

	cols := df.Columns.Copy(elems...)
	return DataFrame{
		Columns: cols,
	}
}

func (df *DataFrame) Len() int {
	return df.Columns.Len()
}

// Log returns a DataFrame copy with logged columns
//
// Parameters:
//   - column indexes
func (df *DataFrame) Log(cols ...int) {
	if len(cols) == 0 {
		return
	}
	length := df.Len()
	for _, c := range cols {
		col := df.Col(c)
		switch colValue := col.(type) {
		case []float64:
			for i := 0; i < length; i++ {
				colValue[i] = math.Log(colValue[i])
			}
		}
	}
}

func (df *DataFrame) Col(num int) any {
	if num > df.Columns.Count() {
		df.err = errors.New("column index greater than columns count")
		return nil
	}
	return df.Columns[num]
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

func (df *DataFrame) Get(row int, col int) any {
	if df.Columns.Count() < col || df.Len() < row {
		return nil
	}
	switch c := df.Columns[col].(type) {
	case []float64:
		return c[row]
	case []string:
		return c[row]
	}

	return nil
}

//MATH FUNC************************************************************************************************

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

func Max[nm ~int | ~float64](numbers ...nm) nm {
	length := len(numbers)
	if length == 0 {
		return 0
	}
	if length == 1 {
		return numbers[0]
	}
	maxNum := numbers[0]
	for n := 1; n < length; n++ {
		if numbers[n] > maxNum {
			maxNum = numbers[n]
		}
	}
	return maxNum
}

func Argmax(args ...float64) int {
	length := len(args)
	if length == 0 {
		return -1
	}
	mIdx := 0
	for i := 1; i < length; i++ {
		if args[i] > args[mIdx] {
			mIdx = i
		}
	}
	return mIdx
}

func Argmin(args ...float64) int {
	if len(args) == 0 {
		return -1
	}
	mIdx := 0
	for i := 1; i < len(args); i++ {
		if args[i] < args[mIdx] {
			mIdx = i
		}
	}
	return mIdx
}

func Arange(length int, f func(t float64, elems ...float64) float64, elems ...float64) []float64 {
	slc := make([]float64, length)
	for i := 0; i < length; i++ {
		slc[i] = f(float64(i), elems...)
	}

	return slc
}

//MATH FUNC************************************************************************************************

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
