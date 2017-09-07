package handlers

import (
	"log"
	"net/http"

	"github.com/alexellis/faas/gateway/requests"
	"github.com/labstack/echo"
)

func (hl *Handler) Deploy(ctx echo.Context) error {
	opts := new(requests.CreateFunctionRequest)
	if err := ctx.Bind(opts); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	if err := hl.Hyper.Create(opts.Service, opts.Image, buildEnvVars(opts)); err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	log.Println("Deployed function - " + opts.Service)
	return ctx.NoContent(http.StatusAccepted)
}

func buildEnvVars(request *requests.CreateFunctionRequest) []string {
	envVars := []string{}

	if len(request.EnvProcess) > 0 {
		envVar := "fprocess=" + request.EnvProcess
		envVars = append(envVars, envVar)
	}

	for k, v := range request.EnvVars {
		envVar := k + "=" + v
		envVars = append(envVars, envVar)
	}
	return envVars
}
