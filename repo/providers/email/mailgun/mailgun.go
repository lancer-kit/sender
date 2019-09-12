package mailgun

import (
	"github.com/lancer-kit/sender/config"
	"github.com/lancer-kit/sender/models/email"
	emailp "github.com/lancer-kit/sender/repo/providers/email"
	"github.com/pkg/errors"
	"gopkg.in/mailgun/mailgun-go.v1"
)

type sender struct {
	cfg  *config.Mailgun
	conn mailgun.Mailgun
}

func New(cfg *config.Mailgun) emailp.Sender {
	return &sender{
		cfg:  cfg,
		conn: mailgun.NewMailgun(cfg.Domain, cfg.PrivateKey, cfg.PublicKey),
	}
}

func (m *sender) SendEmail(email *email.Message) error {

	message := m.conn.NewMessage(
		m.cfg.Sender,
		email.Subject,
		email.Body,
		email.To,
	)
	message.SetHtml(email.Body)

	_, _, err := m.conn.Send(message)
	if err != nil {
		return errors.Wrap(err, "cant send message")
	}

	return nil
}
