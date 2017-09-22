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
		ctx.NoContent(http.StatusBadRequest)
		return err
	}
	envs, config := buildConfig(opts)
	if err := hl.Hyper.Create(opts.Service, opts.Image, envs, config); err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return err
	}
	log.Println("Deployed function - " + opts.Service)
	return ctx.NoContent(http.StatusCreated)
}

func buildConfig(request *requests.CreateFunctionRequest) ([]string, map[string]string) {
	envVars := []string{}
	config := make(map[string]string)

	if len(request.EnvProcess) > 0 {
		envVar := "fprocess=" + request.EnvProcess
		envVars = append(envVars, envVar)
	}

	for k, v := range request.EnvVars {
		if k == "hyper_size" {
			config["hyper_size"] = v
			continue
		}
		envVar := k + "=" + v
		envVars = append(envVars, envVar)
	}
	return envVars, config
}
