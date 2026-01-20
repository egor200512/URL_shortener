package tests

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/lpernett/godotenv"
)

func TestMain(m *testing.M) {
	if err := godotenv.Load(filepath.Join("..", "..", "..", "..", "..", "..", ".env")); err != nil {
		log.Fatalf("failed to load .env for testing: %s", err.Error())
	}

	exitCode := m.Run()

	os.Exit(exitCode)
}
