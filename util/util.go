package util

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func DirExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return info.IsDir(), nil
}

func AllFilesInDir(path string, extension string) ([]string, error) {
	var xmlFiles []string

	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error: accessing path of %s %w\n", path, err)
		}

		if !d.IsDir() && strings.HasSuffix(strings.ToLower(d.Name()), extension) {
			xmlFiles = append(xmlFiles, path)
		}

		return nil
	})

	return xmlFiles, err
}

func FindMatchingDirs(patternPath string) ([]string, error) {
	matches, err := filepath.Glob(patternPath)
	if err != nil {
		return nil, err
	}

	var dirs []string
	for _, match := range matches {
		info, err := os.Stat(match)
		if err != nil {
			return nil, err
		}

		if info.IsDir() {
			dirs = append(dirs, match)
		}
	}

	return dirs, nil
}

func RemoveAtIndex(s [][]byte, i int) [][]byte {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func ReadAllFiles(files []string) ([][]byte, error) {
	var contents [][]byte

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}

		contents = append(contents, data)
	}

	return contents, nil
}
