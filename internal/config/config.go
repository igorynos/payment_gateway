package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `yaml:"env" env-default:"local"`
	Storage    DATABASEURL
	HTTPServer `yaml:"http_server"`
}

type HTTPServer struct {
	Address      string        `yaml:"address" env-default:"localhost"`
	Port         int           `yaml:"port" env-default:"9000"`
	Timeout      time.Duration `yaml:"timeout" env-default:"3s"`
	Idle_timeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type DATABASEURL struct {
	Name     string `env:"DATABASE_NAME" env-required:"true"`
	User     string `env:"DATABASE_USER" env-required:"true"`
	Password string `env:"DATABASE_PASSWORD" env-required:"true"`
	Host     string `env:"DATABASE_HOST" env-required:"true"`
	Port     string `env:"DATABASE_PORT" env-required:"true"`
}

func (d DATABASEURL) URL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		d.User,
		d.Password,
		d.Host,
		d.Port,
		d.Name,
	)
}

func LoadConfig() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("Config is no set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist: %s", configPath)
	}

	var config Config

	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		log.Fatalf("Cannot read config %s", err)
	}

	return &config
}
