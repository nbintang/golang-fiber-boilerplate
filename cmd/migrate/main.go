package main

import (
	"context"
	"rest-fiber/pkg/env"

	"github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()
	env.Load()
	if err := InitMigrate(ctx); err != nil {
		logrus.Warnf("Migration failed: %v", err)
	}
	logrus.Println("Migration Succeed")
}
