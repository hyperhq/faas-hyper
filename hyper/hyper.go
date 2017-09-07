package hyper

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"strings"

	"github.com/alexellis/faas/gateway/requests"
	"github.com/docker/go-connections/nat"
	"github.com/hyperhq/hyper-api/client"
	"github.com/hyperhq/hyper-api/types"
	"github.com/hyperhq/hyper-api/types/container"
	"github.com/hyperhq/hyper-api/types/filters"
	"github.com/hyperhq/hyper-api/types/network"
)

type Hyper struct {
	*client.Client
}

func New() (*Hyper, error) {
	var (
		host          = "tcp://us-west-1.hyper.sh:443"
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

	return &Hyper{client}, nil
}

func (hyper *Hyper) Create(name, image string, envs []string) error {
	hostName := "faas-function-" + name
	res, err := hyper.ContainerCreate(
		context.Background(),
		&container.Config{
			Image:    image,
			Hostname: hostName,
			Env:      envs,
			Labels: map[string]string{
				"sh_hyper_instancetype": "s4",
				"faas-function":         "true",
			},
			ExposedPorts: map[nat.Port]struct{}{
				"8080/tcp": {},
			},
		},
		&container.HostConfig{
			PortBindings: nat.PortMap{
				"8080/tcp": []nat.PortBinding{
					nat.PortBinding{HostPort: "8080"},
				},
			},
			RestartPolicy: container.RestartPolicy{
				Name: "always",
			},
		},
		&network.NetworkingConfig{},
		hostName,
	)
	if err != nil {
		return err
	}

	return hyper.ContainerStart(context.Background(), res.ID, "")
}

func (hyper *Hyper) List() ([]requests.Function, error) {
	args := filters.NewArgs()
	args.Add("label", "faas-function=true")

	containers, err := hyper.ContainerList(context.Background(), types.ContainerListOptions{
		All:    true,
		Filter: args,
	})
	if err != nil {
		return nil, err
	}

	functions := make([]requests.Function, 0)
	for _, container := range containers {
		var replicas uint64

		function := requests.Function{
			Name:            strings.TrimPrefix(container.Names[0], "/faas-function-"),
			Replicas:        replicas,
			Image:           container.Image,
			InvocationCount: 0,
		}
		functions = append(functions, function)
	}

	return functions, nil
}

func (hyper *Hyper) Delete(name string) error {
	args := filters.NewArgs()
	args.Add("v", "1")
	args.Add("force", "1")

	_, err := hyper.ContainerRemove(context.Background(), name, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		RemoveLinks:   true,
		Force:         true,
	})
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
