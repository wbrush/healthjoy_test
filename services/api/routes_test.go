//  +build unit

package api

import (
	"bitbucket.org/optiisolutions/go-template-service/configuration"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"testing"
)

func TestAPI_initRoutes(t *testing.T) {
	cfg := configuration.InitConfig("", "")

	api := NewAPI(cfg, "../../docs/", nil, nil)
	api.router = mux.NewRouter()
	api.initRoutes(negroni.New())
}
