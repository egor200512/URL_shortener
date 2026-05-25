package configs

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/lpernett/godotenv"
)

func LoadConfig(path string) {
	if err := godotenv.Load(path); err != nil {
		fallbackPath, fallbackErr := findEnvFile()
		if fallbackErr != nil {
			panic(fmt.Errorf("failed to load config: %s", err.Error()))
		}
		if err := godotenv.Load(fallbackPath); err != nil {
			panic(fmt.Errorf("failed to load config: %s", err.Error()))
		}
	}
}

func findEnvFile() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		path := filepath.Join(dir, ".env")
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", os.ErrNotExist
		}
		dir = parent
	}
}
