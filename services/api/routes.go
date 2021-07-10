package api

import (
	"bitbucket.org/optiisolutions/go-common/httphelper"
	"bitbucket.org/optiisolutions/go-template-service/configuration"
	"github.com/urfave/negroni"
)

const (
	TemplatePath = "/template"
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
			Name:        "Create Template",
			Method:      "POST",
			Pattern:     TemplatePath,
			HandlerFunc: api.CreateTemplate,
			Middleware:  []negroni.HandlerFunc{httphelper.MWUserInfoHeader},
		},
		{
			Name:        "Get Template",
			Method:      "GET",
			Pattern:     TemplatePath + "/{id}",
			HandlerFunc: api.GetTemplate,
			Middleware:  []negroni.HandlerFunc{httphelper.MWUserInfoHeader},
		},
		{
			Name:        "List Templates",
			Method:      "GET",
			Pattern:     TemplatePath,
			HandlerFunc: api.ListTemplates,
			Middleware:  []negroni.HandlerFunc{httphelper.MWUserInfoHeader},
		},
		{
			Name:        "Update Template",
			Method:      "PUT",
			Pattern:     TemplatePath + "/{id}",
			HandlerFunc: api.UpdateTemplate,
			Middleware:  []negroni.HandlerFunc{httphelper.MWUserInfoHeader},
		},
		{
			Name:        "Remove Template",
			Method:      "DELETE",
			Pattern:     TemplatePath + "/{id}",
			HandlerFunc: api.DeleteTemplate,
			Middleware:  []negroni.HandlerFunc{httphelper.MWUserInfoHeader},
		},
	})
}
