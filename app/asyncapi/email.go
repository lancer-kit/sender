package asyncapi

import (
	"encoding/json"

	"github.com/lancer-kit/armory/natsx"
	"github.com/lancer-kit/sender/config"
	"github.com/lancer-kit/sender/models/email"
	emailp "github.com/lancer-kit/sender/repo/providers/email"
	"github.com/lancer-kit/sender/repo/providers/email/mailgun"
	"github.com/lancer-kit/sender/repo/providers/email/sendgrid"
	"github.com/lancer-kit/sender/repo/providers/email/smtp"
	"github.com/lancer-kit/uwe/v2"
	"github.com/nats-io/go-nats"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type EmailService struct {
	cfg    *config.Cfg
	logger *logrus.Entry
	sender emailp.Sender
}

func NewEmail(cfg *config.Cfg, logger *logrus.Entry) uwe.Worker {
	return &EmailService{
		logger: logger.WithField("worker", config.WorkerAsyncAPIEmail),
		cfg:    cfg,
	}
}

func (s *EmailService) Init() error {
	if s.initEmailProvider(); s.sender == nil {
		return errors.New("email providers does not set")
	}

	return nil
}

func (s *EmailService) Run(ctx uwe.Context) error {
	bus := make(chan *nats.Msg)
	sub, err := natsx.Subscribe(email.Topic, bus)
	if err != nil {
		return errors.Wrap(err, "unable to open subscription")
	}
	defer func() {
		if err = sub.Unsubscribe(); err != nil {
			s.logger.WithError(err).Info("unable to unsubscribe")
		}
	}()

	s.logger.Info("Starting Async-API-Email Service")

	for {
		select {
		case msg := <-bus:
			if msg == nil {
				continue
			}
			s.logger.Debug("got new email message")

			if err = s.processMsg(msg.Data); err != nil {
				s.logger.WithError(err).Error("message processing failed")
				continue
			}

			s.logger.Debug("email was sent")

		case <-ctx.Done():
			s.logger.Info("Async-API-Email gracefully stopped")
			return nil
		}
	}
}

func (s *EmailService) initEmailProvider() {
	switch {
	case s.cfg.Providers.Email.SMTP.Available:
		s.sender = smtp.New(s.cfg.Providers.Email.SMTP)
	case s.cfg.Providers.Email.Sendgrid.Available:
		s.sender = sendgrid.New(s.cfg.Providers.Email.Sendgrid)
	case s.cfg.Providers.Email.Mailgun.Available:
		s.sender = mailgun.New(s.cfg.Providers.Email.Mailgun)
	}
}

func (s *EmailService) processMsg(data []byte) (err error) {
	msg := new(email.Message)

	if err = json.Unmarshal(data, msg); err != nil {
		return errors.Wrap(err, "invalid message data")
	}

	if err = msg.Validate(); err != nil {
		return errors.Wrap(err, "json validation failed")
	}

	if err = s.sender.SendEmail(msg); err != nil {
		return errors.Wrap(err, "message sending failed")
	}

	return err
}
