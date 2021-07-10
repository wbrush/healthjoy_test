//  +build integration

package integrations

import (
	"github.com/wbrush/go-template-service/configuration"
	"github.com/wbrush/go-template-service/datamodels"
	"strconv"
	"testing"
)

func TestCreateTemplate(t *testing.T) {
	cfg := configuration.GetCurrentCfg()
	a, err := strconv.Atoi(cfg.ServiceParams.Port)
	if err != nil {
		t.Errorf("Received error getting port number: %s", err.Error())
	}
	port := strconv.Itoa(a)

	cRecord1 := datamodels.Template{
		Name:         "Template 1",
		Status:       datamodels.TemplateStatusNew,
		Description:  "Integration test Template 1",
		TemplateSelf: "http://localhost:8000/template/1",
	}
	cReturn1 := datamodels.Template{}

	err = httphelper.MakeHTTPRequest(httphelper.RequestData{
		Headers: MakeHeaders(CreateTemplateShard),
		URL:     "http://localhost:" + port + configuration.APIBasePath + configuration.APIVersion + "/template",
		Json:    &cRecord1,
	}, &cReturn1, "POST")
	if err != nil {
		t.Errorf("Received error creating template (%s): %s", cRecord1.Name, err.Error())
	}

	if cReturn1.Name != cRecord1.Name ||
		cReturn1.Description != cRecord1.Description {
		t.Errorf("Error creating Template record!")
	}
}
