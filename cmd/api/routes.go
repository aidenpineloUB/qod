// Filename: cmd/api/routes.go

package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
)

func (app *applicationDependencies)routes() http.Handler  {

   // setup a new router
   router := httprouter.New()
   router.NotFound = http.HandlerFunc(app.notFoundResponse)
   router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
   // setup routes
   router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
   router.HandlerFunc(http.MethodPost, "/v1/comments", app.createCommentHandler)

   return app.recoverPanic(router)      
  
}






