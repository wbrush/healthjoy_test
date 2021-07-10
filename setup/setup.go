package setup

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/wbrush/healthjoy_test/configuration"
	"github.com/wbrush/healthjoy_test/services/api"
)

func SetupAndRun(block bool, commit, builtAt, swaggerLoc string) {
	//  perform initialization here
	cfg := configuration.InitConfig(commit, builtAt)

	err := cfg.ConfigureLogger()
	if err != nil {
		logrus.Fatalf("Cannot ConfigureLogger: %s", err.Error())
	}

	logrus.Info("------------------------------")
	logrus.Info("Starting " + configuration.ServiceName)
	logrus.Info("Version:", cfg.Version, "; Build Date:", cfg.BuiltAt)
	logrus.Info("------------------------------")

	time.Sleep(2 * time.Millisecond)

	//  init necessary processing routines (modules)
	//  init REST server
	apiModule := api.NewAPI(cfg, swaggerLoc)

	StartUp(cfg)

	RunModules(apiModule)
	if block {
		WaitForDone()
	}
}

func RunModules(modules ...helpers.Module) {
	if len(modules) > 0 {
		for _, m := range modules {
			if m != nil {
				logrus.Infof("Starting module %s", m.Title())
				go m.Run()
				runningModules = append(runningModules, m)
			}
		}
	}
}

func WaitForDone() {
	if len(runningModules) > 0 {
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt)

		select {
		case <-interrupt:
			fmt.Println() //is used to embellish the output after ^C
			CloseAllModules()
		} //lock execution
	}
}

func CloseAllModules() {
	for _, m := range runningModules {
		logrus.Warnf("Stopping module %s", m.Title())
		ctx, _ := context.WithTimeout(context.Background(), configuration.GracefulStopTimeoutSec*time.Second)
		err := m.GracefulStop(ctx)
		if err != nil {
			logrus.Errorf("can't stop the module %s: %s", m.Title(), err.Error())
		}
	}
}
