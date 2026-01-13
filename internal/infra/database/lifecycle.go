package database

import (
	"context"
	"time"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

func RegisterLifecycle(lc fx.Lifecycle, db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return sqlDB.PingContext(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return sqlDB.Close()
		},
	})

	return nil
}
