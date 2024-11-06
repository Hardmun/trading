package visual

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"image/color"
)

type Items struct {
	Candles []Candle
	//Lines     [][]f
}

type Candle struct {
	Open, Close, High, Low float64
}

func DrawGraph(itm Items) {
	p := plot.New()
	p.Title.Text = "Candlestick Chart"
	p.X.Label.Text = "Time"
	p.Y.Label.Text = "Price"

	for i, c := range itm.Candles {
		highLowLine, _ := plotter.NewLine(plotter.XYs{{X: float64(i), Y: c.Low}, {X: float64(i), Y: c.High}})
		highLowLine.Color = color.Black
		p.Add(highLowLine)

		body, err := plotter.NewBoxPlot(vg.Points(10), float64(i), plotter.Values{c.Open, c.Close})
		if err != nil {
			//TODO:make log
			panic(err)
		}
		if c.Close > c.Open {
			body.FillColor = color.RGBA{R: 0, G: 200, B: 0, A: 255} // Green for up
		} else {
			body.FillColor = color.RGBA{R: 200, G: 0, B: 0, A: 255} // Red for down
		}
		p.Add(body)
	}
	_ = p.Save(6*vg.Inch, 4*vg.Inch, "candles.png")
}
