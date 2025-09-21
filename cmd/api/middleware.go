// Filename: cmd/api/middleware.go
package main

import (
  "fmt"
  "net/http"
)

func (a *applicationDependencies)recoverPanic(next http.Handler)http.Handler  {
   return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
   // defer will be called when the stack unwinds
       defer func() {
           // recover() checks for panics
           err := recover();
           if err != nil {
               w.Header().Set("Connection", "close")
               a.serverErrorResponse(w, r, fmt.Errorf("%s", err))
           }
       }()
       next.ServeHTTP(w,r)
   })  
}

func (a *applicationDependencies) enableCORS (next http.Handler) http.Handler {                             
   return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

        w.Header().Add("Vary", "Origin")
       // The request method can vary so don't rely on cache
       w.Header().Add("Vary", "Access-Control-Request-Method")
       origin := r.Header.Get("Origin")

       if origin != "" {
         for i := range a.config.cors.trustedOrigins {
              if origin == a.config.cors.trustedOrigins[i] {
                   w.Header().Set("Access-Control-Allow-Origin", origin)  
                   // check if it is a Preflight CORS request
                     if r.Method == http.MethodOptions && 
                        r.Header.Get("Access-Control-Request-Method") != "" {
        if r.Method == http.MethodOptions && 
           r.Header.Get("Access-Control-Request-Method") != "" {
              w.Header().Set("Access-Control-Allow-Methods",
                             "OPTIONS, PUT, PATCH, DELETE")
              w.Header().Set("Access-Control-Allow-Headers",
                             "Authorization, Content-Type")
        
              // we need to send a 200 OK status. Also since there
              // is no need to continue the middleware chain we
              // we leave  - remember it is not a real 'comments' request but
              // only a preflight CORS request 
              w.WriteHeader(http.StatusOK)
              return
          }
                       return
                    }

                    break
                }
            }
        }

        next.ServeHTTP(w, r)
    })
}





