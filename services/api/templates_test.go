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

	cdatamodels "bitbucket.org/optiisolutions/go-common/datamodels"
	"bitbucket.org/optiisolutions/go-common/httphelper"
	"github.com/urfave/negroni"

	"bitbucket.org/optiisolutions/go-common/helpers"
	"bitbucket.org/optiisolutions/go-template-service/configuration"
	"bitbucket.org/optiisolutions/go-template-service/dao/mock_dao"
	"bitbucket.org/optiisolutions/go-template-service/datamodels"
	"github.com/golang/mock/gomock"
)

const headerData = "{\"shard\":[2]}"

var (
	testTemplateData datamodels.Template
)

func fillTemplateTestData() {
	testTemplateData = datamodels.Template{
		Id:          1,
		Name:        "test",
		Status:      datamodels.TemplateStatusNew,
		Description: "a long time ago ina galaxy far far away...",
	}
}

func TestAPI_CreateTemplate(t *testing.T) {
	oldRequestFunc := requestFunc

	//mock the request function
	requestFunc = func(r httphelper.RequestData, dest interface{}, method ...string) error {
		return nil
	}

	fillTemplateTestData()
	cfg := configuration.InitConfig("", "")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name        string
		template    *datamodels.Template
		body        string
		createError error
		wantStatus  int
	}{
		{
			name:        "creating template good case",
			template:    &testTemplateData,
			body:        helpers.HelperJsonMarshallMust(testTemplateData),
			createError: nil,
			wantStatus:  http.StatusOK,
		},
		{
			name:        "bad input data case",
			template:    nil,
			body:        "",
			createError: nil,
			wantStatus:  http.StatusBadRequest,
		},
		{
			name:        "creating template error case",
			template:    &testTemplateData,
			body:        helpers.HelperJsonMarshallMust(testTemplateData),
			createError: fmt.Errorf("test"),
			wantStatus:  http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mdao := mock_dao.NewMockDataAccessObject(ctrl)
			api := NewAPI(cfg, "../../docs/", mdao, nil)

			mdao.
				EXPECT().
				CreateTemplate(gomock.Any(), gomock.Any()).
				Return(false, tt.createError).
				MaxTimes(1)

			helpers.HelperHttpHandlerTest(helpers.HandlerTestData{
				T:           t,
				Method:      "POST",
				Middleware:  []negroni.HandlerFunc{httphelper.MWUserInfoHeader},
				URL:         "/api/v1/template",
				Headers:     map[string]string{httphelper.OPTiiUserInfoHeader: headerData},
				Body:        strings.NewReader(tt.body),
				HandlerFunc: api.CreateTemplate,
				ExpStatus:   tt.wantStatus,
				BodyCheckFunc: func(body *bytes.Buffer) error {
					if tt.wantStatus == http.StatusOK {
						response := body.Bytes()
						values := make(map[string]interface{})
						err := json.Unmarshal(response, &values)
						if err != nil {
							return fmt.Errorf("bad json returned from unmarshal = %s", err.Error())
						}
						if Id, ok := values["templateId"]; !ok || (tt.template != nil && int64(Id.(float64)) != tt.template.Id) {
							return fmt.Errorf("handler returned unexpected body: %#v, when Id expected to be %d, but get %d", values, tt.template.Id, Id)
						}
					}

					return nil
				},
			})
		})
	}

	requestFunc = oldRequestFunc
}

func TestAPI_GetTemplate(t *testing.T) {
	fillTemplateTestData()
	cfg := configuration.InitConfig("", "")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name       string
		id         string
		template   *datamodels.Template
		isFound    bool
		getError   error
		wantStatus int
	}{
		{
			name:       "finding template good case",
			id:         "1",
			template:   &testTemplateData,
			isFound:    true,
			getError:   nil,
			wantStatus: http.StatusOK,
		},
		{
			name:       "bad id format case",
			id:         "bad",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "finding template error",
			id:         "2",
			template:   nil,
			isFound:    false,
			getError:   fmt.Errorf("test error"),
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "template not found",
			id:         "2",
			template:   nil,
			isFound:    false,
			getError:   nil,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mdao := mock_dao.NewMockDataAccessObject(ctrl)
			api := NewAPI(cfg, "../../docs/", mdao, nil)

			mdao.
				EXPECT().
				GetTemplateById(gomock.Any(), gomock.Any()).
				Return(tt.template, tt.isFound, tt.getError).
				MaxTimes(1)

			helpers.HelperHttpHandlerTest(helpers.HandlerTestData{
				T:           t,
				Method:      "GET",
				Middleware:  []negroni.HandlerFunc{httphelper.MWUserInfoHeader},
				URLMask:     "/api/v1/template/{id}",
				URL:         "/api/v1/template/" + tt.id,
				Headers:     map[string]string{httphelper.OPTiiUserInfoHeader: headerData},
				Body:        nil,
				HandlerFunc: api.GetTemplate,
				ExpStatus:   tt.wantStatus,
				BodyCheckFunc: func(body *bytes.Buffer) error {
					if tt.wantStatus == http.StatusOK {
						response := body.Bytes()
						values := make(map[string]interface{})
						err := json.Unmarshal(response, &values)
						if err != nil {
							return fmt.Errorf("bad json returned from unmarshal = %s", err.Error())
						}
						if id, ok := values["templateId"]; !ok || !strings.EqualFold(fmt.Sprintf("%d", int64(id.(float64))), tt.id) {
							return fmt.Errorf("handler returned unexpected body: %#v", values)
						}
					}

					return nil
				},
			})
		})
	}
}

func TestAPI_ListTemplates(t *testing.T) {
	fillTemplateTestData()
	cfg := configuration.InitConfig("", "")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name       string
		filters    url.Values
		entries    []datamodels.Template
		total      int
		hasMore    bool
		getError   error
		wantStatus int
	}{
		{
			name: "finding template good case",
			filters: url.Values{
				"name": []string{string(testTemplateData.Name)},
			},
			entries:    []datamodels.Template{testTemplateData},
			total:      1,
			hasMore:    false,
			getError:   nil,
			wantStatus: http.StatusOK,
		},
		{
			name:       "unknown filters case",
			filters:    helpers.HelperSchemaEncodeMust(struct{ f string }{f: "test"}),
			entries:    []datamodels.Template{testTemplateData},
			total:      0,
			hasMore:    false,
			getError:   nil,
			wantStatus: http.StatusOK,
		},
		{
			name: "finding template fails case",
			filters: url.Values{
				"name": []string{string(testTemplateData.Name)},
			},
			entries:    []datamodels.Template{},
			total:      0,
			hasMore:    false,
			getError:   fmt.Errorf("test"),
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mdao := mock_dao.NewMockDataAccessObject(ctrl)
			api := NewAPI(cfg, "../../docs/", mdao, nil)

			mdao.
				EXPECT().
				ListTemplates(gomock.Any(), gomock.Any()).
				Return(tt.entries, tt.total, tt.hasMore, tt.getError).
				MaxTimes(1)

			helpers.HelperHttpHandlerTest(helpers.HandlerTestData{
				T:           t,
				Method:      "GET",
				Middleware:  []negroni.HandlerFunc{httphelper.MWUserInfoHeader},
				URL:         "/api/v1/template",
				Headers:     map[string]string{httphelper.OPTiiUserInfoHeader: headerData},
				Query:       tt.filters,
				HandlerFunc: api.ListTemplates,
				ExpStatus:   tt.wantStatus,
				BodyCheckFunc: func(body *bytes.Buffer) error {
					if tt.wantStatus == http.StatusOK {
						response := body.Bytes()
						var list cdatamodels.List
						err := json.Unmarshal(response, &list)
						if err != nil {
							return fmt.Errorf("bad json returned from unmarshal = %s", err.Error())
						}
						if len(list.Edges) != len(tt.entries) {
							return fmt.Errorf("handler returned unexpected body list want len %d got len %d: ", len(list.Edges), len(tt.entries))
						}
						if inv, ok := list.Edges[0].Node.(map[string]interface{}); !ok || (len(tt.entries) > 0 && int64(inv["templateId"].(float64)) != tt.entries[0].Id) {
							return fmt.Errorf("handler returned unexpected body want %d got %d: %#v", int64(inv["templateId"].(float64)), tt.entries[0].Id, list)
						}
					}

					return nil
				},
			})
		})
	}
}

func TestAPI_UpdateTemplate(t *testing.T) {
	fillTemplateTestData()
	cfg := configuration.InitConfig("", "")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	newStatus := datamodels.TemplateStatusActive
	tests := []struct {
		name       string
		id         string
		template   *datamodels.Template
		isFound    bool
		getError   error
		body       string
		isFoundUpd bool
		updError   error
		wantStatus int
	}{
		{
			name:     "updating template good case",
			id:       "1",
			template: &testTemplateData,
			isFound:  true,
			getError: nil,
			body: helpers.HelperJsonMarshallMust(datamodels.TemplateUpdate{
				Status: &newStatus,
			}),
			isFoundUpd: true,
			updError:   nil,
			wantStatus: http.StatusOK,
		},
		{
			name:       "bad id format case",
			id:         "bad",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "finding template error",
			id:         "2",
			template:   nil,
			isFound:    false,
			getError:   fmt.Errorf("test error"),
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "template not found",
			id:         "2",
			template:   nil,
			isFound:    false,
			getError:   nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "bad update data case",
			id:         "1",
			template:   &testTemplateData,
			isFound:    true,
			getError:   nil,
			body:       "",
			isFoundUpd: false,
			updError:   nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "nothing to update case",
			id:         "1",
			template:   &testTemplateData,
			isFound:    true,
			getError:   nil,
			body:       helpers.HelperJsonMarshallMust(struct{ f *datamodels.TemplateStatus }{f: &newStatus}),
			isFoundUpd: false,
			updError:   nil,
			wantStatus: http.StatusOK,
		},
		{
			name:     "update error case",
			id:       "1",
			template: &testTemplateData,
			isFound:  true,
			getError: nil,
			body: helpers.HelperJsonMarshallMust(datamodels.TemplateUpdate{
				Status: &newStatus,
			}),
			isFoundUpd: true,
			updError:   fmt.Errorf("test"),
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mdao := mock_dao.NewMockDataAccessObject(ctrl)
			api := NewAPI(cfg, "../../docs/", mdao, nil)

			mdao.
				EXPECT().
				GetTemplateById(gomock.Any(), gomock.Any()).
				Return(tt.template, tt.isFound, tt.getError).
				MaxTimes(1)

			mdao.
				EXPECT().
				UpdateTemplate(gomock.Any(), gomock.Any()).
				Return(tt.updError).
				MaxTimes(1)

			helpers.HelperHttpHandlerTest(helpers.HandlerTestData{
				T:           t,
				Method:      "PUT",
				Middleware:  []negroni.HandlerFunc{httphelper.MWUserInfoHeader},
				URLMask:     "/api/v1/template/{id}",
				URL:         "/api/v1/template/" + tt.id,
				Headers:     map[string]string{httphelper.OPTiiUserInfoHeader: headerData},
				Body:        strings.NewReader(tt.body),
				HandlerFunc: api.UpdateTemplate,
				ExpStatus:   tt.wantStatus,
				BodyCheckFunc: func(body *bytes.Buffer) error {
					if tt.wantStatus == http.StatusOK {
						response := body.Bytes()
						values := make(map[string]interface{})
						err := json.Unmarshal(response, &values)
						if err != nil {
							return fmt.Errorf("bad json returned from unmarshal = %s", err.Error())
						}
						if id, ok := values["templateId"]; !ok || !strings.EqualFold(fmt.Sprintf("%d", int64(id.(float64))), tt.id) {
							return fmt.Errorf("handler returned unexpected body: %#v", values)
						}
					}

					return nil
				},
			})
		})
	}
}

func TestAPI_DeleteTemplate(t *testing.T) {
	fillTemplateTestData()
	cfg := configuration.InitConfig("", "")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name       string
		id         string
		isFound    bool
		delError   error
		wantStatus int
	}{
		{
			name:       "deleting template good case",
			id:         "1",
			isFound:    true,
			delError:   nil,
			wantStatus: http.StatusOK,
		},
		{
			name:       "bad id format case",
			id:         "bad",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "finding template error",
			id:         "2",
			isFound:    false,
			delError:   fmt.Errorf("test error"),
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "template not found",
			id:         "2",
			isFound:    false,
			delError:   nil,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mdao := mock_dao.NewMockDataAccessObject(ctrl)
			api := NewAPI(cfg, "../../docs/", mdao, nil)

			mdao.
				EXPECT().
				DeleteTemplateById(gomock.Any(), gomock.Any()).
				Return(tt.isFound, tt.delError).
				MaxTimes(1)

			helpers.HelperHttpHandlerTest(helpers.HandlerTestData{
				T:           t,
				Method:      "DELETE",
				Middleware:  []negroni.HandlerFunc{httphelper.MWUserInfoHeader},
				URLMask:     "/api/v1/template/{id}",
				URL:         "/api/v1/template/" + tt.id,
				Headers:     map[string]string{httphelper.OPTiiUserInfoHeader: headerData},
				Body:        nil,
				HandlerFunc: api.DeleteTemplate,
				ExpStatus:   tt.wantStatus,
				BodyCheckFunc: func(body *bytes.Buffer) error {
					if tt.wantStatus == http.StatusOK {
						response := body.Bytes()
						values := make(map[string]interface{})
						err := json.Unmarshal(response, &values)
						if err != nil {
							return fmt.Errorf("bad json returned from unmarshal = %s", err.Error())
						}
						if id, ok := values["id"]; !ok || !strings.EqualFold(fmt.Sprintf("%d", int64(id.(float64))), tt.id) {
							return fmt.Errorf("handler returned unexpected body: %#v", values)
						}
					}

					return nil
				},
			})
		})
	}
}
