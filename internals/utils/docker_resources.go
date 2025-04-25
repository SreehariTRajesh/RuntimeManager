package utils

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func UpdateCorePool(core_pool []int, container_id string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return fmt.Errorf("error creating docker client: %w", err)
	}
	defer cli.Close()
	// update core pool of the container through this function
	update_config := container.UpdateConfig{
		Resources: container.Resources{
			CpusetCpus: getCoreSet(core_pool),
		},
	}

	_, err = cli.ContainerUpdate(ctx, container_id, update_config)

	if err != nil {
		return fmt.Errorf("error updating container configs: %w", err)
	}
	return nil
}

func UpdateMemory(memory int64, container_id string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return fmt.Errorf("error creating docker client: %w", err)
	}
	defer cli.Close()
	update_config := container.UpdateConfig{
		Resources: container.Resources{
			Memory: memory,
		},
	}
	_, err = cli.ContainerUpdate(ctx, container_id, update_config)

	if err != nil {
		return fmt.Errorf("error updating container configs: %w", err)
	}
	return nil
}
