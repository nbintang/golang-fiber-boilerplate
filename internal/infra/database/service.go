package database

import (
	"rest-fiber/config"
	"gorm.io/gorm"
)


func NewService(env config.Env, logger *DBLogger) (*gorm.DB, error) {
	return GetStandalone(env, logger)
}
