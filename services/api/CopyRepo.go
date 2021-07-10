package api

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

// swagger:operation POST /api/v1/copy_repo repo CopyFile
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
func (api *API) CopyRepo(w http.ResponseWriter, r *http.Request) {
	logrus.Debug("requested CopyRepo")

	// var newTemplate datamodels.Template
	// decoder := json.NewDecoder(r.Body)
	// err := decoder.Decode(&newTemplate)
	// if err != nil {
	// 	logrus.Warnf("wrong input data provided: %s", err.Error())
	// 	httphelper.JsonError(w, errorhandler.NewError(errorhandler.ErrBadParam, "body"))
	// 	return
	// }

	// err = newTemplate.Validate()
	// if err != nil {
	// 	logrus.Warnf("Template data is invalid: %s", err.Error())
	// 	httphelper.JsonError(w, errorhandler.NewError(errorhandler.ErrBadParam, "body"))
	// 	return
	// }

	// //get shard id
	// shards, shardsExists := r.Context().Value(httphelper.ShardsCtx).([]int64)
	// if !shardsExists || len(shards) < 1 {
	// 	logrus.Errorf("no %s value in request context. Probably middleware was not called?", httphelper.ShardsCtx)
	// 	httphelper.JsonError(w, errorhandler.NewError(errorhandler.ErrBadRequest, httphelper.ShardsCtx))
	// 	return
	// }

	// isDuplicate, err := api.dao.CreateTemplate(shards[0], &newTemplate)
	// if err != nil {
	// 	logrus.Errorf("Template creation error: %s", err.Error())
	// 	httphelper.JsonError(w, errorhandler.NewError(errorhandler.ErrService))
	// 	return
	// }
	// if isDuplicate {
	// 	logrus.Errorf("Template already exist: %d", newTemplate.Id)
	// 	httphelper.JsonError(w, errorhandler.NewError(errorhandler.ErrAlreadyExists, strconv.FormatInt(newTemplate.Id, 10)))
	// 	return
	// }

	w.WriteHeader(http.StatusOK)
	logrus.Debug("finished CopyFile")
	return
}
