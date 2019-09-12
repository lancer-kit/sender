package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/lancer-kit/armory/api/render"
	"github.com/lancer-kit/armory/log"
	"github.com/lancer-kit/sender/config"
	"github.com/lancer-kit/sender/models/sms"
	"github.com/lancer-kit/sender/repo/providers/email"
	"github.com/lancer-kit/sender/repo/providers/email/mailgun"
	"github.com/lancer-kit/sender/repo/providers/email/sendgrid"
	"github.com/lancer-kit/sender/repo/providers/email/smtp"
	smsp "github.com/lancer-kit/sender/repo/providers/sms"
	"github.com/lancer-kit/sender/repo/providers/sms/twilio"
	"github.com/lancer-kit/sender/repo/providers/sms/viber"
	"github.com/lancer-kit/sender/repo/providers/sms/whatsapp"
	"github.com/lancer-kit/uwe/v2"
	"github.com/lancer-kit/uwe/v2/presets/api"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Service struct {
	cfg    *config.Cfg
	logger *logrus.Entry

	emailSender email.Sender
	smsSenders  map[sms.Provider]smsp.Sender
}

func New(cfg *config.Cfg, logger *logrus.Entry) *Service {
	return &Service{
		logger: logger.WithField("worker", config.WorkerAPIServer),
		cfg:    cfg,
	}

}

func (s *Service) Init() error {
	if s.initSMSProviders(); len(s.smsSenders) == 0 {
		return errors.New("sms providers does not set")
	}

	if s.initEmailProvider(); s.emailSender == nil {
		return errors.New("email provider does not set")
	}

	return nil
}

func (s *Service) Run(ctx uwe.Context) error {
	server := api.NewServer(
		api.Config{
			Host: s.cfg.API.Host,
			Port: s.cfg.API.Port,
		},
		s.router(s.cfg.API.APIRequestTimeout),
	)
	s.logger.Info("Starting API Service")
	if err := server.Run(ctx); err != nil {
		return err
	}
	s.logger.Info("API gracefully stopped")
	return nil
}

func (s *Service) initSMSProviders() {
	s.smsSenders = make(map[sms.Provider]smsp.Sender)
	if s.cfg.Providers.SMS.Twilio.Available {
		s.smsSenders[sms.ProviderSMS] = twilio.New(s.cfg.Providers.SMS.Twilio)
	}
	if s.cfg.Providers.SMS.Viber.Available {
		s.smsSenders[sms.ProviderViber] = viber.New(s.cfg.Providers.SMS.Viber)
	}
	if s.cfg.Providers.SMS.Whatsapp.Available {
		s.smsSenders[sms.ProviderSMS] = whatsapp.New(s.cfg.Providers.SMS.Whatsapp)
	}
}

func (s *Service) initEmailProvider() {
	switch {
	case s.cfg.Providers.Email.SMTP.Available:
		s.emailSender = smtp.New(s.cfg.Providers.Email.SMTP)
	case s.cfg.Providers.Email.Sendgrid.Available:
		s.emailSender = sendgrid.New(s.cfg.Providers.Email.Sendgrid)
	case s.cfg.Providers.Email.Mailgun.Available:
		s.emailSender = mailgun.New(s.cfg.Providers.Email.Mailgun)
	}
}

func (s *Service) router(requestTimeout int) http.Handler {
	r := chi.NewRouter()

	r.Use(log.NewRequestLogger(s.logger.Logger))
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	if requestTimeout > 0 {
		r.Use(middleware.Timeout(time.Second * time.Duration(requestTimeout)))
	}

	r.Route("/v1/sender", func(r chi.Router) {
		h := &handler{
			logger:      s.logger,
			smsSenders:  s.smsSenders,
			emailSender: s.emailSender,
		}
		r.Post("/email", h.SendEmail)
		r.Post("/sms", h.SendSms)
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		render.ResultNotFound.Render(w)
	})

	return r
}
