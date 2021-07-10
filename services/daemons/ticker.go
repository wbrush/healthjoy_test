package daemons

import (
	"context"
	"log"
	"time"

	"github.com/wbrush/healthjoy_test/configuration"
)

type ExampleTicker struct {
	interval time.Duration

	ticker *time.Ticker
}

func NewExampleTicker(interval time.Duration) *ExampleTicker {
	return &ExampleTicker{
		interval: interval,
	}
}

func (et ExampleTicker) Title() string {
	return "Example Ticker"
}

func (et ExampleTicker) Run() {
	et.ticker = time.NewTicker(time.Second * configuration.ExampleTickerIntervalSec)

	for {
		select {
		case t := <-et.ticker.C:
			log.Printf("tick %v\n", t)
		}
	}
}

func (et *ExampleTicker) GracefulStop(ctx context.Context) error {
	if et.ticker != nil {
		et.ticker.Stop()
	}

	return nil
}

func (et *ExampleTicker) SetInterval(interval time.Duration) {
	et.interval = interval
}

func (et ExampleTicker) GetInterval() time.Duration {
	return et.interval
}
