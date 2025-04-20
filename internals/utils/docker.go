package utils

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
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

func CreateAndStartContainer(cores []int, memory int, vittural_ip string, function_bundle string) {
	// functionality to start a container
	/*
		@param: cpu_cores
		@param: memory
		@param: virtual_ip
		@param: function_bundle
		# default port assigned 80
	*/

}

func ExecCommandContainer() {
	// functionality to execute container functions
}
