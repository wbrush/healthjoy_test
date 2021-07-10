//  +build unit

package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/wbrush/go-template-service/configuration"
	"net/http"
	"strings"
	"testing"
)

func TestHandleInfo(t *testing.T) {
	cfg := configuration.InitConfig("", "")
	api := NewAPI(cfg, "../../docs/", nil, nil)

	helpers.HelperHttpHandlerTest(helpers.HandlerTestData{
		T:           t,
		Method:      "GET",
		URL:         "/api/v1/info",
		Body:        nil,
		HandlerFunc: api.HandleInfo,
		ExpStatus:   http.StatusOK,
		BodyCheckFunc: func(body *bytes.Buffer) error {
			response := body.Bytes()
			values := JsonValues{}
			err := json.Unmarshal(response, &values)
			if err != nil {
				return fmt.Errorf("bad json returned from unmarshal = %s", err.Error())
			}
			if values.Num_go_routines < 2 || values.Num_go_routines > 5 {
				return fmt.Errorf("handler returned unexpected body: got text, want json: %#v", values)
			}

			return nil
		},
	})
}

func TestHandlePing(t *testing.T) {
	cfg := configuration.InitConfig("", "")
	api := NewAPI(cfg, "../../docs/", nil, nil)

	helpers.HelperHttpHandlerTest(helpers.HandlerTestData{
		T:           t,
		Method:      "GET",
		URL:         "/api/v1/ping",
		Body:        nil,
		HandlerFunc: api.HandlePing,
		ExpStatus:   http.StatusOK,
		BodyCheckFunc: func(body *bytes.Buffer) error {
			expected := "{}"
			if !strings.Contains(body.String(), expected) {
				return fmt.Errorf("got %#v want %#v", body.String(), expected)
			}
			return nil
		},
	})
}
