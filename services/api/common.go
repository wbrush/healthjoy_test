package api

import (
	"encoding/json"
	"net/http"
	"runtime"

	"github.com/sirupsen/logrus"
)

// swagger:model BasicInfoStruct
type JsonValues struct {
	Version string `json:"version"`
	BuiltOn string `json:"build_date"`

	Cpus            int   `json:"cpus"`
	Num_go_routines int   `json:"num_go_routines"`
	Num_cgo_calls   int64 `json:"num_cgo_calls"`

	Alloc_heap_total   uint64 `json:"alloc_heap_total"`
	Alloc_system_total uint64 `json:"alloc_system_total"`
	Est_max_heap       uint64 `json:"est_max_heap"`
	Used_stack         uint64 `json:"used_stack"`
	Stack_max          uint64 `json:"max_stack"`
}

// swagger:operation GET /info basicCommands basicInfo
// ---
// summary: Get basic service information to aid in debugging.
// description: returns memory stats, cpu stats, and other information useful in debugging problems in deployed service.
// parameters:
// - name: Accept
//   in: header
//   description: standard {Accept} header values
//   type: string
//   required: true
// produces:
//   - application/json
//   - application/text
// responses:
//   200:
//     description: "OK"
//     schema:
//       type: object
//       $ref: "#/definitions/BasicInfoStruct"
func (api *API) HandleInfo(w http.ResponseWriter, r *http.Request) {
	logrus.Debugf("HandleInfo called")

	var memstat runtime.MemStats
	var retValues JsonValues

	runtime.ReadMemStats(&memstat)
	cfg := api.config

	retValues.Version = cfg.Commit
	retValues.BuiltOn = cfg.BuiltAt
	retValues.Cpus = runtime.NumCPU()
	retValues.Num_go_routines = runtime.NumGoroutine()
	retValues.Num_cgo_calls = runtime.NumCgoCall()
	retValues.Alloc_heap_total = memstat.HeapAlloc / 1024
	retValues.Alloc_system_total = memstat.Sys / 1024
	retValues.Est_max_heap = memstat.HeapSys / 1024
	retValues.Used_stack = memstat.StackInuse / 1024
	retValues.Stack_max = memstat.StackSys / 1024

	w.WriteHeader(http.StatusOK)
	jsonResp, err := json.Marshal(retValues)
	if err != nil {
		logrus.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
	return
}

// swagger:operation GET /ping basicCommands basicPing
// ---
// summary: Get basic service information to aid in debugging.
// description: returns clean JSON object to check that everything is fine
// parameters:
// - name: Accept
//   in: header
//   description: standard "Accept" header values
//   type: string
//   required: true
// produces:
//   - application/json
// responses:
//   '200':
//     description: "OK; returns empty page or json structure"
func (api *API) HandlePing(w http.ResponseWriter, r *http.Request) {
	logrus.Debugf("HandlePing called")

	w.WriteHeader(http.StatusOK)
	return
}
