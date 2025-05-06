package persistence

import (
	"encoding/json"
	"errors"
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

func (s *JSONStateStorage) LoadState() (state *State, err error) {
	file, err := os.Open(s.path)
	if err != nil {
		return nil, fmt.Errorf("storage file opening: %w", err)
	}
	defer func() {
		if errClose := file.Close(); errClose != nil {
			err = errors.Join(err, fmt.Errorf("storage file closing: %w", errClose))
		}
	}()

	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&state); err != nil {
		return nil, fmt.Errorf("storage decoding: %w", err)
	}
	return
}

func (s *JSONStateStorage) StoreState(state State) error {
	// write to temp file, then move it to destination,
	// to avoid half-written file
	file, creationErr := os.CreateTemp("", "")
	if creationErr != nil {
		return fmt.Errorf("storage file creation: %w", creationErr)
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(state); err != nil {
		return fmt.Errorf("storage encoding: %w", err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("storage file closing: %w", err)
	}

	if renameErr := os.Rename(file.Name(), s.path); renameErr != nil {
		if removeErr := os.Remove(file.Name()); removeErr != nil {
			return fmt.Errorf("storage file replacing: %w", errors.Join(renameErr, removeErr))
		}
		return fmt.Errorf("storage file replacing: %w", renameErr)
	}

	return nil
}
