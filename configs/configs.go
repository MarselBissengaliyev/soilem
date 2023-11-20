package configs

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Config struct {
	Port     string
	Database struct {
		Username string
		Password string
		Host     string
		Port     string
		DBName   string
	}
	Twilio struct {
		AccountSid string
		AuthToken  string
		FromNumber string
	}
	Gmail struct {
		Password string
		Username string
		Host     string
		Port     int
	}
}

func InitConfig(path string, ctype string, name string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigType(ctype)
	viper.SetConfigName(name)

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "Unable to init config")
	}

	return &Config{
		Port: viper.GetString("port"),
	}, nil
}
