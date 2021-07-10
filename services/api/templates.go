package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	cdatamodels "bitbucket.org/optiisolutions/go-common/datamodels"
	"bitbucket.org/optiisolutions/go-common/db"
	"bitbucket.org/optiisolutions/go-common/errorhandler"
	"bitbucket.org/optiisolutions/go-common/helpers"
	"bitbucket.org/optiisolutions/go-common/httphelper"
	"bitbucket.org/optiisolutions/go-template-service/datamodels"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var requestFunc = httphelper.MakeHTTPRequest

// swagger:operation POST /template entries createTemplate
// ---
// summary: Create new Template
// description: returns new Template
// parameters:
// - name: body
//   in: body
//   description: Template object that needs to be added
//   schema:
//       $ref: "#/definitions/GlobalTemplateStruct"
//   required: true
// consumes:
//   - application/json
// produces:
//   - application/json
// responses:
//   200:
//     description: "OK"
//     schema:
//       type: object
//       $ref: "#/definitions/GlobalTemplateStruct"
func (api *API) CreateTemplate(w http.ResponseWriter, r *http.Request) {
	logrus.Trace("requested CreateTemplate")

	var newTemplate datamodels.Template
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newTemplate)
	if err != nil {
		logrus.Warnf("wrong input data provided: %s", err.Error())
		httphelper.JsonError(w, errorhandler.NewError(errorhandler.ErrBadParam, "body"))
		return
	}

	err = newTemplate.Validate()
	if err != nil {
		logrus.Warnf("Template data is invalid: %s", err.Error())
		httphelper.JsonError(w, errorhandler.NewError(errorhandler.ErrBadParam, "body"))
		return
	}

	//get shard id
	shards, shardsExists := r.Context().Value(httphelper.ShardsCtx).([]int64)
	if !shardsExists || len(shards) < 1 {
		logrus.Errorf("no %s value in request context. Probably middleware was not called?", httphelper.ShardsCtx)
		httphelper.JsonError(w, errorhandler.NewError(errorhandler.ErrBadRequest, httphelper.ShardsCtx))
		return
	}

	isDuplicate, err := api.dao.CreateTemplate(shards[0], &newTemplate)
	if err != nil {
		logrus.Errorf("Template creation error: %s", err.Error())
		httphelper.JsonError(w, errorhandler.NewError(errorhandler.ErrService))
		return
	}
	if isDuplicate {
		logrus.Errorf("Template already exist: %d", newTemplate.Id)
		httphelper.JsonError(w, errorhandler.NewError(errorhandler.ErrAlreadyExists, strconv.FormatInt(newTemplate.Id, 10)))
		return
	}

	httphelper.Json(w, newTemplate)
	logrus.Trace("finished CreateTemplate")
	return
}

// swagger:operation GET /template/{id} entries getTemplate
// ---
// summary: Get Template by given id
// description: returns Template
// parameters:
// - name: id
//   in: path
//   description: Numeric ID of the Template to get
//   type: number
//   required: true
// produces:
//   - application/json
// responses:
//   200:
//     description: "OK"
//     schema:
//       type: object
//       $ref: "#/definitions/GlobalTemplateStruct"
func (api *API) GetTemplate(w http.ResponseWriter, r *http.Request) {
	logrus.Trace("requested GetTemplate")

	ids := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(ids, 10, 64)
	if err != nil {
		logrus.Warnf("wrong Template id %s provided: %s", ids, err.Error())
		httphelper.JsonError(w, errorhandler.NewError(errorhandler.ErrBadParam, ids))
		return
	}

	//get shard id
	shards, shardsExists := r.Context().Value(httphelper.ShardsCtx).([]int64)
	if !shardsExists || len(shards) < 1 {
		logrus.Errorf("no %s value in request context. Probably middleware was not called?", httphelper.ShardsCtx)
		httphelper.JsonError(w, errorhandler.NewError(errorhandler.ErrBadRequest, httphelper.ShardsCtx))
		return
	}

	template, found, err := api.dao.GetTemplateById(shards[0], id)
	if err != nil {
		logrus.Errorf("Template id %d select error: %s", id, err.Error())
		httphelper.JsonError(w, errorhandler.NewError(errorhandler.ErrService))
		return
	}
	if !found || template == nil {
		logrus.Warnf("No such Template by id %d was found", id)
		httphelper.JsonError(w, errorhandler.NewError(errorhandler.ErrNotFound, ids))
		return
	}

	httphelper.Json(w, template)
	logrus.Trace("finished GetTemplate")
	return
}

// swagger:operation GET /template entries listTemplates
// ---
//description: returns Templates list
//produces:
//- application/json
//tags:
//- Template
//summary: Return list of Templates
//operationId: listTemplates
//parameters:
//- in: query
//  name: first
//  description: how many objects must to return
//  schema:
//    type: integer
//- in: query
//  name: after
//  description: cursor of object that must to be returned
//  schema:
//    type: string
//- in: query
//  name: orderBy
//  description: on which field results must to be ordered
//  schema:
//    type: string
//- in: query
//  name: filters
//  description: 'filters on fields, which can be provided in format described here:
//    https://godoc.org/github.com/go-pg/pg/urlvalues#Filter'
//  schema:
//    type: object
//responses:
//  '200':
//    description: OK
//    schema:
//      type: array
//      items:
//        "$ref": "#/definitions/GlobalList"
func (api *API) ListTemplates(w http.ResponseWriter, r *http.Request) {
	logrus.Trace("requested ListTemplates")

	//get shard id
	shards, shardsExists := r.Context().Value(httphelper.ShardsCtx).([]int64)
	if !shardsExists || len(shards) < 1 {
		logrus.Errorf("no %s value in request context. Probably middleware was not called?", httphelper.ShardsCtx)
		httphelper.JsonError(w, errorhandler.NewError(errorhandler.ErrBadRequest, httphelper.ShardsCtx))
		return
	}

	Templates, total, hasNext, err := api.dao.ListTemplates(shards[0], r.URL.Query())
	if err != nil {
		logrus.Errorf("ListTemplateEntries list select error: %s", err.Error())
		httphelper.JsonError(w, errorhandler.NewError(errorhandler.ErrService))
		return
	}

	var list cdatamodels.List
	list.TotalCount = total
	list.Edges = make([]cdatamodels.Edge, 0)
	for i := range Templates {
		list.Edges = append(list.Edges, cdatamodels.Edge{
			Node:   Templates[i],
			Cursor: db.EncodeIdToCursor(Templates[i].Id),
		})
	}

	list.PageInfo.HasNextPage = hasNext
	if len(Templates) > 0 {
		list.PageInfo.EndCursor = db.EncodeIdToCursor(Templates[len(Templates)-1].Id)
	}

	httphelper.Json(w, list)
	logrus.Trace("finished ListTemplates")
	return
}

// swagger:operation PUT /template/{id} entries updateTemplate
// ---
//summary: Update Template by given id
//description: returns Template
//parameters:
//- name: id
//  in: path
//  description: Numeric ID of the Template to update
//  type: number
//  required: true
//- name: body
//  in: body
//  description: Template update object that needs to be updated
//  schema:
//  $ref: "#/definitions/GlobalTemplateUpdateStruct"
//  required: true
//produces:
//  - application/json
//responses:
//  200:
//    description: "OK"
//    schema:
//      type: object
//      $ref: "#/definitions/GlobalTemplateStruct"
func (api *API) UpdateTemplate(w http.ResponseWriter, r *http.Request) {
	logrus.Trace("requested UpdateTemplate")

	ids := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(ids, 10, 64)
	if err != nil {
		logrus.Warnf("wrong pii info id %s provided by pii: %s", ids, err.Error())
		httphelper.JsonError(w, errorhandler.NewError(errorhandler.ErrBadParam, ids))
		return
	}

	//get shard id
	shards, shardsExists := r.Context().Value(httphelper.ShardsCtx).([]int64)
	if !shardsExists || len(shards) < 1 {
		logrus.Errorf("no %s value in request context. Probably middleware was not called?", httphelper.ShardsCtx)
		httphelper.JsonError(w, errorhandler.NewError(errorhandler.ErrBadRequest, httphelper.ShardsCtx))
		return
	}

	template, found, err := api.dao.GetTemplateById(shards[0], id)
	if err != nil {
		logrus.Errorf("Template id %d select error: %s", id, err.Error())
		httphelper.JsonError(w, errorhandler.NewError(errorhandler.ErrService))
		return
	}
	if !found || template == nil {
		logrus.Warnf("No such Template by id %d was found", id)
		httphelper.JsonError(w, errorhandler.NewError(errorhandler.ErrNotFound, ids))
		return
	}

	var updTemplate datamodels.TemplateUpdate
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&updTemplate)
	if err != nil {
		logrus.Warnf("wrong input data provided: %s", err.Error())
		httphelper.JsonError(w, errorhandler.NewError(errorhandler.ErrBadParam, "body"))
		return
	}

	isFoundUpd, err := helpers.NullableFieldsToStruct(updTemplate, template)
	if err != nil {
		logrus.Warnf("cannot fill updateable fields: %s", err.Error())
		httphelper.JsonError(w, errorhandler.NewError(errorhandler.ErrService))
		return
	}
	if !isFoundUpd {
		logrus.Infof("nothing to update for Template %d", template.Id)
		httphelper.Json(w, template)
		return
	}

	err = api.dao.UpdateTemplate(shards[0], template)
	if err != nil {
		logrus.Errorf("Template with id %d update error: %s", id, err.Error())
		httphelper.JsonError(w, errorhandler.NewError(errorhandler.ErrService))
		return
	}

	httphelper.Json(w, template)
	logrus.Trace("finished UpdateTemplate")
	return
}

// swagger:operation DELETE /template/{id} entries deleteTemplate
// ---
//summary: Delete Template by given id
//description: returns id
//parameters:
//- name: id
// in: path
// description: Numeric ID of the Template to delete
// type: number
// required: true
//produces:
//  - application/json
//responses:
//  200:
//    description: "id"
//    schema:
//      type: integer
func (api *API) DeleteTemplate(w http.ResponseWriter, r *http.Request) {
	logrus.Trace("requested DeleteTemplate")

	ids := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(ids, 10, 64)
	if err != nil {
		logrus.Warnf("wrong Template id %s provided: %s", ids, err.Error())
		httphelper.JsonError(w, errorhandler.NewError(errorhandler.ErrBadParam, ids))
		return
	}

	//get shard id
	shards, shardsExists := r.Context().Value(httphelper.ShardsCtx).([]int64)
	if !shardsExists || len(shards) < 1 {
		logrus.Errorf("no %s value in request context. Probably middleware was not called?", httphelper.ShardsCtx)
		httphelper.JsonError(w, errorhandler.NewError(errorhandler.ErrBadRequest, httphelper.ShardsCtx))
		return
	}

	found, err := api.dao.DeleteTemplateById(shards[0], id)
	if err != nil {
		logrus.Errorf("Template id %d delete error: %s", id, err.Error())
		httphelper.JsonError(w, errorhandler.NewError(errorhandler.ErrService))
		return
	}
	if !found {
		logrus.Warnf("No such Template by id %d was found", id)
		httphelper.JsonError(w, errorhandler.NewError(errorhandler.ErrNotFound, ids))
		return
	}

	httphelper.Json(w, map[string]interface{}{"id": id})
	logrus.Trace("finished DeleteTemplate")
	return
}
