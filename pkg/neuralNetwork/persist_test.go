package neuralNetwork

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	protoV1 "github.com/binuud/ai-green-field/gen/go/v1/neuralNetwork"
)

func Test_NN_SaveLoad(t *testing.T) {

	dir := t.TempDir()
	filename := filepath.Join(dir, "LinearRegression.model")
	// create random training weights
	model := NewNeuralNetwork(&protoV1.ModelConfig{
		LearningRate: .01,
		Name:         "LinearRegression",
		Epochs:       3000,
		Seed:         42.0,
	})
	model.Model.Uuid = "fd878f87-b1bf-4848-bdbf-64374f2f0e2b"

	// Test Save()
	err := model.Save(filename)
	if err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Fatal("Save() did not create file")
	}

	// test model load
	newModel, err := NewNeuralNetworkFromModel(filename)
	if err != nil {
		t.Errorf("Cannot load model from saved file")
	}

	fmt.Printf("Loaded model %v", newModel)

	// For your Save/Load tests
	if reflect.DeepEqual(model.Model.Config, newModel.Model.Config) {
		fmt.Println("Neural nets match exactly")
	} else {
		t.Errorf("Load model data does not match saved model")
	}

}
