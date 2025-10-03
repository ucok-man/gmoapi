package main

import (
	"net/http"

	"github.com/ucok-man/gmoapi/cmd/api/config"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	env := envelope{
		"status": "available",
		"system_info": map[string]any{
			"environment": app.config.Env,
			"version":     config.APP_VERSION,
		},
	}

	err := app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
