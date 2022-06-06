package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	App App
	DB DB
	Auth Auth
}

type App struct {
	Endpoint   string  `mapstructure:"endpoint"   validate:"required,url"`
	MinRating  float64 `mapstructure:"minRating"  validate:"required,numeric,max=5"`
	MaxWorkers int     `mapstructure:"maxWorkers" validate:"required,numeric,min=1"`
	Latitude   float64 `mapstructure:"latitude"   validate:"required,latitude"`
	Longitude  float64 `mapstructure:"longitude"  validate:"required,longitude"`
}

type DB struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

type Auth struct {
	Name string
	Password string
}

func SetUpConfig() (*Config, error) {
	if err := setUpViper(); err != nil {
		return nil, err
	}

	var cfg Config

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	if err := parseEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func setUpViper() error {
	viper.AddConfigPath("./configs")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}

func parseEnv(cfg *Config) error {
	if err := viper.BindEnv("DB_PASSWORD"); err != nil {
		return err
	}

	if err := viper.BindEnv("AUTH_USERNAME"); err != nil {
		return err
	}

	if err := viper.BindEnv("AUTH_PASSWORD"); err != nil {
		return err
	}

	cfg.DB.Password = viper.GetString("DB_PASSWORD")
	cfg.Auth.Name = viper.GetString("AUTH_USERNAME")
	cfg.Auth.Password = viper.GetString("AUTH_PASSWORD")

	return nil
}


