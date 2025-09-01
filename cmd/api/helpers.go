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

func (app *applicationDependencies)healthcheckHandler(w http.ResponseWriter,r *http.Request) {
   panic("Apples & Oranges")   // deliberate panic
   data := envelope {
                     "status": "available",
                     "system_info": map[string]string{
                             "environment": app.config.environment,
                             "version": appVersion,
                    },
   }
   err := app.writeJSON(w, http.StatusOK, data, nil)
   if err != nil {
    app.serverErrorResponse(w, r, err)
   }
}

