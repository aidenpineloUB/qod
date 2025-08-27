// Filename: cmd/api/helpers.go
package main

import (
  "encoding/json"
  "net/http"
)

  // create an envelope type
  type envelope map[string]any
func (a *applicationDependencies)writeJSON(w http.ResponseWriter,
                                           status int, data envelope,
                                           headers http.Header) error  {
    jsResponse, err := json.MarshalIndent(data, "", "\t")
    if err != nil {
        return err
    }
    jsResponse = append(jsResponse, '\n')
    // additional headers to be set
    for key, value := range headers {
        w.Header()[key] = value
    }
    // set content type header
    w.Header().Set("Content-Type", "application/json")
    // explicitly set the response status code
    w.WriteHeader(status) 
    _, err = w.Write(jsResponse) // 👈 This line is the fix
    if err != nil {
        return err
    }

    return nil
}

func (a *applicationDependencies)healthcheckHandler(w http.ResponseWriter,
                                               r *http.Request) {
   data := envelope {
                     "status": "available",
                     "system_info": map[string]string{
                             "environment": a.config.environment,
                             "version": appVersion,
                    },
   }
   err := a.writeJSON(w, http.StatusOK, data, nil)
   if err != nil {
    a.logger.Error(err.Error())
    http.Error(w, "The server encountered a problem and couldnot process your request", http.StatusInternalServerError)
   }
}
