package api

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/wbrush/healthjoy_test/models"
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

	var inputData models.InputData
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&inputData)
	if err != nil {
		logrus.Warnf("wrong input data provided: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	logrus.WithFields(logrus.Fields{
		"source":      inputData.Source,
		"destination": inputData.Destination,
	}).Debug("requested CopyRepo")

	//  build command - https://stackoverflow.com/questions/44274188/forking-a-github-repo-using-from-the-command-line-with-bash-curl-and-the-githu
	//curl -u 'my_user_name' https://api.github.com/repos/$upstream_repo_username/$upstream_repo_name/forks -d ''

	//  perform command

	//  check status

	w.WriteHeader(http.StatusOK)
	// w.Write([]byte("hello"))

	logrus.Debug("finished CopyFile")
	return
}
