package utils

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
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
		# default port assigned 80
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

	networking_config := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			network_name: network_config,
		},
	}

	container_config := &container.Config{
		Image: image_name,
	}

	host_config := &container.HostConfig{
		Resources: resource_config,
	}

	// create the container
	res, err := cli.ContainerCreate(ctx, container_config, host_config, networking_config, nil, "")

	if err != nil {
		return "", fmt.Errorf("error while creating container ")
	}

	//start the container
	if err := cli.ContainerStart(ctx, res.ID, container.StartOptions{}); err != nil {
		return "", fmt.Errorf("error starting the container: %w", err)
	}
	return res.ID, nil
}

func ExecCommandContainer() {
	// functionality to execute container functions
}

func getCoreSet(cores []int) string {
	core_array := make([]string, len(cores))
	for i, core := range cores {
		core_array[i] = strconv.Itoa(core)
	}
	return strings.Join(core_array, ",")
}

func AssignName() {
	// assign a name for the container -> name must be unique
}
