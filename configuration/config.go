package configuration

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type (
	Config struct {
		Commit  string `json:"version"`
		BuiltAt string `json:"built_at"`

		Host string `json:"host"`
		Port string `json:"port"`

		isLoaded bool `json:"is_loaded"`
	}
)

var cfg Config

func InitConfig(commit, builtAt string) *Config {
	ConfigureLogger()

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

	hostStr := os.Getenv("HOST")
	if len(hostStr) > 0 {
		cfg.Host = hostStr
	} else {
		logrus.Warnf("HOST env variable not found. Using Default!")
	}

	portStr := os.Getenv("PORT")
	if len(portStr) > 0 {
		_, err := strconv.Atoi(portStr)
		if err != nil {
			logrus.Errorf("Error converting the port env variable to integer: %s", err.Error())
		} else {
			cfg.Port = portStr
		}
	} else {
		logrus.Warnf("PORT env variable not found. Using Default!")
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
