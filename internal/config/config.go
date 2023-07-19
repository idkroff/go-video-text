package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string  `yaml:"env" env:"ENV" env-required:"true"`
	FontPath string  `yaml:"font_path" env-required:"true"`
	FontSize float64 `yaml:"font_size" env-required:"true"`
	MaxWidth int     `yaml:"max_width" env-required:"true"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./configs/local.yaml"
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist on %s", configPath)
	}

	var config Config
	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &config
}
