package envx

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	loadOnce sync.Once
	loadErr  error
)

func Load() error {
	loadOnce.Do(func() {
		path, err := findDotEnv()
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return
			}
			loadErr = err
			return
		}

		loadErr = loadFile(path)
	})

	return loadErr
}

func findDotEnv() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		path := filepath.Join(dir, ".env")
		info, err := os.Stat(path)
		if err == nil && !info.IsDir() {
			return path, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", os.ErrNotExist
}

func loadFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "export ") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		value = strings.Trim(value, `"'`)

		if key == "" {
			continue
		}

		if os.Getenv(key) != "" {
			continue
		}

		if err := os.Setenv(key, value); err != nil {
			return err
		}
	}

	return scanner.Err()
}
