package main

import (
	"fmt"
	"image/color"
	"log"

	bTensor "github.com/binuud/ai-green-field/pkg/bTensor"
	"github.com/sirupsen/logrus"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func plotGraph(xTrain []float64, yTrain []float64, xTest []float64, yTest []float64) {

	// Convert to XYs for plotting
	pts := make(plotter.XYs, len(xTrain))
	for i := range xTrain {
		pts[i].X = xTrain[i]
		pts[i].Y = yTrain[i]
	}

	// Convert to XYs for plotting
	testPts := make(plotter.XYs, len(xTest))
	for i := range xTest {
		testPts[i].X = xTest[i]
		testPts[i].Y = yTest[i]
	}

	// Create plot
	p := plot.New()
	p.Title.Text = "Dot Graph from Arrays"
	p.X.Label.Text = "X Values"
	p.Y.Label.Text = "Y Values"

	// Add scatter points
	s1, err := plotter.NewScatter(pts)
	if err != nil {
		log.Fatal(err)
	}
	s1.GlyphStyle.Color = color.RGBA{255, 0, 0, 255} // Red dots
	s1.Radius = vg.Length(1 * vg.Millimeter)         // Dot size
	p.Add(s1)

	// Add Test scatter points
	s2, err := plotter.NewScatter(testPts)
	if err != nil {
		log.Fatal(err)
	}
	s2.GlyphStyle.Color = color.RGBA{0, 0, 255, 255} // Blue
	s2.Radius = vg.Length(1 * vg.Millimeter)         // Dot size
	p.Add(s2)

	// Save as PNG
	if err := p.Save(5*vg.Inch, 5*vg.Inch, "dotgraph.png"); err != nil {
		logrus.Fatal(err)
	}

}

func main() {

	logrus.Info("Init")

	// Create data
	// start := float32(0.0)
	// end := float32(1.0)
	// step := float32(0.02)

	x := bTensor.NewFromArange(0.0, 1.0, 0.02)
	fmt.Println("Result X:", x.Data[:10])

	y := x.ApplyLinearEquation(0.7, 0.3)
	fmt.Println("Result Y:", y[:10])

	// Create train/test split
	train_split := int(0.8 * float32(len(x.Data))) // 80% of data used for training set, 20% for testing
	X_train, y_train := x.Data[:train_split], y[:train_split]
	X_test, y_test := x.Data[train_split:], y[train_split:]

	fmt.Printf("\n Training Data len %d, %d", len(X_train), len(y_train))
	fmt.Printf("\n Test Data len %d, %d", len(X_test), len(y_test))

	plotGraph(X_train, y_train, X_test, y_test)

}
