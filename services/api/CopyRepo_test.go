//  +build unit

package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/urfave/negroni"

]	"github.com/wbrush/go-template-service/configuration"
	"github.com/wbrush/go-template-service/dao/mock_dao"
	"github.com/wbrush/go-template-service/models"
	"github.com/golang/mock/gomock"
)

func TestAPI_CreateTemplate(t *testing.T) {
	// oldRequestFunc := requestFunc

	// //mock the request function
	// requestFunc = func(r httphelper.RequestData, dest interface{}, method ...string) error {
	// 	return nil
	// }

	// fillTemplateTestData()
	// cfg := configuration.InitConfig("", "")

	// ctrl := gomock.NewController(t)
	// defer ctrl.Finish()

	// tests := []struct {
	// 	name        string
	// 	template    *datamodels.Template
	// 	body        string
	// 	createError error
	// 	wantStatus  int
	// }{
	// 	{
	// 		name:        "creating template good case",
	// 		template:    &testTemplateData,
	// 		body:        helpers.HelperJsonMarshallMust(testTemplateData),
	// 		createError: nil,
	// 		wantStatus:  http.StatusOK,
	// 	},
	// 	{
	// 		name:        "bad input data case",
	// 		template:    nil,
	// 		body:        "",
	// 		createError: nil,
	// 		wantStatus:  http.StatusBadRequest,
	// 	},
	// 	{
	// 		name:        "creating template error case",
	// 		template:    &testTemplateData,
	// 		body:        helpers.HelperJsonMarshallMust(testTemplateData),
	// 		createError: fmt.Errorf("test"),
	// 		wantStatus:  http.StatusInternalServerError,
	// 	},
	// }

	// for _, tt := range tests {
	// 	t.Run(tt.name, func(t *testing.T) {
	// 		mdao := mock_dao.NewMockDataAccessObject(ctrl)
	// 		api := NewAPI(cfg, "../../docs/", mdao, nil)

	// 		mdao.
	// 			EXPECT().
	// 			CreateTemplate(gomock.Any(), gomock.Any()).
	// 			Return(false, tt.createError).
	// 			MaxTimes(1)

	// 		helpers.HelperHttpHandlerTest(helpers.HandlerTestData{
	// 			T:           t,
	// 			Method:      "POST",
	// 			Middleware:  []negroni.HandlerFunc{httphelper.MWUserInfoHeader},
	// 			URL:         "/api/v1/template",
	// 			Headers:     map[string]string{httphelper.OPTiiUserInfoHeader: headerData},
	// 			Body:        strings.NewReader(tt.body),
	// 			HandlerFunc: api.CreateTemplate,
	// 			ExpStatus:   tt.wantStatus,
	// 			BodyCheckFunc: func(body *bytes.Buffer) error {
	// 				if tt.wantStatus == http.StatusOK {
	// 					response := body.Bytes()
	// 					values := make(map[string]interface{})
	// 					err := json.Unmarshal(response, &values)
	// 					if err != nil {
	// 						return fmt.Errorf("bad json returned from unmarshal = %s", err.Error())
	// 					}
	// 					if Id, ok := values["templateId"]; !ok || (tt.template != nil && int64(Id.(float64)) != tt.template.Id) {
	// 						return fmt.Errorf("handler returned unexpected body: %#v, when Id expected to be %d, but get %d", values, tt.template.Id, Id)
	// 					}
	// 				}

	// 				return nil
	// 			},
	// 		})
	// 	})
	// }

	// requestFunc = oldRequestFunc
}
