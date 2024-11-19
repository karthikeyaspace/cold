package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Email string
	Pass  string
	Port  string
	Host  string
	Batch int
	AiKey string
	Delay int
	ResumePath string
	DataPath string
}

func LoadEnv() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return &Config{}, fmt.Errorf("error loading .env file")
	}

	batchSize := 2
	if val := os.Getenv("BATCH"); val != "" {
		batchSize, _ = strconv.Atoi(val)
	}

	delay, err := strconv.Atoi(os.Getenv("DELAY"))
	if err != nil {
		return &Config{}, fmt.Errorf("error parsing delay: %v, not present in env", err)
	}

	cfg := &Config{
		Email: os.Getenv("EMAIL"),
		Pass:  os.Getenv("PASS"),
		Port:  "587",
		Host:  "smtp.gmail.com",
		Batch: batchSize,
		AiKey: os.Getenv("AI_KEY"),
		Delay: delay,
		ResumePath: os.Getenv("data/resume.pdf"),
		DataPath: os.Getenv("data/data.xlsx"), 
	}

	if cfg.Email == "" || cfg.Pass == "" {
		return &Config{}, fmt.Errorf("email and password are required")
	}

	return cfg, nil

}
