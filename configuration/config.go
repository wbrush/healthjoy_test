package configuration

import (
	"os"

	"bitbucket.org/optiisolutions/go-common/config"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type (
	Config struct {
		config.ServiceParams
		config.DbParams
		config.GCP

		DbMigrationPath string  `json:"db_migration_path"`
	}
)

var cfg Config

func InitConfig(commit, builtAt string) *Config {
	//  load space .env variables first if available
	filename := "./files/.env"
	if _, err := os.Stat(filename); err == nil {
		_ = godotenv.Load(filename)
	}

	cfg := &Config{
		ServiceParams: config.ServiceParams{Environment: "local",
			Host:     "",
			Port:     "8000",
			LogLevel: "debug",
		},

		DbParams: config.DbParams{
			Host:     "",
			Port:     "5432",
			User:     "",
			Password: "",
			Database: "",
		},
		DbMigrationPath: "./dao/postgres",
	}

	err := cfg.LoadEnvVariables(cfg, commit, builtAt)
	if err != nil {
		logrus.Fatalf("cannot load config: %s", err.Error())
	}

	err = cfg.ConfigureLogger()
	if err != nil {
		logrus.Fatalf("cannot ConfigureLogger: %s", err.Error())
	}

	SetCurrentCfg(*cfg)

	return cfg
}

func GetCurrentCfg() Config {
	return cfg
}

func SetCurrentCfg(c Config) {
	cfg = c
}
