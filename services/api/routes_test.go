//  +build unit

package api

import (
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"github.com/wbrush/go-template-service/configuration"
	"testing"
)

func TestAPI_initRoutes(t *testing.T) {
	cfg := configuration.InitConfig("", "")

	api := NewAPI(cfg, "../../docs/", nil, nil)
	api.router = mux.NewRouter()
	api.initRoutes(negroni.New())
}
