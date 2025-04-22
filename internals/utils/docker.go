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
	"os/exec"
	"runtime-manager/internals/pkg"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
)

// Following class contains utils for docker
func CreateMacVLANNetwork(network_name string, subnet string, gateway string, parent_device_name string) (*types.NetworkCreateResponse, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("error creating docker client: %w", err)
	}
	defer cli.Close()

	ipamConfig := &network.IPAMConfig{
		Subnet:  subnet,
		Gateway: gateway,
	}

	ipam := &network.IPAM{
		Driver: "default",
		Config: []network.IPAMConfig{*ipamConfig},
	}

	networkOptions := types.NetworkCreate{
		CheckDuplicate: true,
		Driver:         "macvlan",
		Options: map[string]string{
			"parent": parent_device_name,
		},
		IPAM: ipam,
	}

	resp, err := cli.NetworkCreate(ctx, network_name, networkOptions)
	if err != nil {
		return nil, fmt.Errorf("error creating macvlan network %s: %w", network_name, err)
	}
	return &resp, nil
}

// Assign IP Configuration and create containers

func DestroyMacVLANNetwork(network_id string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return fmt.Errorf("error creating docker client: %w", err)
	}
	err = cli.NetworkRemove(ctx, network_id)
	if err != nil {
		return fmt.Errorf("error removing the network from host: %w", err)
	}
	return nil
}

func CreateAndStartContainer(image_name string, cores []int, memory int, virtual_ip string, network_name string, function_bundle string) (string, error) {
	// functionality to start a container
	/*
			@param: cpu_cores
			@param: memory
			@param: virtual_ip
			@param: function_bundle
			# default port assigned 80if err != nil {
			return nil, fmt.Errorf("error creating docker client: %w", err)
		}
		defer cli.Close()
	*/
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return "", fmt.Errorf("error creating docker client: %w", err)
	}

	_, err = cli.ImagePull(ctx, image_name, types.ImagePullOptions{})

	if err != nil {
		return "", fmt.Errorf("error pulling image %s: %w", image_name, err)
	}

	// Get the network name from the network_manager
	//network_name := macvlan.GetNetworkname()
	cpuSet := getCoreSet(cores)

	resource_config := container.Resources{
		Memory:     int64(memory),
		CpusetCpus: cpuSet,
	}

	network_config := &network.EndpointSettings{
		IPAddress: virtual_ip,
	}

	exposed_ports := nat.PortSet{
		nat.Port(pkg.DEFAULT_CONTAINER_PORT): struct{}{},
	}

	networking_config := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			network_name: network_config,
		},
	}

	host_config := &container.HostConfig{
		Resources: resource_config,
		Binds: []string{
			fmt.Sprintf("%s:/app", function_bundle),
		},
	}

	container_config := &container.Config{
		Image:        image_name,
		Cmd:          []string{"sh", "-c", "chmod +x /app/startup.sh && /app/startup.sh"},
		ExposedPorts: exposed_ports,
	}

	// create the container
	res, err := cli.ContainerCreate(ctx, container_config, host_config, networking_config, nil, "")

	if err != nil {
		return "", fmt.Errorf("error while creating container: %w", err)
	}

	endpoint_config := &network.EndpointSettings{
		IPAddress: virtual_ip,
	}

	err = cli.NetworkConnect(ctx, network_name, res.ID, endpoint_config)
	if err != nil {
		log.Fatalf("Failed to connect to network: %v", err)
	}

	//start the container
	if err := cli.ContainerStart(ctx, res.ID, container.StartOptions{}); err != nil {
		return "", fmt.Errorf("error starting the container: %w", err)
	}
	return res.ID, nil
}

func InvokeFunction(container_ip string, params map[string]any) (map[string]any, error) {
	// functionality to execute container functions
	return MakeHttpRequest(container_ip, 80, params)
}

// set of util functions;
func getCoreSet(cores []int) string {
	core_array := make([]string, len(cores))
	for i, core := range cores {
		core_array[i] = strconv.Itoa(core)
	}
	return strings.Join(core_array, ",")
}

// assign a name for the container -> name must be unique
func AssignName() string {
	unique_id := uuid.New().String()
	return unique_id[:16]
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

func MigrateContainer(source_ip string, dest_ip string, container_id string, checkpoint_dir string, image_name string) error {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv)

	if err != nil {
		return fmt.Errorf("error creating a docker client: %w", err)
	}

	checkPointOpts := types.CheckpointCreateOptions{
		CheckpointID: "migration-checkpoint",
		Exit:         true,
	}

	if err := cli.CheckpointCreate(ctx, container_id, checkPointOpts); err != nil {
		return fmt.Errorf("error creating checkpoint for container %s: %w", container_id, err)
	}

	log.Println("checkpoint created for container: ", container_id)
	err = TransferCheckpointFiles(checkpoint_dir, dest_ip)

	if err != nil {
		return fmt.Errorf("error transferring checkpoint files to destination node: %w", err)
	}

	log.Println("checkpoint files transferred to destination node:", dest_ip)

	err = StartMigratedContainerOnRemoteHost(dest_ip, container_id, image_name)
	if err != nil {
		return fmt.Errorf("error starting migrated container on remote host: %w", err)
	}
	return nil
}

func TransferCheckpointFiles(checkpoint_dir string, dest_ip string) error {
	cmd := fmt.Sprintf("scp -r %s %s:%s", checkpoint_dir, dest_ip, checkpoint_dir)
	if err := ExecCommand(cmd); err != nil {
		return fmt.Errorf("error transferring checkpoint files: %w", err)
	}
	return nil
}

func ExecCommand(command string) error {
	out, err := exec.Command("bash", "-c", command).CombinedOutput()
	if err != nil {
		return fmt.Errorf("command execution failed: %s, output: %s", err, string(out))
	}
	log.Println("command output:", string(out))
	return nil
}

func StartMigratedContainerOnRemoteHost(remote_ip string, container_id string, image_name string) error {

	ctx := context.Background()
	remote_addr := fmt.Sprintf("tcp://%s:2375", remote_ip)
	cli, err := client.NewClientWithOpts(client.WithHost(remote_addr))
	if err != nil {
		return fmt.Errorf("error creating docker client for remote host: %w", err)
	}
	// check if image exists
	_, err = cli.ImagePull(ctx, image_name, types.ImagePullOptions{})

	if err != nil {
		return fmt.Errorf("error pulling image %s on remote host: %w", image_name, err)
	}

	if err := cli.ContainerStart(ctx, container_id, container.StartOptions{}); err != nil {
		return fmt.Errorf("error starting the container on remote host: %w", err)
	}
	return nil
}
