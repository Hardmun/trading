package visual

import "trendlines/dataframe"

type Items struct {
	DataFrame dataframe.DataFrame
	//Lines     [][]f
}

type Candle struct {
	Open, Close, High, Low float64
}

func DrawGraph() {

}
