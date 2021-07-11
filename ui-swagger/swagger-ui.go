package swagger_ui

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"runtime"

	"github.com/gorilla/mux"
	// include static files
)

const (
	SWAGGER_FILE = "swagger.json"
)

func AttachSwaggerUI(router *mux.Router, base_path string, swaggerBase string) (err error) {

	// set swagger-ui routes
	staticPath, err := getWorkingDirectory()
	if err != nil {
		return err
	}

	_, err = os.Stat(swaggerBase + SWAGGER_FILE)
	if err != nil {
		return fmt.Errorf("swagger-ui.AttachSwaggerUI() -> ERROR: swagger.json file does not exists: %s ", err.Error())
	}

	// set swagger.json file route
	router.PathPrefix(base_path + "help/data").Handler(http.StripPrefix(base_path+"help/data", http.FileServer(http.Dir(swaggerBase))))

	router.PathPrefix(base_path + "help/node-modules").Handler(http.StripPrefix(base_path+"help/node-modules", http.FileServer(http.Dir(staticPath+"node-modules"))))
	router.PathPrefix(base_path + "help/swagger-ui-dist").Handler(http.StripPrefix(base_path+"help/swagger-ui-dist", http.FileServer(http.Dir(staticPath+"swagger-ui-dist"))))
	router.PathPrefix(base_path + "help/next-tick").Handler(http.StripPrefix(base_path+"help/next-tick", http.FileServer(http.Dir(staticPath+"next-tick"))))
	router.PathPrefix(base_path + "help/es6-symbol").Handler(http.StripPrefix(base_path+"help/es6-symbol", http.FileServer(http.Dir(staticPath+"es6-symbol"))))
	router.PathPrefix(base_path + "help/es6-iterator").Handler(http.StripPrefix(base_path+"help/es6-iterator", http.FileServer(http.Dir(staticPath+"es6-iterator"))))
	router.PathPrefix(base_path + "help/es5-ext").Handler(http.StripPrefix(base_path+"help/es5-ext", http.FileServer(http.Dir(staticPath+"es5-ext"))))
	router.PathPrefix(base_path + "help/d").Handler(http.StripPrefix(base_path+"help/d", http.FileServer(http.Dir(staticPath+"d"))))
	router.PathPrefix(base_path + "help").Handler(http.StripPrefix(base_path+"help", http.FileServer(http.Dir(staticPath))))

	return nil
}

func getWorkingDirectory() (staticPath string, err error) {
	// get static path from calling lib otherwise
	_, packagePath, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("swagger-ui.AttachSwaggerUI() -> ERROR: Could not get swagger-ui package path: %s", err.Error())
	}

	// set swagger-ui routes
	return path.Dir(packagePath) + "/swagger-ui-static/", nil
}
