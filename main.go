package main

import (
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

	srv.GET("/system/functions", hl.List)
	srv.POST("/system/functions", hl.Deploy)
	srv.DELETE("/system/functions/:name", hl.Delete)

	srv.GET("/system/function/:name", hl.Inspect)
	srv.POST("/system/scale-function/:name", hl.Scale)

	srv.Any("/function/:name", hl.Proxy)

	srv.Logger.Fatal(srv.Start(":8080"))
}
