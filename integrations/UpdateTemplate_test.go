//  +build integration

package integrations

import (
	"github.com/wbrush/go-template-service/configuration"
	"strconv"
	"testing"

	"github.com/wbrush/go-template-service/models"
)

func TestUpdateTemplate_break(t *testing.T) {
	target := testRecordNum
	cfg := configuration.GetCurrentCfg()
	a, err := strconv.Atoi(cfg.ServiceParams.Port)
	if err != nil {
		t.Errorf("Received error getting port number: %s", err.Error())
	}
	port := strconv.Itoa(a)

	//  get baseline record
	cRecord1 := datamodels.Template{}
	err = httphelper.MakeHTTPRequest(httphelper.RequestData{
		Headers: MakeHeaders(UpdateTemplateShard),
		URL:     "http://localhost:" + port + configuration.APIBasePath + configuration.APIVersion + "/template/" + strconv.Itoa(int(target)),
	}, &cRecord1, "GET")
	if err != nil {
		t.Errorf("Received error getting template #%d: %s", target, err.Error())
	}

	//  make a change (or two)
	cRecord1.Status = datamodels.TemplateStatusActive

	//  write updates
	cReturn1 := datamodels.Template{}
	err = httphelper.MakeHTTPRequest(httphelper.RequestData{
		Headers: MakeHeaders(UpdateTemplateShard),
		URL:     "http://localhost:" + port + configuration.APIBasePath + configuration.APIVersion + "/template/" + strconv.Itoa(int(target)),
		Json:    &cRecord1,
	}, &cReturn1, "PUT")
	if err != nil {
		t.Errorf("Received error updating template #%d: %s", target, err.Error())
	}

	if cReturn1.Status != datamodels.TemplateStatusActive {
		t.Errorf("Error applying updates to template record!")
	}

	//  back out the changes
	cRecord1.Status = datamodels.TemplateStatusNew

	//  write updates
	cReturn1 = datamodels.Template{}
	err = httphelper.MakeHTTPRequest(httphelper.RequestData{
		Headers: MakeHeaders(UpdateTemplateShard),
		URL:     "http://localhost:" + port + configuration.APIBasePath + configuration.APIVersion + "/template/" + strconv.Itoa(int(target)),
		Json:    &cRecord1,
	}, &cReturn1, "PUT")
	if err != nil {
		t.Errorf("Received error updating template #%d: %s", target, err.Error())
	}

	if cReturn1.Status != datamodels.TemplateStatusNew {
		t.Errorf("Error applying updates to template record!")
	}
}
