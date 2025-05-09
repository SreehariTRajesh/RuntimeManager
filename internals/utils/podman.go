package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/containers/common/libnetwork/types"
	"github.com/containers/podman/v5/pkg/bindings"
	"github.com/containers/podman/v5/pkg/bindings/containers"
	"github.com/containers/podman/v5/pkg/bindings/images"
	"github.com/containers/podman/v5/pkg/specgen"

	"github.com/opencontainers/runtime-spec/specs-go"
)

func CreateContainerFunction(fn_name string, fn_bundle string, image string, cpu []int, mem int64) (string, string, error) {
	socket := "unix:///run/podman/podman.sock"
	ctx, err := bindings.NewConnection(context.Background(), socket)
	if err != nil {
		return "", "", fmt.Errorf("error while connecting to podman socket: %w", err)
	}
	imageExists, err := images.Exists(ctx, image, &images.ExistsOptions{})
	if err != nil {
		return "", "", fmt.Errorf("error while checking if image exists: %w", err)
	}
	if !imageExists {
		log.Printf("image %s does not exist, pulling the image", image)
		_, err = images.Pull(ctx, image, &images.PullOptions{})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
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
	if err != nil {
		return "", "", fmt.Errorf("error while parsing mac address: %w", err)
	}

	spec.Networks = map[string]types.PerNetworkOptions{
		"vxlan-overlay": {},
	}

	res, err := containers.CreateWithSpec(ctx, spec, &containers.CreateOptions{})
	if err != nil {
		return "", "", fmt.Errorf("error while creating a container with spec: %w", err)
	}
	err = containers.Start(ctx, res.ID, &containers.StartOptions{})
	if err != nil {
		return "", "", fmt.Errorf("error while creating a container with spec: %w", err)
	}

	inspect, err := containers.Inspect(ctx, res.ID, &containers.InspectOptions{})
	if err != nil {
		return "", "", fmt.Errorf("error while inspecting the container: %w", err)
	}
	network := inspect.NetworkSettings
	return res.ID, network.IPAddress, nil
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

func MigrateContainerFunction(container_id string, source_ip string, destination_ip string) error {
	socket := "unix:///run/podman/podman.sock"
	ctx, err := bindings.NewConnection(context.Background(), socket)
	if err != nil {
		return fmt.Errorf("error while connecting to podman socket: %w", err)
	}
	tcp_established := true
	leave_running := false
	export_file_path := fmt.Sprintf("/tmp/%s.tar.gz", container_id)
	_, err = containers.Checkpoint(ctx, container_id, &containers.CheckpointOptions{
		TCPEstablished: &tcp_established,
		LeaveRunning:   &leave_running,
		Export:         &export_file_path,
	})
	if err != nil {
		return fmt.Errorf("error creating checkpoint: %w", err)
	}
	err = TransferCheckpointFiles(container_id, export_file_path, destination_ip)
	if err != nil {
		return fmt.Errorf("error transferring checkpoint fiels: %w", err)
	}

	_, err = containers.Remove(ctx, container_id, &containers.RemoveOptions{})
	if err != nil {
		return fmt.Errorf("error while trying to remove container: %w", err)
	}
	return nil
}

func TransferCheckpointFiles(container_id string, checkpoint_dir string, destination_ip string) error {
	user := "root"
	password := "root"
	return TransferFilesSFTP(container_id, destination_ip, user, password)
}

func ExecFunction(virtual_ip string, params map[string]any) (map[string]any, error) {
	return MakeHttpRequest(virtual_ip, 80, params)
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

func StartMigratedContainer(container_id string) (string, error) {
	socket := "unix:///run/podman/podman.sock"
	ctx, err := bindings.NewConnection(context.Background(), socket)
	if err != nil {
		return "", fmt.Errorf("error while connecting to podman socket: %w", err)
	}
	import_archive := fmt.Sprintf("/tmp/%s.tar.gz", container_id)
	_, err = containers.Restore(ctx, container_id, &containers.RestoreOptions{
		ImportArchive: &import_archive,
	})
	if err != nil {
		return "", fmt.Errorf("error while restoring container from checkpoint: %w", err)
	}
	return container_id, nil
}
