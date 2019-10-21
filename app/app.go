package app

import (
	"github.com/lancer-kit/armory/log"
	"github.com/lancer-kit/armory/natsx"
	"github.com/lancer-kit/sender/app/api"
	"github.com/lancer-kit/sender/app/asyncapi"
	"github.com/lancer-kit/sender/config"
	"github.com/lancer-kit/uwe/v2"
	"github.com/sirupsen/logrus"
)

type App struct {
	cfg    *config.Cfg
	logger *logrus.Entry
}

func New(cfg *config.Cfg) *App {
	a := new(App)
	a.cfg = cfg
	a.logger = log.Default.WithField("app", config.ServiceName)
	return a
}

func (a *App) Run() {
	natsx.SetConfig(a.cfg.NATS)

	chief := uwe.NewChief()
	workers := a.workers()
	for _, name := range a.cfg.Workers {
		chief.AddWorker(uwe.WorkerName(name), workers[uwe.WorkerName(name)])
	}
	chief.UseDefaultRecover()
	chief.SetEventHandler(a.eventHandler)
	chief.Run()
}

func (a *App) workers() map[uwe.WorkerName]uwe.Worker {
	return map[uwe.WorkerName]uwe.Worker{
		config.WorkerAPIServer: api.New(
			a.cfg,
			a.logger.WithField("worker", config.WorkerAPIServer),
		),
		config.WorkerAsyncAPIEmail: asyncapi.NewEmail(
			a.cfg,
			a.logger.WithField("worker", config.WorkerAsyncAPIEmail),
		),
		config.WorkerAsyncAPISms: asyncapi.NewSms(
			a.cfg,
			a.logger.WithField("worker", config.WorkerAsyncAPISms),
		),
	}
}

func (a *App) eventHandler(event uwe.Event) {
	logger := a.logger.WithField("level", event.Level).WithField("worker", event.Worker)
	if err := event.ToError(); err != nil {
		logger.WithError(err).Error("error in chief")
		return
	}

	logger.WithField("message", event.Message).Info("chief log")
}
