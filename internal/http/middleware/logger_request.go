package middleware

import ( 
	"rest-fiber/internal/infra/infraapp"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func LoggerRequest(l *infraapp.AppLogger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()

		l.WithFields(logrus.Fields{
			"method":  c.Method(),
			"path":    c.OriginalURL(),
			"endpoint": c.Route().Path,
			"status":  c.Response().StatusCode(),
			"latency": time.Since(start).String(),
			"ip":      c.IP(),
		}).Info("http_request")

		return err
	}
}
