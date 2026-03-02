package neuralNetwork

import (
	"encoding/gob"
	"fmt"
	"os"

	protoV1 "github.com/binuud/ai-green-field/gen/go/v1/neuralNetwork"
	"github.com/sirupsen/logrus"
)

const (
	MODEL_DIRECTORY = "nn_models"
)

func (nn *NeuralNetwork) GetFileName() string {

	// FILE_STORE_PATH is the folder where data is stored as json files
	FILE_STORE_PATH := os.Getenv("FILE_STORE_PATH")
	if FILE_STORE_PATH == "" {
		logrus.Fatal("env FILE_STORE_PATH not set")
	}

	return fmt.Sprintf("%s/%s/%s", FILE_STORE_PATH, MODEL_DIRECTORY, nn.Model.Uuid)

}

// internal function to load from model file
func loadFromModelFile(fileName string) (error, *protoV1.Model) {

	file, err := os.Open(fileName)
	if err != nil {
		return err, nil
	}
	defer file.Close() // Ensure the file is closed

	model := &protoV1.Model{}
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(model) // Decode into the pointer
	if err != nil {
		return err, nil
	}
	return nil, model

}

func (nn *NeuralNetwork) Save(fileName string) error {

	if fileName == "" {
		fileName = nn.GetFileName()
	}

	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := gob.NewEncoder(f)

	// save the file and return error if present
	err = enc.Encode(nn.Model)
	if err != nil {
		logrus.Errorf("Cannot save file %s %v", fileName, err)
		return err
	}
	logrus.Infof("Saved Model to file %s", fileName)
	return nil
}
