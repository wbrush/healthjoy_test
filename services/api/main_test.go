//  +build unit

package api

import (
	"net/http"
	"testing"

	"bitbucket.org/optiisolutions/go-template-service/configuration"
	"github.com/gorilla/mux"
)

func TestNewRouter(t *testing.T) {
	cfg := configuration.InitConfig("", "")

	cfg.ServiceParams.Port = ""
	api := NewAPI(cfg, "../../docs/", nil, nil)
	if api.host != "" {
		t.Error("Default DbHost Name Should Not Be Active")
	}
	if api.port == "" {
		t.Error("Default DbPort Number Is Not Active")
	}

	cfg.ServiceParams.Port = "8090"
	api = NewAPI(cfg, "../../docs/", nil, nil)
	if api.host != "" {
		t.Error("Invalid DbHost Name Returned: ", api.host)
	}
	if api.port != "8090" {
		t.Error("Invalid DbPort Number Returned: ", api.port)
	}

	cfg.ServiceParams.Port = "8100"
	cfg.ServiceParams.Host = "KellerRealty"
	api = NewAPI(cfg, "../../docs/", nil, nil)
	if api.host != "KellerRealty" {
		t.Error("Invalid DbHost Name Returned: ", api.host)
	}
	if api.port != "8100" {
		t.Error("Invalid DbPort Number Returned: ", api.port)
	}
}

func TestAPI_Title(t *testing.T) {
	type fields struct {
		config      *configuration.Config
		host        string
		port        string
		swaggerPath string
		router      *mux.Router
		server      *http.Server
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "name return test",
			want: "HTTP REST API",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := &API{
				config:      tt.fields.config,
				host:        tt.fields.host,
				port:        tt.fields.port,
				swaggerPath: tt.fields.swaggerPath,
				router:      tt.fields.router,
				server:      tt.fields.server,
			}
			if got := api.Title(); got != tt.want {
				t.Errorf("API.Title() = %v, want %v", got, tt.want)
			}
		})
	}
}
