// keeps track of all the containers on the host and destroys them when the SIGINT/SIGTERM comes
package manager

import (
	"runtime-manager/internals/lifecycle"
	"runtime-manager/internals/pkg"
	"runtime-manager/internals/utils"
)

type ContainerStatusRegistry struct {
	ContainerIds []string
}

func (registry *ContainerStatusRegistry) Add(container_id string) {
	registry.ContainerIds = append(registry.ContainerIds, container_id)
}

func (registry *ContainerStatusRegistry) Cleanup() {
	for _, id := range registry.ContainerIds {
		utils.DeleteContainer(id)
	}
}

func (registry *ContainerStatusRegistry) Order() int {
	return pkg.ORDER_0
}

func init() {
	container_registry := &ContainerStatusRegistry{}
	lifecycle.RegisterCleanable(container_registry)
}
