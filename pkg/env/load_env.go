package env

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func Load(filenames ...string) {
	if len(filenames) == 0 {
		if err := godotenv.Load(".env"); err != nil {
			logrus.Warn("No .env file found, using environment variables")
		}
		return
	}

	if err := godotenv.Load(filenames...); err != nil {
		logrus.Warn("No .env file found, using environment variables")
	}
}
