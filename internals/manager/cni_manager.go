package manager

import (
	"fmt"
	"log"
	"runtime-manager/internals/lifecycle"
	"runtime-manager/internals/pkg"
	"runtime-manager/internals/utils"
)

type CNINetwork struct {
	Name      string
	Type      string
	Subnet    string
	Interface string
}

func (cni *CNINetwork) Initialize() error {
	appConfig := utils.GetConfig()
	network_name := appConfig.CNI.Name
	interface_name := appConfig.CNI.Interface
	subnet := appConfig.CNI.Subnet
	gateway := appConfig.CNI.Gateway
	driver := appConfig.CNI.Driver
	err := utils.CreatePodmanNetwork(network_name, subnet, gateway, driver, interface_name)
	if err != nil {
		return fmt.Errorf("error while creating macvlan network %s: %w", network_name, err)
	}
	return nil
}

func (cni *CNINetwork) Cleanup() {
	err := utils.DestroyPodmanNetwork(cni.Name)
	if err != nil {
		log.Fatalf("error while destroying macvlan network %s: %v", cni.Name, err)
	}
}

func (cni *CNINetwork) Order() int {
	return pkg.ORDER_2
}

func init() {
	cni := &CNINetwork{}
	lifecycle.RegisterInitializable(cni)
	lifecycle.RegisterCleanable(cni)
}
