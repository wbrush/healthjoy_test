//  +build integration

package integrations

import (
	"github.com/wbrush/go-template-service/configuration"
	"strconv"
	"testing"
)

func TestListTemplates(t *testing.T) {
	cfg := configuration.GetCurrentCfg()
	a, err := strconv.Atoi(cfg.ServiceParams.Port)
	if err != nil {
		t.Errorf("Received error getting port number: %s", err.Error())
	}
	port := strconv.Itoa(a)

	cReturn1 := datamodels.List{}

	qd := make(map[string]interface{})
	qd["orderBy"] = []string{"name_ASC"}

	err = httphelper.MakeHTTPRequest(httphelper.RequestData{
		Headers:   MakeHeaders(ListTemplateShard),
		URL:       "http://localhost:" + port + configuration.APIBasePath + configuration.APIVersion + "/template",
		QueryData: httphelper.MakeQueryData(qd),
	}, &cReturn1, "GET")
	if err != nil {
		t.Errorf("Received error listing templates: %s", err.Error())
	}

	if len(cReturn1.Edges) != 1 {
		t.Errorf("Error listing templates records %d!", len(cReturn1.Edges))
	}
}
