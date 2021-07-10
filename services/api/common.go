package api

import (
	"net/http"
	"runtime"

	"bitbucket.org/optiisolutions/go-common/httphelper"
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
	logrus.Tracef("HandleInfo called")

	var memstat runtime.MemStats
	var retValues JsonValues

	runtime.ReadMemStats(&memstat)
	cfg := api.config

	retValues.Version = cfg.Version
	retValues.BuiltOn = cfg.BuiltAt
	retValues.Cpus = runtime.NumCPU()
	retValues.Num_go_routines = runtime.NumGoroutine()
	retValues.Num_cgo_calls = runtime.NumCgoCall()
	retValues.Alloc_heap_total = memstat.HeapAlloc / 1024
	retValues.Alloc_system_total = memstat.Sys / 1024
	retValues.Est_max_heap = memstat.HeapSys / 1024
	retValues.Used_stack = memstat.StackInuse / 1024
	retValues.Stack_max = memstat.StackSys / 1024

	//msg := []string{
	//	fmt.Sprintf("Info Command\n"),
	//	fmt.Sprintf("pii-info-mgr\n"),
	//	fmt.Sprintf("Version: %s; Built on: %s\n", retValues.Version, retValues.BuiltOn),
	//	fmt.Sprintf("Copyright %d Optii\n", time.Now().Year()),
	//	fmt.Sprintf("\n"),
	//	fmt.Sprintf("Service Info\n"),
	//	fmt.Sprintf("    # CPUs:       %d\n", retValues.Cpus),
	//	fmt.Sprintf("    # GoRoutines: %d\n", retValues.Num_go_routines),
	//	fmt.Sprintf("    # cGo Calls:  %d\n", retValues.Num_cgo_calls),
	//	fmt.Sprintf("Memory Info\n"),
	//	fmt.Sprintf("    Total Heap Allocated(KB):         %d\n", retValues.Alloc_heap_total),
	//	fmt.Sprintf("    Total Allocated from System(KB):  %d\n", retValues.Alloc_system_total),
	//	fmt.Sprintf("    Estimated Max Allocated Heap(KB): %d\n", retValues.Est_max_heap),
	//	fmt.Sprintf("    Stack In Use(KB):                 %d\n", retValues.Used_stack),
	//	fmt.Sprintf("    Stack Memory(KB):                 %d\n", retValues.Stack_max),
	//}

	httphelper.Json(w, retValues)
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
	logrus.Tracef("HandlePing called")

	httphelper.Json(w, map[string]interface{}{})
	return
}
