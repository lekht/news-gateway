package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App    `yaml:"app"`
		API    `yaml:"api"`
		Server `yaml:"server"`
	}

	App struct {
		Name   string `yaml:"name"`
		Versin string `yaml:"version"`
	}

	API struct {
		NewsAddr      string `yaml:"newsAddr"`
		CommentsAddr  string `yaml:"commentsAddr"`
		FormatterAddr string `yaml:"formatterAddr"`
	}

	Server struct {
		Listen string `yaml:"listen"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}
	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config - NewConfig - ReadConfig error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
