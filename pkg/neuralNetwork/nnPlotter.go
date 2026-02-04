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
	GraphType string
	plotter   *plot.Plot
}

func NewPlotter() *nnPlotter {

	// Create plot
	p := plot.New()
	p.Title.Text = "Scatter Plot Graph"
	p.X.Label.Text = "X Values"
	p.Y.Label.Text = "Y Values"

	return &nnPlotter{
		GraphType: "ScatterPlot",
		plotter:   p,
	}
}

func (p *nnPlotter) addPoints(xData []float64, yData []float64, plotColor color.Color) {

	// Convert to XYs for plotting
	pts := make(plotter.XYs, len(xData))
	for i := range xData {
		pts[i].X = xData[i]
		pts[i].Y = yData[i]
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
	if err := p.plotter.Save(5*vg.Inch, 5*vg.Inch, "dotgraph.png"); err != nil {
		logrus.Fatal(err)
	}

}
