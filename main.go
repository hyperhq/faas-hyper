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

	srv.Any("/function/:name", hl.Proxy)

	srv.POST("/system/functions", hl.Deploy)
	srv.GET("/system/functions", hl.List)
	srv.GET("/system/functions/:name", hl.Inspect)
	srv.DELETE("/system/functions/:name", hl.Delete)

	srv.Logger.Fatal(srv.Start(":8080"))
}
