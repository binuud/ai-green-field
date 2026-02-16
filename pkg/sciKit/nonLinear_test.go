package scikit

import (
	"fmt"
	"testing"
)

func Test_MakeCircles(t *testing.T) {

	points, labels := MakeCircles(100, 0.1, 0.5, 42)

	for i := 0; i < 50; i++ {
		fmt.Printf("Point %d: (%.3f, %.3f), Label: %d\n", i, points[i].X, points[i].Y, labels[i])
	}
	// graph := nn.NewPlotter("circles.png")
	// points, color := MakeCircles(100, 0.03, 0.5, 42)
	// graph.addPoints(xTrain, yTrain, color.RGBA{0, 0, 255, 255})
	// graph.addPoints(xTest, yTest, color.RGBA{0, 255, 0, 255})
	// graph.saveGraph()
}
