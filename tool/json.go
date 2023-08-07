package tool

import (
	"encoding/json"
	"os"
)

func LoadJSONFile(path string, data interface{}) error {
	filePtr, err := os.Open(path)
	if err != nil {
		return err
	}
	defer filePtr.Close()
	err = json.NewDecoder(filePtr).Decode(data)
	if err != nil {
		return err
	}
	return nil
}
