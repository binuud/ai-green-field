package scikit

import (
	"math"
	"math/rand"
)

type Point struct {
	X, Y  float64
	Label int
}

// Ai generated code below
func MakeCircles(nSamples int, noise, factor float64, randomState int64) ([]Point, []int) {
	rand.Seed(randomState)
	nOuter, nInner := nSamples/2, (nSamples+1)/2 // Inner gets one more if odd

	points := make([]Point, 0, nSamples)
	labels := make([]int, 0, nSamples)

	// Outer circle (radius 1, label 0)
	for i := 0; i < nOuter; i++ {
		theta := rand.Float64() * 2 * math.Pi
		r := 1.0
		x := r * math.Cos(theta)
		y := r * math.Sin(theta)
		if noise > 0 {
			x += rand.NormFloat64() * noise
			y += rand.NormFloat64() * noise
		}
		points = append(points, Point{X: x, Y: y, Label: 0})
		labels = append(labels, 0)
	}

	// Inner circle (radius factor, label 1)
	for i := 0; i < nInner; i++ {
		theta := rand.Float64() * 2 * math.Pi
		r := factor
		x := r * math.Cos(theta)
		y := r * math.Sin(theta)
		if noise > 0 {
			x += rand.NormFloat64() * noise
			y += rand.NormFloat64() * noise
		}
		points = append(points, Point{X: x, Y: y, Label: 1})
		labels = append(labels, 1)
	}

	// Shuffle
	rand.Shuffle(len(points), func(i, j int) {
		points[i], points[j] = points[j], points[i]
		labels[i], labels[j] = labels[j], labels[i]
	})

	return points[:nSamples], labels[:nSamples] // Trim if needed
}
