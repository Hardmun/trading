package visual

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"image/color"
)

type Plots string

const (
	Line   Plots = "line"
	Candle Plots = "Candle"
)

type CandleType struct {
	Open, Close, High, Low float64
}

type Items[t []CandleType | []float64] struct {
	Plot Plots
	Data t
}

func DrawGraph(itm Items[[]CandleType]) {
	p := plot.New()
	p.BackgroundColor = color.RGBA{R: 195, G: 195, B: 195, A: 255}

	p.Title.Text = "BTCUSDT"
	//p.X.Label.Text = "Time"
	//p.Y.Label.Text = "Price"

	for i, c := range itm.Data {
		highLowLine, _ := plotter.NewLine(plotter.XYs{{X: float64(i), Y: c.Low}, {X: float64(i), Y: c.High}})
		if c.Close > c.Open {
			highLowLine.Color = color.RGBA{R: 0, G: 200, B: 0, A: 255} // Green for up
		} else {
			highLowLine.Color = color.RGBA{R: 200, G: 0, B: 0, A: 255} // Red for down
		}

		p.Add(highLowLine)

		body, err := plotter.NewBoxPlot(vg.Points(10), float64(i), plotter.Values{c.Open, c.Close})
		body.BoxStyle.Color = color.Transparent
		body.MedianStyle.Color = color.Transparent
		body.WhiskerStyle.Color = color.Transparent

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
	_ = p.Save(10*vg.Inch, 6*vg.Inch, "candles.png")

}
