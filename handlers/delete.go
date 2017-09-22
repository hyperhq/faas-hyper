package handlers

import (
	"log"
	"net/http"

	"github.com/alexellis/faas/gateway/requests"
	"github.com/labstack/echo"
)

func (hl *Handler) Delete(ctx echo.Context) error {
	opts := new(requests.DeleteFunctionRequest)
	if err := ctx.Bind(opts); err != nil {
		ctx.NoContent(http.StatusBadRequest)
		return err
	}
	err := hl.Hyper.Delete("faas-function-" + opts.FunctionName)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return err
	}
	log.Println("Deleted function - " + opts.FunctionName)
	return ctx.NoContent(http.StatusOK)
}
