package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env          string  `yaml:"env" env:"ENV" env-required:"true"`
	FontPath     string  `yaml:"font_path" env-required:"true"`
	FontSize     float64 `yaml:"font_size" env-required:"true"`
	MaxWidth     int     `yaml:"max_width" env-required:"true"`
	VideoOptions `yaml:"video_options"`
}

type VideoOptions struct {
	FPS         int     `yaml:"fps" env-default:"30"`
	RandomDelay bool    `yaml:"random_delay" env-default:"false"`
	Delay       float64 `yaml:"delay" env-default:"1"`
	MinDelay    float64 `yaml:"min_delay" env-default:"0.8"`
	MaxDelay    float64 `yaml:"max_delay" env-default:"1.2"`
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
