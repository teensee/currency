package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"sync"
)

type Config struct {
	AppDebug bool   `env:"APP_DEBUG" env-default:"true"`
	AppEnv   string `env:"APP_ENV" env-default:"dev"`
	Vault    struct {
		UseVault    bool   `env:"USE_VAULT" env-default:"false"`
		MountPath   string `env:"VAULT_MOUNT_PATH"`
		SecretPath  string `env:"VAULT_SECRET_PATH"`
		VaultHost   string `env:"VAULT_HOST"`
		VaultApiKey string `env:"VAULT_API_KEY"`
	}
	Listen struct {
		BindIP string `env:"BACKEND_IP" env-default:"10.10.10.1"`
		Port   string `env:"BACKEND_PORT" env-default:"10000"`
	}
	AppConfig struct {
		LogLever              string
		SyncRatesAfterStartup bool `env:"SYNC_RATES_AFTER_STARTUP" env-default:"false"`
		Database              struct {
			Host     string `env:"DB_HOST"`
			Port     string `env:"DB_PORT"`
			DbName   string `env:"DB_NAME"`
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
		_ = cleanenv.ReadConfig("./.env", instance)

		if err := cleanenv.ReadEnv(instance); err != nil {

			description, _ := cleanenv.GetDescription(instance, nil)

			log.Print(description)
			log.Fatal(err)
		}
	})

	return instance
}

func (c *Config) GetDsn() string {
	db := c.AppConfig.Database
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Moscow", db.Host, db.User, db.Password, db.DbName, db.Port)
}
