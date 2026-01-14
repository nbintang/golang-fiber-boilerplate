package main

import (
	"context"
	"rest-fiber/config"
	"rest-fiber/internal/infra/database"
	"rest-fiber/internal/user"

	"gorm.io/gorm"
)

func InitMigrate(ctx context.Context) error {
	dbLogger := database.NewLogger()
	env, err := config.NewEnvs()
	if err != nil {
		return err
	}
	db, err := database.GetStandalone(env, dbLogger)
	if err != nil {
		return err
	}

	if err := createEnums(ctx, db); err != nil {
		return err
	}

	return db.WithContext(ctx).Debug().AutoMigrate(
		&user.User{},
		// add more models / entities in here
	)
}

func createEnums(ctx context.Context, db *gorm.DB) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Debug().Exec(`
		DO $$ BEGIN
			CREATE TYPE role_type AS ENUM ('ADMIN', 'MEMBER');
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;`).Error; err != nil {
			return err
		}

		if err := tx.Debug().Exec(`
		DO $$ BEGIN
			CREATE TYPE status_type AS ENUM ('PUBLISHED', 'DRAFT');
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;`).Error; err != nil {
			return err
		}
		return nil
	})
}
