package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexellis/faas/gateway/requests"
)

func request(path, method string, data interface{}) ([]byte, *http.Response, error) {
	res, err := json.Marshal(data)
	if err != nil {
		return nil, nil, err
	}
	input := bytes.NewBuffer(res)

	gatewayAddr := os.Getenv("FAAS_GATEWAY_ADDR")
	if gatewayAddr == "" {
		gatewayAddr = "localhost"
	}

	request, err := http.NewRequest(method, "http://"+gatewayAddr+":8080"+path, input)
	if err != nil {
		return nil, nil, err
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	return body, resp, nil
}

func TestBefore(t *testing.T) {
	deleteReq := requests.DeleteFunctionRequest{
		FunctionName: "nodeinfo",
	}

	_, _, err := request("/system/functions", "DELETE", deleteReq)
	if err != nil {
		t.Log(err)
	}
}

func TestDeploy(t *testing.T) {
	deploy := requests.CreateFunctionRequest{
		Image:      "functions/nodeinfo:latest",
		Service:    "nodeinfo",
		Network:    "func_functions",
		EnvProcess: "",
	}

	_, res, err := request("/system/functions", "POST", deploy)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	if res != nil && res.StatusCode != http.StatusCreated {
		t.Logf("got %d, wanted %d", res.StatusCode, http.StatusCreated)
		t.Fail()
	}

	time.Sleep(time.Second * 20) // Waiting for service scale ready
}

func TestList(t *testing.T) {
	data, res, err := request("/system/functions", "GET", nil)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	if res.StatusCode != http.StatusOK {
		t.Logf("got %d, wanted %d", res.StatusCode, http.StatusOK)
		t.Fail()
	}

	functions := []requests.Function{}
	err = json.Unmarshal(data, &functions)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	if len(functions) == 0 {
		t.Log("List functions got: 0, want: > 0")
		t.Fail()
	}
}

func TestInvoke(t *testing.T) {
	_, res, err := request("/function/nodeinfo", "POST", nil)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	if res.StatusCode != http.StatusOK {
		t.Logf("got %d, wanted %d", res.StatusCode, http.StatusOK)
		t.Fail()
	}
}

func TestDelete(t *testing.T) {
	deleteReq := requests.DeleteFunctionRequest{
		FunctionName: "nodeinfo",
	}

	_, res, err := request("/system/functions", "DELETE", deleteReq)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	if res.StatusCode != http.StatusOK {
		t.Logf("got %d, wanted %d", res.StatusCode, http.StatusOK)
		t.Fail()
	}
}
