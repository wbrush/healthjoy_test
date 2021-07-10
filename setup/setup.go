package setup

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/wbrush/go-template-service/configuration"
	"github.com/wbrush/go-template-service/dao/postgres"
	"github.com/wbrush/go-template-service/services/api"
	"github.com/wbrush/go-template-service/services/subscriber"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"time"
)

var (
	runningModules []helpers.Module
)

func init() {
	if runningModules == nil {
		runningModules = make([]helpers.Module, 0)
	}
}

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

	ps, err := messaging.NewGPubSub(cfg.GCP.ProjectID)
	if err != nil {
		msg := "Failed to create pubsub client: %v"
		logrus.Errorf(msg, err)
	}
	if ps != nil {
		nc, nce := strconv.Atoi(cfg.NumConns)
		if nce == nil {
			if nc > 2 && nc < 6 {
				nc -= 2
			} else if nc >= 6 && nc < 13 {
				nc -= 3
			} else if nc >= 13 {
				nc = 10
			}
			logrus.Debugf("main(): Setting number go threads to %d", nc)
			ps.SetNumGoThreads(nc)
		}
	}

	//connect to db
	dao, err := initDAO(cfg, ps)
	if err != nil {
		logrus.Fatalf("Cannot init PostgressDAO: %s", err.Error())
	}
	if block {
		defer dao.Close()
	}

	//for integrations
	postgres.SetPgDao(dao)

	// init subscriptions
	subModule := subscriber.NewSub(cfg, dao, ps)

	//  init necessary processing routines (modules)
	//  init REST server
	apiModule := api.NewAPI(cfg, swaggerLoc, dao, ps)

	StartUp(cfg, dao, ps)

	RunModules(apiModule, subModule)
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

func initDAO(cfg *configuration.Config, ps *messaging.GPubSub) (*postgres.PgDAO, error) {
	d, err := postgres.NewPgDao(cfg)
	if err != nil && err.Error() == db.ErrNoShardsYet { //only if no shards
		if d != nil && ps != nil {
			d.LockMux.Lock() // lock mutex before shards init
			var mf messaging.FilterType
			err = mf.NewV2("#", reflect.TypeOf(datamodels.Property{}).Name(), string(messaging.CommandMessageActionSync),
				"#", "#", configuration.ServiceName)
			if err != nil {
				return nil, fmt.Errorf("failed to create base message: %s", err.Error())
			}

			pubSubId, err := ps.Publish(cfg.GCP.ServicePubTopic, &messaging.CommandMessage{
				BaseMessageData: messaging.BaseMessageData{
					Filter: mf,
					Type:   messaging.MessageTypeChangeMessage,
				},
				Action: messaging.CommandMessageActionSync,
			})

			if err != nil {
				return nil, fmt.Errorf("failed to publish base message id %s: %s", pubSubId, err.Error())
			}
			logrus.Infof("message %s published with id %s", mf, pubSubId)
		}

		return d, nil //no shards yet, but base dao exists
	}

	if err != nil { //for all other errors
		return nil, err
	}

	return d, nil //everything fine
}
