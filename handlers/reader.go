package handlers

import (
	"log"
	"net/http"

	"github.com/alexellis/faas/gateway/requests"
	"github.com/labstack/echo"
)

func (hl *Handler) List(ctx echo.Context) error {
	functions, err := hl.Hyper.List()
	if err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, functions)
}

func (hl *Handler) Inspect(ctx echo.Context) error {
	name := ctx.Param("name")
	function, err := hl.Hyper.Inspect(name)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	log.Println("Inspected function - " + name)
	return ctx.JSON(http.StatusOK, function)
}

func (hl *Handler) Scale(ctx echo.Context) error {
	name := ctx.Param("name")

	opts := new(requests.ScaleServiceRequest)
	if err := ctx.Bind(opts); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	err := hl.Hyper.Scale(name, opts.Replicas)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	log.Println("Inspected function - " + name)
	return ctx.NoContent(http.StatusOK)
}
