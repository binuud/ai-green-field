package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

func linearEquation(m float32, x float32, b float32) float32 {
	return (m*x + b)
}

func main() {
	logrus.Info("Init")
	for i := 0; i <= 100; i++ {
		y := linearEquation(3, float32(i), 5)
		fmt.Printf("\n x : %d, y : %f", i, y)
	}
}
