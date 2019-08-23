package asyncapi

import (
	"context"
	"encoding/json"

	"github.com/lancer-kit/armory/natsx"
	"github.com/lancer-kit/sender/config"
	"github.com/lancer-kit/sender/models/sms"
	smsp "github.com/lancer-kit/sender/repo/providers/sms"
	"github.com/lancer-kit/sender/repo/providers/sms/twilio"
	"github.com/lancer-kit/sender/repo/providers/sms/viber"
	"github.com/lancer-kit/sender/repo/providers/sms/whatsapp"
	"github.com/nats-io/go-nats"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type SmsService struct {
	ctx       context.Context
	cfg       *config.Cfg
	logger    *logrus.Entry
	providers map[sms.Provider]smsp.Sender
}

func NewSms(ctx context.Context, cfg *config.Cfg, logger *logrus.Entry) *SmsService {
	return &SmsService{
		logger: logger.WithField("worker", config.WorkerAsyncAPISms),
		ctx:    ctx,
		cfg:    cfg,
	}
}

func (s *SmsService) Init() error {

	return nil
}

func (s *SmsService) Run(errChan chan<- error) {
	bus := make(chan *nats.Msg)
	sub, err := natsx.Subscribe(sms.Topic, bus)
	if err != nil {
		errChan <- errors.Wrap(err, "unable to open subscription")
		return
	}
	defer func() {
		if err = sub.Unsubscribe(); err != nil {
			errChan <- errors.Wrap(err, "unable to unsubscribe")
		}
	}()

	s.logger.Info("Starting Async-API-SMS Service")

	for {
		select {
		case msg := <-bus:
			if msg == nil {
				continue
			}
			s.logger.Debug("got new sms message")

			if err = s.processMsg(msg.Data); err != nil {
				errChan <- errors.Wrap(err, "msg processing failed")
				continue
			}

			s.logger.Debug("sms was sent")

		case <-s.ctx.Done():
			s.logger.Info("async api sms gracefully stopped")
			return
		}
	}
}

func (s *SmsService) processMsg(data []byte) (err error) {
	msg := new(sms.Message)

	if err = json.Unmarshal(data, msg); err != nil {
		return errors.Wrap(err, "invalid message data")
	}

	err = msg.Validate()
	if err != nil {
		return errors.Wrap(err, "json validation failed")
	}

	if _, ok := s.providers[msg.Provider]; !ok {
		return errors.Wrap(err, "cannot send sms, unavailable provider")
	}

	if err = s.providers[msg.Provider].SendSms(msg); err != nil {
		return errors.Wrap(err, "cant send sms")
	}

	return nil
}

func (s *SmsService) initSMSProviders() {
	s.providers = make(map[sms.Provider]smsp.Sender)
	if s.cfg.Providers.SMS.Twilio.Available {
		s.providers[sms.ProviderSMS] = twilio.New(s.cfg.Providers.SMS.Twilio)
	}
	if s.cfg.Providers.SMS.Viber.Available {
		s.providers[sms.ProviderViber] = viber.New(s.cfg.Providers.SMS.Viber)
	}
	if s.cfg.Providers.SMS.Whatsapp.Available {
		s.providers[sms.ProviderSMS] = whatsapp.New(s.cfg.Providers.SMS.Whatsapp)
	}
}
