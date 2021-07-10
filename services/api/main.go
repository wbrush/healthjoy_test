package api

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
	"github.com/wbrush/go-template-service/configuration"
	"github.com/wbrush/go-template-service/dao"
	"net/http"
	"strings"
)

const DefaultServicePort = "8080"

type (
	Route struct {
		Name        string
		Method      string
		Pattern     string
		HandlerFunc http.HandlerFunc
		Middleware  []negroni.HandlerFunc
	}

	// API serves the end users requests.
	API struct {
		config *configuration.Config

		host        string
		port        string
		swaggerPath string

		router *mux.Router
		server *http.Server
	}
)

// NewAPI initializes a new instance of API with needed fields, but doesn't start listening,
// nor creates the router.
func NewAPI(cfg *configuration.Config, swaggerPath string, dao dao.DataAccessObject, ps messaging.PublisherSubscriber) *API {
	api := &API{
		config: cfg,

		host:        cfg.ServiceParams.Host,
		port:        cfg.ServiceParams.Port,
		swaggerPath: swaggerPath,

		dao: dao,
		ps:  ps,
	}

	if api.port == "" || len(strings.TrimSpace(api.port)) == 0 {
		api.port = DefaultServicePort
	}

	return api
}

// GracefulStop shuts down the server without interrupting any
// active connections.
func (api *API) GracefulStop(ctx context.Context) error {
	cont, cancel := context.WithTimeout(ctx, configuration.GracefulStopTimeoutSec)
	defer cancel()
	return api.server.Shutdown(cont)
}

// Title returns the title.
func (api *API) Title() string {
	return "HTTP REST API"
}

// Run starts the http server and binds the handlers.
func (api *API) Run() {
	api.initialize()
	api.startServe()
}

func (api *API) initialize() {
	api.router = mux.NewRouter()

	wrapper := negroni.New()

	wrapper.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
	}))

	api.initRoutes(wrapper)

	// attach the swagger routes
	err := sw.AttachSwaggerUI(api.router,
		fmt.Sprintf("%s/", configuration.APIBasePath), api.swaggerPath)
	if err != nil {
		logrus.Errorf("error attaching the swagger UI: %s ", err.Error())
	}
}

func (api *API) startServe() {
	logrus.Infof("Starting REST Server on port %s...", api.port)

	connAddress := api.host + ":" + api.port

	api.server = &http.Server{Addr: connAddress, Handler: api.router}
	err := api.server.ListenAndServe()
	if err != nil {
		logrus.Fatalf("cannot start REST Server: %s", err.Error())
	}
}

// HandleActions is used to handle all given routes
func (api *API) HandleActions(wrapper *negroni.Negroni, prefix string, routes []Route) {
	for _, r := range routes {
		w := wrapper.With()
		for _, m := range r.Middleware {
			w.Use(m)
		}

		w.Use(negroni.Wrap(http.HandlerFunc(r.HandlerFunc)))
		api.router.Handle(prefix+r.Pattern, w).Methods(r.Method, "OPTIONS")
	}
}
