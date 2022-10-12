package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"sync"
)

type Config struct {
	AppDebug bool   `env:"APP_DEBUG" env-default:"true"`
	AppEnv   string `env:"APP_ENV" env-default:"dev"`
	Listen   struct {
		BindIP string `env:"BACKEND_IP" env-default:"10.10.10.1"`
		Port   string `env:"BACKEND_PORT" env-default:"10000"`
	}
	AppConfig struct {
		LogLever string
		Database struct {
			Host     string `env:"DB_HOST"`
			Port     string `env:"DB_PORT"`
			DbName   string `env:"DB_NAME" env-default:"currency"`
			User     string `env:"DB_USER"`
			Password string `env:"DB_PASSWORD"`
		}
	}
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		log.Print("Populate config")

		instance = &Config{}

		cleanenv.ReadConfig("./.env", instance)
		if err := cleanenv.ReadEnv(instance); err != nil {

			description, _ := cleanenv.GetDescription(instance, nil)

			log.Print(description)
			log.Fatal(err)
		}
	})

	return instance
}
