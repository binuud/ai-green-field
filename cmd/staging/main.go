package main

import (
	"fmt"

	bTensor "github.com/binuud/ai-green-field/pkg/bTensor"
	"github.com/sirupsen/logrus"
)

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

}
