// Filename: cmd/api/routes.go

package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
)

// routes specifies our routes (routes.go)
func (app *applicationDependencies) routes() http.Handler {
   // setup a new router
   router := httprouter.New()
   // handle 404
   router.NotFound = http.HandlerFunc(app.notFoundResponse)
  // handle 405
   router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
  
   router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

   return router
}



