//go:build debug

package env

import (
	"path/filepath"
	"runtime"
)

func GetExecDir() (string, error) {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(b), "../"), nil
}
