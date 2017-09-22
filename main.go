package main

import (
	"log"

	"github.com/hyperhq/faas-hyper/handlers"
	"github.com/labstack/echo"
)

type Server struct {
	*echo.Echo
}

func main() {
	srv := Server{}
	srv.Echo = echo.New()

	hl, err := handlers.New()
	if err != nil {
		srv.Logger.Fatal(err)
	}

	srv.HTTPErrorHandler = func(err error, c echo.Context) {
		log.Printf("%s %s %s", c.Request().Method, c.Path(), err)
	}

	srv.GET("/system/functions", hl.List)
	srv.POST("/system/functions", hl.Deploy)
	srv.DELETE("/system/functions", hl.Delete)

	srv.GET("/system/function/:name", hl.Inspect)
	srv.POST("/system/scale-function/:name", hl.Scale)

	srv.Any("/function/:name", hl.Proxy)

	srv.Logger.Fatal(srv.Start(":8080"))
}
