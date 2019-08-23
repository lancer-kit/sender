package main

import (
	"os"

	"github.com/lancer-kit/armory/log"
	"github.com/lancer-kit/sender/app"
	"github.com/lancer-kit/sender/config"
	"github.com/urfave/cli"
)

const (
	cfgFlag = "config"

	serveCommand = "serve"
)

func main() {
	a := cli.NewApp()
	a.Name = config.ServiceName
	a.Commands = append(
		a.Commands,
		serve(),
	)
	a.Run(os.Args)
}

func serve() cli.Command {
	var cfgFlag = []cli.Flag{
		cli.StringFlag{
			Name:  cfgFlag,
			Value: "./config.yaml",
		},
	}
	var serveCommand = cli.Command{
		Name:   serveCommand,
		Usage:  "starts " + config.ServiceName + " workers",
		Flags:  cfgFlag,
		Action: serveAction,
	}
	return serveCommand
}

func serveAction(c *cli.Context) error {
	cfgPath := c.String(cfgFlag)
	cfg, err := config.Config(cfgPath)
	if err != nil {
		log.Default.WithError(err).Fatal("cannot get configs")
	}

	app.New(cfg).Run()

	return nil
}
