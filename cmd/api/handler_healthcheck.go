package main

import (
	"net/http"

	"github.com/ucok-man/gmoapi/cmd/api/config"
)

// @Summary      System Health Check
// @Description  Returns the API health status, environment, and version information. This endpoint does not require authentication and can be used for monitoring.
// @Tags         Health
// @Produce      json
// @Success      200  {object}  object{status=string, system_info=object{environment=string, version=string}}  "API is healthy and operational"
// @Failure      500  {object}  object{error=string}  "Internal server error"
// @Router       / [get]
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
