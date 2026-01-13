package database

import (
	"fmt"
	"rest-fiber/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetStandalone(env config.Env, logger *DBLogger) (*gorm.DB, error) {
	sslMode := env.DatabaseSSLMode
	if sslMode == "" {
		sslMode = "disable"
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		env.DatabaseHost,
		env.DatabaseUser,
		env.DatabasePassword,
		env.DatabaseName,
		env.DatabasePort,
		sslMode,
	)

	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger,
	})
}