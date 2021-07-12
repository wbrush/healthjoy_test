// Package classification Template Service for Go-based Development
//
// the purpose of this service is to provide a template for Go-based
// microservice development. This can be used to start development on
// a new service or as an example of Optii accepted best practices for
// developing Go-based microservices
//
// Terms Of Service:
//
// there are no TOS at this moment, use at your own risk we take no responsibility
//
//     Schemes: http
//     BasePath: /api
//     Version: 0.0.1
//
//     Consumes:
//     - application/vnd.api+json
//     - application/json
//
//     Produces:
//     - application/vnd.api+json
//     - application/json
//
//     security:
//     - wbrush_apikey:
//	   - wbrush_oauth2:
//	     read_scope
//
// 	   securityDefinitions:
//       wbrush_apikey:
//         type: apiKey
//         name: KEY
//         in: header
//       wbrush_oauth2:
//         type: oauth2
//         description: example
//         flow: accessCode
//         authorizationUrl: 'https://localhost/oauth2/auth'
//         tokenUrl: 'https://localhost/oauth2/token'
//         scopes:
//           read_scope: description here
//           write_scope: description here
//
// swagger:meta
package main

//go:generate swagger generate spec -m -o ./docs/swagger.json

/*
for more information about generating swagger.json from comments, see:
	https://www.ribice.ba/swagger-golang/
*/
import (
	"github.com/wbrush/healthjoy_test/setup"
)

var (
	commit  string
	builtAt string
)

func main() {
	setup.SetupAndRun(true, commit, builtAt, "./docs/")
}
