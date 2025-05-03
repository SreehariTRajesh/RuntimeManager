package utils

import (
	"context"
	"fmt"
	"log"
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
	log.Println("getting new connection")
	ctx, err := bindings.NewConnection(context.Background(), socket)
	log.Println("got new connection")
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

func GetCoreSet(cpus []int) string {
	var cpu_array []string
	for _, core := range cpus {
		cpu_array = append(cpu_array, strconv.Itoa(core))
	}
	return strings.Join(cpu_array, ",")
}
