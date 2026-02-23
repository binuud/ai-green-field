package neuralNetwork

import (
	"image/color"
	"log"

	"github.com/sirupsen/logrus"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

type nnPlotter struct {
	Filename  string
	GraphType string
	plotter   *plot.Plot
}

func NewPlotter(Filename string) *nnPlotter {

	// Create plot
	p := plot.New()
	p.Title.Text = "Scatter Plot Graph"
	p.X.Label.Text = "X Values"
	p.Y.Label.Text = "Y Values"

	return &nnPlotter{
		Filename:  Filename,
		GraphType: "ScatterPlot",
		plotter:   p,
	}
}

func (p *nnPlotter) addPoints(xData []float32, yData []float32, plotColor color.Color) {

	// Convert to XYs for plotting
	pts := make(plotter.XYs, len(xData))
	for i := range xData {
		pts[i].X = float64(xData[i])
		pts[i].Y = float64(yData[i])
	}

	s1, err := plotter.NewScatter(pts)
	if err != nil {
		log.Fatal(err)
	}
	s1.GlyphStyle.Color = plotColor
	s1.Radius = vg.Length(1 * vg.Millimeter) // Dot size
	p.plotter.Add(s1)

}

func (p *nnPlotter) saveGraph() {

	// Save as PNG
	if err := p.plotter.Save(5*vg.Inch, 5*vg.Inch, p.Filename); err != nil {
		logrus.Fatal(err)
	}

}
