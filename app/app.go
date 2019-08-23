package app

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/lancer-kit/armory/log"
	"github.com/lancer-kit/armory/natsx"
	"github.com/lancer-kit/sender/app/api"
	"github.com/lancer-kit/sender/app/asyncapi"
	"github.com/lancer-kit/sender/config"
	"github.com/sirupsen/logrus"
)

type App struct {
	cfg     *config.Cfg
	ctx     context.Context
	cancel  context.CancelFunc
	wg      *sync.WaitGroup
	logger  *logrus.Entry
	workers WorkerList
	errors  map[string]chan error
}

func New(cfg *config.Cfg) *App {
	ctx, cancel := context.WithCancel(context.Background())
	logger := log.Default.WithField("app", config.ServiceName)
	workers := WorkerList{
		config.WorkerAPIServer:     api.New(ctx, cfg, logger),
		config.WorkerAsyncAPIEmail: asyncapi.NewEmail(ctx, cfg, logger),
		config.WorkerAsyncAPISms:   asyncapi.NewSms(ctx, cfg, logger),
	}
	return &App{
		cfg:     cfg,
		workers: workers,
		logger:  logger,
		cancel:  cancel,
		ctx:     ctx,
		wg:      new(sync.WaitGroup),
	}
}

func (a *App) Run() {
	a.errors = make(map[string]chan error)
	for name := range a.workers {
		a.errors[name] = make(chan error)
	}

	natsx.SetConfig(a.cfg.NATS)

	a.workers.RunAll(a.cfg.Workers, a.errors)

	go a.checkWorkerErrors(a.errors)

	a.GracefulStop()
}

func (a *App) GracefulStop() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	<-signalChan

	a.logger.Logger.Debugln("Begin graceful shutdown...")
	a.cancel()
	a.wg.Wait()
	a.logger.Logger.Debugln("Service successfully stopped")
}

func (a App) Context() context.Context {
	return a.ctx
}

func (a App) Logger() *logrus.Entry {
	return a.logger
}

func (a *App) WAdd(delta int) {
	a.wg.Add(delta)
}

func (a *App) WDone() {
	a.wg.Done()
}

func (a *App) checkWorkerErrors(errors map[string]chan error) {
	for name, err := range errors {
		go a.logErrors(name, err)
	}
	<-a.Context().Done()
}

func (a *App) logErrors(name string, errStream chan error) {
	for err := range errStream {
		if err != nil {
			a.Logger().WithError(err).
				WithField("worker", name).
				Errorln("error from worker")
		}
	}
}
