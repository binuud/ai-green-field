package neuralNetwork

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func Test_NN_SaveLoad(t *testing.T) {

	dir := t.TempDir()
	filename := filepath.Join(dir, "LinearRegression.model")
	// create random training weights
	model := NewNeuralNetwork(&NeuralNetworkConfig{
		LearningRate: .01,
		Name:         "LinearRegression",
		ModelFile:    filename,
		NumEpochs:    3000,
		Seed:         42.0,
	})

	// Test Save()
	err := model.Save()
	if err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Fatal("Save() did not create file")
	}

	// test model load
	err, newModel := NewNeuralNetworkFromModel(filename)
	if err != nil {
		t.Errorf("Cannot load model from saved file")
	}

	fmt.Printf("Loaded model %v", newModel)

	// For your Save/Load tests
	if reflect.DeepEqual(model.config, newModel.config) {
		fmt.Println("Neural nets match exactly")
	} else {
		t.Errorf("Load model data does not match saved model")
	}

}
