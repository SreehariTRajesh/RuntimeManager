package status_table

import "sync"

type MigrationStatusTableEntry struct {
	ContainerIP  string
	FunctionName string
	ContainerId  string
	Occupied     bool
	sync.Mutex
}
