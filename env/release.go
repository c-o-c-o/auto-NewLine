//go:build !debug

package env

import (
	"os"
	"path/filepath"
)

func GetExecDir() (string, error) {
	p, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(p), nil
}
