package asyncapi

import (
	"context"
	"encoding/json"

	"github.com/lancer-kit/sender/config"

	"github.com/lancer-kit/armory/natsx"
	"github.com/lancer-kit/sender/models/email"
	emailp "github.com/lancer-kit/sender/repo/providers/email"
	"github.com/lancer-kit/sender/repo/providers/email/mailgun"
	"github.com/lancer-kit/sender/repo/providers/email/sendgrid"
	"github.com/lancer-kit/sender/repo/providers/email/smtp"
	"github.com/nats-io/go-nats"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type EmailService struct {
	ctx    context.Context
	cfg    *config.Cfg
	logger *logrus.Entry
	sender emailp.Sender
}

func NewEmail(ctx context.Context, cfg *config.Cfg, logger *logrus.Entry) *EmailService {
	return &EmailService{
		logger: logger.WithField("worker", config.WorkerAsyncAPIEmail),
		cfg:    cfg,
		ctx:    ctx,
	}
}

func (s *EmailService) Init() error {
	if s.initEmailProvider(); s.sender == nil {
		return errors.New("sms providers does not set")
	}

	return nil
}

func (s *EmailService) Run(errChan chan<- error) {
	bus := make(chan *nats.Msg)
	sub, err := natsx.Subscribe(email.Topic, bus)
	if err != nil {
		errChan <- errors.Wrap(err, "unable to open subscription")
		return
	}
	defer func() {
		if err = sub.Unsubscribe(); err != nil {
			errChan <- errors.Wrap(err, "unable to unsubscribe")
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
				errChan <- errors.Wrap(err, "msg processing failed")
				continue
			}

			s.logger.Debug("email was sent")

		case <-s.ctx.Done():
			s.logger.Info("email async api gracefully stopped")
			return

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
