package neuralNetwork

import (
	"encoding/gob"
	"fmt"
	"os"
)

func (nn *neuralNetwork) GetFileName() string {
	// if model file is set, use same
	if len(nn.config.ModelFile) > 0 {
		return nn.config.ModelFile
	}
	// if model file is unset, generate file name from model
	nn.config.ModelFile = fmt.Sprintf("%s.model", nn.config.Name)
	return nn.config.ModelFile
}

// internal function to load config from model file
func loadFromModelFile(fileName string) (error, *NeuralNetworkConfig) {

	file, err := os.Open(fileName)
	if err != nil {
		return err, nil
	}
	defer file.Close() // Ensure the file is closed

	config := &NeuralNetworkConfig{}
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(config) // Decode into the pointer
	if err != nil {
		return err, nil
	}
	return nil, config

}

func (nn *neuralNetwork) Save() error {
	f, err := os.Create(nn.GetFileName())
	if err != nil {
		return err
	}
	defer f.Close()
	enc := gob.NewEncoder(f)
	return enc.Encode(nn.config) // save the file and return error if present
}
