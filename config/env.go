package config

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type DatabaseConfig struct {
	DatabaseURL      string `mapstructure:"DATABASE_URL" validate:"omitempty"`
	DatabaseHost     string `mapstructure:"DATABASE_HOST" validate:"omitempty"`
	DatabaseUser     string `mapstructure:"DATABASE_USER" validate:"omitempty"`
	DatabasePassword string `mapstructure:"DATABASE_PASSWORD" validate:"omitempty"`
	DatabaseName     string `mapstructure:"DATABASE_NAME" validate:"omitempty"`
	DatabasePort     int    `mapstructure:"DATABASE_PORT" validate:"omitempty"`
	DatabaseSSLMode  string `mapstructure:"DATABASE_SSL_MODE" validate:"omitempty"`
}

type JWTConfig struct {
	JWTAccessSecret       string `mapstructure:"JWT_ACCESS_SECRET" validate:"omitempty"`
	JWTRefreshSecret      string `mapstructure:"JWT_REFRESH_SECRET" validate:"omitempty"`
	JWTVerificationSecret string `mapstructure:"JWT_VERIFICATION_SECRET" validate:"omitempty"`
}

type SMTPConfig struct {
	SMTPHost     string `mapstructure:"SMTP_HOST" validate:"omitempty"`
	SMTPPort     string `mapstructure:"SMTP_PORT" validate:"omitempty"`
	SMTPSender   string `mapstructure:"SMTP_SENDER" validate:"omitempty"`
	SMTPEmail    string `mapstructure:"SMTP_EMAIL" validate:"omitempty"`
	SMTPPassword string `mapstructure:"SMTP_PASSWORD" validate:"omitempty"`
}

type RedisConfig struct {
	RedisHost     string `mapstructure:"REDIS_HOST" validate:"omitempty"`
	RedisPort     string `mapstructure:"REDIS_PORT" validate:"omitempty"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD" validate:"omitempty"`
}

type AppEnv string

const (
	Development AppEnv = "development"
	Local       AppEnv = "local"
	Production  AppEnv = "production"
)

type AppConfig struct {
	AppEnv      AppEnv `mapstructure:"APP_ENV" validate:"omitempty"`
	AppAddr     string `mapstructure:"APP_ADDR" validate:"omitempty"`
	FrontendURL string `mapstructure:"FRONTEND_URL" validate:"omitempty"`
}

type Env struct {
	AppConfig
	DatabaseConfig
	JWTConfig
	SMTPConfig
	RedisConfig
}

func NewEnvs() (Env, error) {
	viper.Reset() 
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) 
 
	viper.SetDefault("APP_ENV", "local")

	appEnv := viper.GetString("APP_ENV")
 
	switch appEnv {
	case string(Development), string(Local):
		viper.SetConfigFile(".env.local")
	case string(Production):
		viper.SetConfigFile(".env")
	default:
		viper.SetConfigFile(".env")
	}

	viper.SetConfigType("env")
 
	if err := viper.ReadInConfig(); err != nil {
		 
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return Env{}, err
		}
	}

	var env Env
	if err := viper.Unmarshal(&env); err != nil {
		return Env{}, err
	}

	
	validate := validator.New()
	if err := validate.Struct(env); err != nil {
		return Env{}, err
	}

	return env, nil
}

