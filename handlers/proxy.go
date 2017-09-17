package handlers

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"

	"io/ioutil"
)

func (hl *Handler) Proxy(ctx echo.Context) error {
	proxyClient := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   3 * time.Second,
				KeepAlive: 0,
			}).DialContext,
			MaxIdleConns:          1,
			DisableKeepAlives:     true,
			IdleConnTimeout:       120 * time.Millisecond,
			ExpectContinueTimeout: 1500 * time.Millisecond,
		},
	}
	service := ctx.Param("name")

	stamp := strconv.FormatInt(time.Now().Unix(), 10)

	defer func(when time.Time) {
		seconds := time.Since(when).Seconds()
		log.Printf("[%s] took %f seconds\n", stamp, seconds)
	}(time.Now())

	watchdogPort := 8080

	requestBody, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		return err
	}
	defer ctx.Request().Body.Close()

	serviceIP, found := hl.FuncMap["faas-function-"+service]
	if !found {
		return ctx.NoContent(http.StatusNotFound)
	}

	url := fmt.Sprintf("http://%s:%d/", serviceIP, watchdogPort)

	request, _ := http.NewRequest("POST", url, bytes.NewReader(requestBody))

	copyHeaders(&request.Header, &ctx.Request().Header)

	response, err := proxyClient.Do(request)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Can't reach service: "+service)
		return err
	}

	clientHeader := ctx.Response().Header()
	copyHeaders(&clientHeader, &response.Header)

	responseBody, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	return ctx.Blob(response.StatusCode, response.Header.Get("Content-Type"), responseBody)
}

func copyHeaders(destination *http.Header, source *http.Header) {
	for k, vv := range *source {
		vvClone := make([]string, len(vv))
		copy(vvClone, vv)
		(*destination)[k] = vvClone
	}
}
