package config

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Env struct {
	DatabaseURL           string `mapstructure:"DATABASE_URL" validate:"omitempty"`
	DatabaseHost          string `mapstructure:"DATABASE_HOST" validate:"omitempty"`
	DatabaseUser          string `mapstructure:"DATABASE_USER" validate:"omitempty"`
	DatabasePassword      string `mapstructure:"DATABASE_PASSWORD" validate:"omitempty"`
	DatabaseName          string `mapstructure:"DATABASE_NAME" validate:"omitempty"`
	DatabasePort          int    `mapstructure:"DATABASE_PORT" validate:"omitempty"`
	DatabaseSSLMode       string `mapstructure:"DATABASE_SSL_MODE" validate:"omitempty"`
	AppEnv                string `mapstructure:"APP_ENV" validate:"omitempty"`
	AppAddr               string `mapstructure:"APP_ADDR" validate:"omitempty"`
	JWTAccessSecret       string `mapstructure:"JWT_ACCESS_SECRET" validate:"omitempty"`
	JWTRefreshSecret      string `mapstructure:"JWT_REFRESH_SECRET" validate:"omitempty"`
	JWTVerificationSecret string `mapstructure:"JWT_VERIFICATION_SECRET" validate:"omitempty"`
	SMTPHost              string `mapstructure:"SMTP_HOST" validate:"omitempty"`
	SMTPPort              string `mapstructure:"SMTP_PORT" validate:"omitempty"`
	SMTPSender            string `mapstructure:"SMTP_SENDER" validate:"omitempty"`
	SMTPEmail             string `mapstructure:"SMTP_EMAIL" validate:"omitempty"`
	SMTPPassword          string `mapstructure:"SMTP_PASSWORD" validate:"omitempty"`
	RedisHost             string `mapstructure:"REDIS_HOST" validate:"omitempty"`
	RedisPort             string `mapstructure:"REDIS_PORT" validate:"omitempty"`
	RedisPassword         string `mapstructure:"REDIS_PASSWORD" validate:"omitempty"`
	TargetURL             string `mapstructure:"TARGET_URL" validate:"omitempty"`
}

func GetEnvs() (Env, error) {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	appEnv := viper.GetString("APP_ENV")
	
	if appEnv == "development" {
		viper.SetConfigFile(".env.local")
	} else {
		viper.SetConfigFile(".env")
	}

	viper.SetConfigType("env")
	_ = viper.ReadInConfig()

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
