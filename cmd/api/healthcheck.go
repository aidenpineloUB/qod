// Filename: cmd/api/healthcheck.go

package main

import(
	"net/http"
)


func (appInstance *applicationDependencies) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	// Remove the deliberate panic for proper functioning
	// panic("Apples & Oranges")   // deliberate panic
	data := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": appInstance.config.environment,
			"version":     appVersion,
		},
	}
	err := appInstance.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		appInstance.serverErrorResponse(w, r, err)
	}
}
