package email

import (
	"context"
	"fmt"
	"net/smtp"
	"rest-fiber/config"
	"rest-fiber/internal/infra/infraapp"
	"strings"
)

type serviceImpl struct {
	env    config.Env
	logger *infraapp.AppLogger
}

func NewService(env config.Env, logger *infraapp.AppLogger) EmailService {
	return &serviceImpl{env, logger}
}
func (s *serviceImpl) SendEmail(ctx context.Context, params Params) error {
	done := make(chan error, 1)
	go func() {
		addr := fmt.Sprintf("%s:%s", s.env.SMTPHost, s.env.SMTPPort)
		auth := smtp.PlainAuth(
			"",
			s.env.SMTPEmail,
			s.env.SMTPPassword,
			s.env.SMTPHost,
		)
		headers := []string{
			"From: " + s.env.SMTPEmail,
			"To: " + params.Reciever.Email,
			"Subject: " + params.Subject,
			"MIME-Version: 1.0",
			"Content-Type: text/plain; charset=\"utf-8\"",
			"",
			params.Message,
		}

		msg := []byte(strings.Join(headers, "\r\n"))
		done <- smtp.SendMail(addr, auth, s.env.SMTPEmail, []string{params.Reciever.Email}, msg)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}
