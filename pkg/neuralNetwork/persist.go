package neuralNetwork

import (
	"encoding/gob"
	"fmt"
	"os"

	protoV1 "github.com/binuud/ai-green-field/gen/go/v1/neuralNetwork"
)

func (nn *NeuralNetwork) GetFileName() string {

	return fmt.Sprintf("%s.model", nn.Model.Uuid)

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
	return enc.Encode(nn.Model) // save the file and return error if present
}
