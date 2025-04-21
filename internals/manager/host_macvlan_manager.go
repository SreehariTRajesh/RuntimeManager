package manager

import (
	"fmt"
	"runtime-manager/internals/lifecycle"
	"runtime-manager/internals/pkg"
	"runtime-manager/internals/utils"
)

type HostMacVLANInterface struct {
	Name      string
	IPAddr    string
	Parent    string
	NetworkId string
}

func (iface *HostMacVLANInterface) Initialize() error {
	appConfig := utils.GetConfig()
	interface_name := appConfig.HostMacVLANInterfaceConf.Name
	parent_device := appConfig.HostMacVLANInterfaceConf.ParentInterface
	ip_addr := appConfig.HostMacVLANInterfaceConf.IPAddress
	err := utils.CreateMacVLANNetworkInterface(interface_name, ip_addr, parent_device)
	if err != nil {
		return fmt.Errorf("error while creating macvlan network %s: %w", interface_name, err)
	}
	iface.Name = interface_name
	iface.Parent = parent_device
	iface.IPAddr = ip_addr
	return nil
}

func (iface *HostMacVLANInterface) Cleanup() {
	utils.DestroyMacVLANNetworkInterfaceByName(iface.Name)
}

func (iface *HostMacVLANInterface) Order() int {
	return pkg.ORDER_2
}

func init() {
	macvlan_interface := &HostMacVLANInterface{}
	lifecycle.RegisterInitializable(macvlan_interface)
	lifecycle.RegisterCleanable(macvlan_interface)
}
