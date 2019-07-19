package config

import (
	"os"
	"path/filepath"
)

// Create initializes an empty file at the location
func Create(file string) error {
	if err := os.MkdirAll(filepath.Dir(file), os.ModePerm); err != nil {
		return err
	}
	_, err := os.Create(file)
	return err
}
