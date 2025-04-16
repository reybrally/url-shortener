package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env        string `yaml:"env" env:"ENV" env-default:"local"`
	Storage    string `yaml:"storage_path" env-required:"true"`
	HTTPServer `yaml:"http_server"`
}
type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
	User        string        `yaml:"user" env-required:"true"`
	Password    string        `yaml:"password" env-required:"true" env:"HTTP-SERVER_PASSWORD"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH") // получаем путь из переменного окружения
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable not set")
	}
	// check of file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("CONFIG_PATH does not exist")
	}
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal("cannot read config")
	}
	return &cfg
}
