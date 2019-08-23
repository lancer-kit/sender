package app

import (
	"github.com/lancer-kit/sender/config"
	"github.com/pkg/errors"
)

type WorkerList map[string]Worker

func (w WorkerList) InitAll() error {
	var err error

	for name, worker := range w {
		if err = worker.Init(); err != nil {
			return errors.Wrapf(err, "unable to initialize %s worker", name)
		}
	}
	return nil
}

func (w WorkerList) RunAll(available config.Workers, errStreams map[string]chan error) {
	for _, name := range available {
		if worker, ok := w[name]; ok {
			go worker.Run(errStreams[name])
		}
	}
}
