package utils

import (
	"os"
	"path/filepath"
)

func DirPath(path ...string) (string, error) {
	pathDir := filepath.Join(path...)
	if info, errDir := os.Stat(pathDir); errDir != nil || !info.IsDir() {
		if errDir = os.Mkdir(pathDir, os.ModePerm); errDir != nil {
			return "", errDir
		}
	}
	return pathDir, nil
}
