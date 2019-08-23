package smtp

import (
	"fmt"
	"net/smtp"

	"github.com/lancer-kit/sender/config"
	"github.com/lancer-kit/sender/models/email"
	"github.com/pkg/errors"
)

const (
	MIME = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	msgF = "From: %s\n" + "To: %s\n" + "Subject: %s\n" + MIME + "%s"
)

type sender struct {
	cfg *config.SMTP
}

func New(cfg *config.SMTP) *sender {
	return &sender{
		cfg: cfg,
	}
}

func (s *sender) SendEmail(email *email.Message) error {
	msg := fmt.Sprintf(msgF, s.cfg.From, email.To, email.Subject, email.Body)

	err := smtp.SendMail(
		fmt.Sprintf("%s:%s", s.cfg.Host, s.cfg.Port),
		smtp.PlainAuth("", s.cfg.From, s.cfg.Password, s.cfg.Host),
		s.cfg.From,
		[]string{email.To},
		[]byte(msg),
	)

	if err != nil {
		return errors.Wrap(err, "cant send message")
	}

	return nil
}
