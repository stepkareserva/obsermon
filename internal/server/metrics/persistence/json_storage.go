package persistence

import (
	"encoding/json"
	"fmt"
	"os"
)

type JSONStateStorage struct {
	path string
}

var _ StateStorage = (*JSONStateStorage)(nil)

func NewJSONStateStorage(path string) JSONStateStorage {
	return JSONStateStorage{path: path}
}

func (s *JSONStateStorage) LoadState() (*State, error) {
	file, err := os.Open(s.path)
	if err != nil {
		return nil, fmt.Errorf("storage file opening: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var state State
	if err := decoder.Decode(&state); err != nil {
		return nil, fmt.Errorf("storage decoding: %w", err)
	}
	return &state, nil
}

func (s *JSONStateStorage) StoreState(state State) error {
	// write to temp file, then move it to destination,
	// to avoid half-written file
	file, err := os.CreateTemp("", "")
	if err != nil {
		return fmt.Errorf("storage file creation: %w", err)
	}
	defer os.Remove(file.Name())

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(state); err != nil {
		return fmt.Errorf("storage encoding: %w", err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("storage file closing: %w", err)
	}

	if err := os.Rename(file.Name(), s.path); err != nil {
		return fmt.Errorf("storage file replacing: %w", err)
	}

	return nil
}
