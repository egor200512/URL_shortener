package configs

import (
	"fmt"

	"github.com/lpernett/godotenv"
)

func LoadConfig(path string) {
	if err := godotenv.Load(path); err != nil {
		panic(fmt.Errorf("failed to load config: %s", err.Error()))
	}
}
