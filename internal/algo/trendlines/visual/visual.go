package visual

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"image/color"
	"trendlines/dataframe"
)

type GraphPlot struct {
	*plot.Plot
}

// DataFrame reads data from a dataframe and prepares it for plotting.
//
// Parameters:
//   - df: the dataframe containing the data to plot
//   - open: the index of the open price column
//   - close: the index of the close price column
//   - high: the index of the high price column
//   - low: the index of the low price column
func (g *GraphPlot) DataFrame(df dataframe.DataFrame, open, close, high, low int) {
	length := df.Len()
	for i := 0; i < length; i++ {
		cOpen := df.Columns[open].([]float64)[i]
		cClose := df.Columns[close].([]float64)[i]
		cHigh := df.Columns[high].([]float64)[i]
		cLow := df.Columns[low].([]float64)[i]

		highLowLine, _ := plotter.NewLine(plotter.XYs{{X: float64(i), Y: cLow}, {X: float64(i), Y: cHigh}})
		if cClose > cOpen {
			highLowLine.Color = color.RGBA{R: 0, G: 200, B: 0, A: 255} // Green for up
		} else {
			highLowLine.Color = color.RGBA{R: 200, G: 0, B: 0, A: 255} // Red for down
		}

		g.Add(highLowLine)

		body, err := plotter.NewBoxPlot(vg.Points(10), float64(i), plotter.Values{cOpen, cClose})
		body.BoxStyle.Color = color.Transparent
		body.MedianStyle.Color = color.Transparent
		body.WhiskerStyle.Color = color.Transparent

		if err != nil {
			//TODO:write err logs
			panic(err)
		}
		if cClose > cOpen {
			body.FillColor = color.RGBA{R: 0, G: 200, B: 0, A: 255} // Green for up
		} else {
			body.FillColor = color.RGBA{R: 200, G: 0, B: 0, A: 255} // Red for down
		}
		g.Add(body)
	}
}

func (g *GraphPlot) Lines(points []float64) {
	length := len(points)
	xys := make(plotter.XYs, length)
	for i := 0; i < length; i++ {
		xys[i] = plotter.XY{
			X: float64(i),
			Y: points[i],
		}
	}
	//TODO:write err logs
	nL, _ := plotter.NewLine(xys)
	g.Add(nL)
}

func NewPlot() *GraphPlot {
	p := plot.New()
	p.BackgroundColor = color.RGBA{R: 195, G: 195, B: 195, A: 255}
	//p.Title.Text = "BTCUSDT"
	//	//p.X.Label.Text = "Time"
	//	//p.Y.Label.Text = "Price"

	return &GraphPlot{p}
}
