package configs

import (
	"github.com/MarselBissengaliyev/soilem/internal/repo"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Config struct {
	Port     string
	Postgres repo.PostgresConfig
	Twilio   struct {
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
		Postgres: repo.PostgresConfig{
			UserName: viper.GetString("postgres.username"),
			Password: viper.GetString("postgres.password"),
			Host:     viper.GetString("postgres.host"),
			Port:     viper.GetString("postgres.port"),
			DBName:   viper.GetString("postgres.dbname"),
			SSLmode:  viper.GetString("postgres.sslmode"),
		},
	}, nil
}
