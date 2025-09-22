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



func (a *applicationDependencies) enableCORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Add("Vary", "Origin")
        w.Header().Add("Vary", "Access-Control-Request-Method")
        
        origin := r.Header.Get("Origin")
        if origin != "" {
            for i := range a.config.cors.trustedOrigins {
                if origin == a.config.cors.trustedOrigins[i] {
                    w.Header().Set("Access-Control-Allow-Origin", origin)
                    
                    // Handle preflight CORS request
                    if r.Method == http.MethodOptions &&
                        r.Header.Get("Access-Control-Request-Method") != "" {
                        w.Header().Set("Access-Control-Allow-Methods",
                            "OPTIONS, PUT, PATCH, DELETE, GET, POST") // Added GET and POST
                        w.Header().Set("Access-Control-Allow-Headers",
                            "Authorization, Content-Type")
                        w.WriteHeader(http.StatusOK)
                        return
                    }
                    break // Only break here, don't return
                }
            }
        }
        next.ServeHTTP(w, r)
    })
}