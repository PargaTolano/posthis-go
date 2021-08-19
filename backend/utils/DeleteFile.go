package utils

import (
	"os"
	"path/filepath"
)

func DeleteStaticFile(path string) error {
	err := os.Remove(filepath.Join("static", path))
	if err != nil {
		return err
	}
	return nil
}
