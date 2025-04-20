package macvlan

import (
	"fmt"
	"runtime-manager/configs"
	"runtime-manager/internals/lifecycle"
	"runtime-manager/internals/pkg"
	"runtime-manager/internals/utils"
)

// following module initialiases the docker network
type MacVLANNetwork struct {
	Name      string
	Subnet    string
	Gateway   string
	Parent    string
	NetworkId string
}

var appConfig *configs.Config

func (network *MacVLANNetwork) Initialize() error {
	network_name := appConfig.MacVLANNetworkConf.Name
	network_gateway := appConfig.MacVLANNetworkConf.Gateway
	network_parent := appConfig.MacVLANNetworkConf.Parent
	network_subnet := appConfig.MacVLANNetworkConf.Subnet
	resp, err := utils.CreateMacVLANNetwork(network_name, network_subnet, network_gateway, network_parent)
	if err != nil {
		return fmt.Errorf("error while creating macvlan network %s: %w", network.Name, err)
	}
	network.Name = network_name
	network.Gateway = network_gateway
	network.Parent = network_parent
	network.Subnet = network_subnet
	network.NetworkId = resp.ID
	return nil
}

func (network *MacVLANNetwork) Cleanup() {
	utils.DestroyMacVLANNetwork(network.NetworkId)
}

func (network *MacVLANNetwork) Order() int {
	return pkg.ORDER_1
}

func init() {
	appConfig = configs.Parser(pkg.CONFIG_FILE_PATH)
	macvlan_network := &MacVLANNetwork{}
	lifecycle.RegisterInitializable(macvlan_network)
	lifecycle.RegisterCleanable(macvlan_network)
}
