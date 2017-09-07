package handlers

import (
	"log"
	"net/http"

	"github.com/labstack/echo"
)

func (hl *Handler) Delete(ctx echo.Context) error {
	name := ctx.Param("name")
	err := hl.Hyper.Delete("faas-function-" + name)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	log.Println("Deleted function - " + name)
	return ctx.NoContent(http.StatusNoContent)
}
