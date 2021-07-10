package api

import (
	"github.com/urfave/negroni"
	"github.com/wbrush/go-template-service/configuration"
)

const (
	TemplatePath = "/api/v1/copy_repo"
)

func (api *API) initRoutes(wrapper *negroni.Negroni) {
	api.HandleActions(wrapper, configuration.APIBasePath, []Route{
		{
			Name:        "Info",
			Method:      "GET",
			Pattern:     "/info",
			HandlerFunc: api.HandleInfo,
			Middleware:  nil,
		},
		{
			Name:        "Ping",
			Method:      "GET",
			Pattern:     "/ping",
			HandlerFunc: api.HandlePing,
			Middleware:  nil,
		},
	})
	api.HandleActions(wrapper, configuration.APIBasePath+configuration.APIVersion, []Route{
		//  application specific
		{
			Name:        "Copy Github Repo",
			Method:      "POST",
			Pattern:     TemplatePath,
			HandlerFunc: api.CopyRepo,
			Middleware:  []negroni.HandlerFunc{httphelper.MWUserInfoHeader},
		},
	})
}
