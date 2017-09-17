package hyper

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/alexellis/faas/gateway/requests"
	"github.com/hyperhq/hyper-api/client"
	"github.com/hyperhq/hyper-api/types"
)

type Hyper struct {
	*client.Client
	FuncMap map[string]string // function name -> service ip
}

type PrometheusQueryResponse struct {
	Status string
	Data   struct {
		ResultType string
		Result     []struct {
			Metric struct {
				FunctionName string `json:"function_name"`
			}
			Value []interface{}
		}
	}
}

func getInvocationCount(name string) float64 {
	res, err := http.Get("http://faas-prometheus:9090/api/v1/query?query=gateway_function_invocation_total")
	if err != nil {
		return 0
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	query := PrometheusQueryResponse{}
	if err := json.Unmarshal(body, &query); err != nil {
		return 0
	}
	for _, item := range query.Data.Result {
		if item.Metric.FunctionName == name {
			if len(item.Value) == 2 {
				count := item.Value[1].(string)
				if ret, err := strconv.ParseFloat(count, 64); err == nil {
					return ret
				}
			}
		}
	}

	return 0
}

func New() (*Hyper, error) {
	region := os.Getenv("HYPER_REGION")
	if region == "" {
		region = "us-west-1"
	}
	var (
		host          = "tcp://" + region + ".hyper.sh:443"
		customHeaders = map[string]string{}
		verStr        = "v1.23"
		accessKey     = os.Getenv("HYPER_ACCESS_KEY")
		secretKey     = os.Getenv("HYPER_SECRET_KEY")
	)

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	client, err := client.NewClient(host, verStr, httpClient, customHeaders, accessKey, secretKey)
	if err != nil {
		return nil, err
	}

	hyper := &Hyper{client, make(map[string]string)}

	if err = hyper.RefreshFuncMap(); err != nil {
		return nil, err
	}
	time.AfterFunc(10*time.Second, func() {
		hyper.RefreshFuncMap()
	})

	return hyper, nil
}

func (hyper *Hyper) Create(name, image string, envs []string, config map[string]string) error {
	size := "s4"
	if _, ok := config["hyper_size"]; ok {
		size = config["hyper_size"]
	}

	fullName := "faas-function-" + name
	service, err := hyper.ServiceCreate(
		context.Background(),
		types.Service{
			Name:          fullName,
			Image:         image,
			ContainerSize: size,
			ContainerPort: 8080,
			Replicas:      1,
			Protocol:      "http",
			ServicePort:   8080,
			Env:           envs,
			Labels: map[string]string{
				"faas-function": "true",
				fullName:        "true",
			},
		},
	)
	if err != nil {
		return err
	}
	hyper.FuncMap[name] = service.IP

	return nil
}

func (hyper *Hyper) RefreshFuncMap() error {
	services, err := hyper.ServiceList(context.Background(), types.ServiceListOptions{})
	if err != nil {
		return err
	}
	funcMap := make(map[string]string)
	for _, service := range services {
		isFunc := strings.HasPrefix(service.Name, "faas-function-")
		if !isFunc {
			continue
		}
		funcMap[service.Name] = service.IP
	}
	hyper.FuncMap = funcMap
	return nil
}

func (hyper *Hyper) List() ([]requests.Function, error) {
	services, err := hyper.ServiceList(context.Background(), types.ServiceListOptions{})
	if err != nil {
		return nil, err
	}

	functions := make([]requests.Function, 0)
	for _, service := range services {
		isFunc := strings.HasPrefix(service.Name, "faas-function-")
		if !isFunc {
			continue
		}
		name := strings.TrimPrefix(service.Name, "faas-function-")
		function := requests.Function{
			Name:            name,
			Replicas:        uint64(service.Replicas),
			Image:           service.Image,
			InvocationCount: getInvocationCount(name),
		}
		functions = append(functions, function)
	}

	return functions, nil
}

func (hyper *Hyper) Delete(name string) error {
	err := hyper.ServiceDelete(context.Background(), name, false)
	if err != nil {
		return err
	}
	return nil
}

func (hyper *Hyper) Inspect(name string) (*requests.Function, error) {
	functions, err := hyper.List()
	if err != nil {
		return nil, err
	}
	for _, function := range functions {
		if function.Name == name {
			return &function, nil
		}
	}
	return nil, nil
}

func (hyper *Hyper) Scale(name string, replica uint64) error {
	count := int(replica)
	_, err := hyper.ServiceUpdate(context.Background(), "faas-function-"+name, types.ServiceUpdate{
		Replicas: &count,
	})
	return err
}
