package main

import (
	"os"
	"path/filepath"
)

func appPath() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", err
	}

	return filepath.Dir(execPath), nil
}
