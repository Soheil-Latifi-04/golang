package main

import (
	"encoding/json"
	"os"
)

type JSONStorage[T any] struct {
	FileName string
}

func NewJSONStorage[T any](fileName string) *JSONStorage[T] {
	return &JSONStorage[T]{FileName: fileName}
}

func (s *JSONStorage[T]) Save(data T) error {
	fileData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.FileName, fileData, 0644)
}

func (s *JSONStorage[T]) Load(data *T) error {
	filedata, err := os.ReadFile(s.FileName)
	if err != nil {
		return err
	}
	return json.Unmarshal(filedata, data)
}
