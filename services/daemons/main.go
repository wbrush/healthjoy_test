package daemons

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/wbrush/healthjoy_test/configuration"
)

type Daemons struct {
	config  *configuration.Config
	workers []helpers.Daemon
}

func NewDaemons(cfg *configuration.Config) Daemons {
	return Daemons{
		config: cfg,
		workers: []helpers.Daemon{
			NewExampleTicker(configuration.ExampleTickerIntervalSec),
		},
	}
}

func (d Daemons) GracefulStop(ctx context.Context) error {
	if len(d.workers) > 0 {
		for _, i := range d.workers {
			logrus.Warnf("Stopping daemon %s...", i.Title())
			c, _ := context.WithCancel(ctx)
			if err := i.GracefulStop(c); err != nil {
				return fmt.Errorf("can't stop daemon %s: %s", d.Title(), err.Error())
			}
		}
	}
	return nil
}

// Title returns the title.
func (d Daemons) Title() string {
	return "Daemons"
}

// Run subscribes to the queue.
func (d Daemons) Run() {
	if len(d.workers) > 0 {
		for _, i := range d.workers {
			logrus.Infof("Starting daemon %s", i.Title())
			go i.Run()
		}

		select {} //lock routine execution
	}
}
