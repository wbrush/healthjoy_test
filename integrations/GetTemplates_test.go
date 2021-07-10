//  +build integration

package integrations

import (
	"github.com/wbrush/go-template-service/configuration"
	"github.com/wbrush/go-template-service/datamodels"
	"strconv"
	"testing"
)

func TestGetTemplate(t *testing.T) {
	cfg := configuration.GetCurrentCfg()
	a, err := strconv.Atoi(cfg.ServiceParams.Port)
	if err != nil {
		t.Errorf("Received error getting port number: %s", err.Error())
	}
	port := strconv.Itoa(a)

	cRecord1 := datamodels.Template{
		Id:          testRecordNum,
		Name:        "Front Office Agent",
		Status:      datamodels.TemplateStatusNew,
		Description: "",
	}
	cReturn1 := datamodels.Template{}

	err = httphelper.MakeHTTPRequest(httphelper.RequestData{
		Headers: MakeHeaders(GetTemplateShard),
		URL:     "http://localhost:" + port + configuration.APIBasePath + configuration.APIVersion + "/template/" + strconv.Itoa(int(testRecordNum)),
	}, &cReturn1, "GET")
	if err != nil {
		t.Errorf("Received error getting template #%d: %s", cRecord1.Id, err.Error())
	}

	if cReturn1.Status != cRecord1.Status {
		t.Errorf("Error getting template record #%d", testRecordNum)
	}
}
