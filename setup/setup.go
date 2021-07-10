package setup

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/wbrush/healthjoy_test/configuration"
	"github.com/wbrush/healthjoy_test/services/api"
)

func SetupAndRun(block bool, commit, builtAt, swaggerLoc string) {
	//  perform initialization here
	cfg := configuration.InitConfig(commit, builtAt)

	logrus.Info("------------------------------")
	logrus.Info("Starting " + configuration.ServiceName)
	logrus.Info("Version:", cfg.Commit, "; Build Date:", cfg.BuiltAt)
	logrus.Info("------------------------------")

	time.Sleep(2 * time.Millisecond)

	//  init REST server
	StartUp(cfg)
	apiModule := api.NewAPI(cfg, swaggerLoc)
	apiModule.Initialize()
	apiModule.StartServe()
}
