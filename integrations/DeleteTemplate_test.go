//  +build integration

package integrations

import (
	"bitbucket.org/optiisolutions/go-common/httphelper"
	"bitbucket.org/optiisolutions/go-template-service/configuration"
	"strconv"
	"testing"
)

func TestDeleteTemplate(t *testing.T) {
	target := testRecordNum
	cfg := configuration.GetCurrentCfg()
	a, err := strconv.Atoi(cfg.ServiceParams.Port)
	if err != nil {
		t.Errorf("Received error getting port number: %s", err.Error())
	}
	port := strconv.Itoa(a)

	cReturn1 := make(map[string]interface{})

	err = httphelper.MakeHTTPRequest(httphelper.RequestData{
		Headers: MakeHeaders(DeleteTemplateShard),
		URL:     "http://localhost:" + port + configuration.APIBasePath + configuration.APIVersion + "/template/" + strconv.Itoa(int(target)),
	}, &cReturn1, "DELETE")
	if err != nil {
		t.Errorf("Received error deleting template #%d: %s", target, err.Error())
	}

	if cReturn1 == nil || int64(cReturn1["id"].(float64)) != target {
		t.Errorf("Error deleting template record #%d on shard %d!", target, DeleteTemplateShard)
	}
}
