package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/containers/podman/v5/pkg/bindings"
	"github.com/containers/podman/v5/pkg/bindings/containers"
	"github.com/containers/podman/v5/pkg/bindings/images"
	"github.com/containers/podman/v5/pkg/specgen"
	"github.com/opencontainers/runtime-spec/specs-go"
)

func CreateContainerFunction(fn_name string, fn_bundle string, image string, cpu []int, mem int64, virt_ip string, mac string) (string, error) {
	socket := "unix:///run/podman/podman.sock"
	ctx, err := bindings.NewConnection(context.Background(), socket)
	if err != nil {
		return "", fmt.Errorf("error while connecting to podman socket: %w", err)
	}
	_, err = images.Pull(ctx, image, &images.PullOptions{})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	spec := specgen.NewSpecGenerator(image, false)
	spec.ResourceLimits = &specs.LinuxResources{
		CPU: &specs.LinuxCPU{
			Cpus: GetCoreSet(cpu),
		},
		Memory: &specs.LinuxMemory{
			Limit: &mem,
		},
	}
	res, err := containers.CreateWithSpec(ctx, spec, &containers.CreateOptions{})
	if err != nil {
		return "", fmt.Errorf("error while creating a container with spec: %w", err)
	}
	err = containers.Start(ctx, res.ID, &containers.StartOptions{})
	if err != nil {
		return "", fmt.Errorf("error while creating a container with spec: %w", err)
	}
	return res.ID, nil
}

func DeleteContainerFunction(container_id string) error {
	socket := "unix:///run/podman/podman.sock"
	ctx, err := bindings.NewConnection(context.Background(), socket)
	if err != nil {
		return fmt.Errorf("error while connecting to podman socket: %w", err)
	}
	err = containers.Kill(ctx, container_id, &containers.KillOptions{})
	if err != nil {
		return fmt.Errorf("error while trying to kill the container: %w", err)
	}
	_, err = containers.Remove(ctx, container_id, &containers.RemoveOptions{})
	if err != nil {
		return fmt.Errorf("error while trying to remove container: %w", err)
	}
	return nil
}

func ExecFunction(virtual_ip string, params map[string]any) (map[string]any, error) {

	return nil, nil
}

func MakeHttpRequest(ip string, port int, params map[string]any) (map[string]any, error) {
	host := fmt.Sprintf("%s:%d", ip, port)
	path := "/invoke"
	scheme := "http"
	baseURL := &url.URL{
		Scheme: scheme,
		Host:   host,
		Path:   path,
	}

	http_url := baseURL.String()
	request_body, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("error while marshalling the request body %v: %w", params, err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", http_url, bytes.NewBuffer(request_body))

	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("error making http request: %w", err)
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, fmt.Errorf("error while reading response: %w", err)
	}

	var response map[string]any

	err = json.Unmarshal(body, &response)

	if err != nil {
		return nil, fmt.Errorf("error while umarshalling json response: %w", err)
	}
	return response, nil
}

func GetCoreSet(cpus []int) string {
	var cpu_array []string
	for _, core := range cpus {
		cpu_array = append(cpu_array, strconv.Itoa(core))
	}
	return strings.Join(cpu_array, ",")
}
