// Filename: cmd/api/routes.go

package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
)

// routes specifies our routes (routes.go)
func (app *applicationDependencies) routes() http.Handler {
   router := httprouter.New()
   router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

   return router
}


