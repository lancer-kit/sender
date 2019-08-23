package sendgrid

import (
	"fmt"

	"github.com/lancer-kit/sender/config"
	"github.com/lancer-kit/sender/models/email"
	"github.com/pkg/errors"
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type sender struct {
	cfg *config.Sendgrid
}

func New(cfg *config.Sendgrid) *sender {
	return &sender{
		cfg: cfg,
	}
}

func (s *sender) SendEmail(email *email.Message) error {
	request := s.createRequest(email)
	response, err := sendgrid.API(request)
	if err != nil {
		return errors.Wrap(err, "cant send message")
	}
	if response.StatusCode >= 400 {
		return fmt.Errorf("sendgrid return status code %d", response.StatusCode)
	}

	return nil
}

func (s *sender) createRequest(email *email.Message) rest.Request {

	p := mail.NewPersonalization()
	p.AddTos(mail.NewEmail(email.To, email.To))
	p.SetSubstitution("{{.Text}}", email.Body)
	p.Subject = email.Subject

	m := mail.NewV3Mail()
	m.SetFrom(mail.NewEmail(s.cfg.SenderName, s.cfg.Sender))
	m.AddContent(mail.NewContent("text/html", email.Body))
	m.AddPersonalizations(p)

	request := sendgrid.GetRequest(s.cfg.PrivateKey, s.cfg.Endpoint, s.cfg.Host)
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)

	return request
}
