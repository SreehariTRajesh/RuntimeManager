package utils

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/containers/common/libnetwork/types"
	"github.com/containers/podman/v5/pkg/bindings"
	"github.com/containers/podman/v5/pkg/bindings/containers"
	"github.com/containers/podman/v5/pkg/specgen"
	"github.com/opencontainers/runtime-spec/specs-go"
)

func CreateContainerFunction(fn_name string, fn_bundle string, image string, cpu []int, mem int64, virt_ip string, mac string) (string, error) {
	socket := "unix:///run/podman/podman.sock"
	ctx := context.Background()
	_, err := bindings.NewConnection(ctx, socket)
	if err != nil {
		return "", fmt.Errorf("error while connecting to podman socket: %w", err)
	}
	fmt.Println("Connected to socket successfully")
	spec := specgen.NewSpecGenerator("python:3.9-slim", true)
	*spec.Terminal = true
	spec.ResourceLimits.CPU.Cpus = GetCoreSet(cpu)
	spec.ResourceLimits.Memory = &specs.LinuxMemory{Limit: &mem}
	spec.Networks["vxlan-overlay"] = types.PerNetworkOptions{
		StaticIPs: []net.IP{net.ParseIP(virt_ip)},
		StaticMAC: types.HardwareAddr(mac),
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
	sock_dir := os.Getenv("XDG_RUNTIME_DIR")
	socket := "unix:" + sock_dir + "/podman/podman.sock"
	ctx := context.Background()
	_, err := bindings.NewConnection(ctx, socket)
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

func GetCoreSet(cpus []int) string {
	var cpu_array []string
	for _, core := range cpus {
		cpu_array = append(cpu_array, strconv.Itoa(core))
	}
	return strings.Join(cpu_array, ",")
}
