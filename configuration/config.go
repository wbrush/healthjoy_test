package configuration

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type (
	Config struct {
		Commit  string
		BuiltAt string

		Host string
		Port string
	}
)

var cfg Config

func InitConfig(commit, builtAt string) *Config {
	//  load space .env variables first if available
	filename := "./files/.env"
	if _, err := os.Stat(filename); err == nil {
		logrus.Debugf("loading environmental variables")
		_ = godotenv.Load(filename)
	}

	cfg := &Config{
		Commit:  commit,
		BuiltAt: builtAt,
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
