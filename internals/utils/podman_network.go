package utils

import (
	"context"
	"fmt"
	"net"

	"github.com/containers/common/libnetwork/types"
	"github.com/containers/podman/v5/pkg/bindings"
	"github.com/containers/podman/v5/pkg/bindings/network"
)

func CreatePodmanNetwork(network_name string, subnet string, gateway string, driver string) error {
	socket := "unix:///run/podman/podman.sock"
	ctx, err := bindings.NewConnection(context.Background(), socket)
	if err != nil {
		return fmt.Errorf("error connecting to podman socket '%s': %w", socket, err)
	}
	sub, err := types.ParseCIDR(subnet)
	if err != nil {
		return fmt.Errorf("error parsing subnet '%s': %w", subnet, err)
	}
	gw := net.ParseIP(gateway)
	opts := types.Network{
		Name:             "vxaln-overlay",
		Driver:           "bridge",
		NetworkInterface: "br0",
		Subnets: []types.Subnet{
			{
				Subnet:  sub,
				Gateway: gw,
			},
		},
	}
	report, err := network.Create(ctx, &opts)
	if err != nil {
		return fmt.Errorf("error creating network '%s': %w", network_name, err)
	}
	fmt.Printf("Network %s created successfully: %v\n", network_name, report)
	return nil
}
