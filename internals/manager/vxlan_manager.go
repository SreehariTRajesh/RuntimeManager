package manager

import (
	"log"
	"runtime-manager/internals/lifecycle"
	"runtime-manager/internals/pkg"
	"runtime-manager/internals/utils"
)

type DeviceStatusEntry struct {
	Name string
	Type string
}

type VXLanNetwork struct {
	Devices []DeviceStatusEntry
}

func (vxlan_net *VXLanNetwork) Initialize() error {
	appConfig := utils.GetConfig()
	err := utils.ConfigureBridge(appConfig.VXLanConf.Bridge)
	if err != nil {
		vxlan_net.Cleanup()
		log.Fatalln("error while configuring bridge", err)
	}

	vxlan_net.Devices = append(vxlan_net.Devices, DeviceStatusEntry{Name: appConfig.VXLanConf.Bridge, Type: "bridge"})
	for _, peer := range appConfig.VXLanConf.VXLanPeers {
		err = utils.ConfigureVXLan(peer.Name, peer.VXLanId, peer.Device, peer.Remote, peer.DstPort, appConfig.VXLanConf.Bridge)
		if err != nil {
			vxlan_net.Cleanup()
			log.Fatalln("error while configuring vxlan", err)
		}
		vxlan_net.Devices = append(vxlan_net.Devices, DeviceStatusEntry{Name: peer.Name, Type: "vxlan"})
	}
	return nil
}

func (vxlan_net *VXLanNetwork) Cleanup() {
	for _, dev := range vxlan_net.Devices {
		if dev.Type == "bridge" {
			utils.DestroyBridge(dev.Name)
		} else if dev.Type == "vxlan" {
			utils.DestroyVXLan(dev.Name)
		}
	}
	clear(vxlan_net.Devices)
}

func (vxlan_net *VXLanNetwork) Order() int {
	return pkg.ORDER_1
}

func init() {
	vxlan_network := &VXLanNetwork{}
	lifecycle.RegisterInitializable(vxlan_network)
	lifecycle.RegisterCleanable(vxlan_network)
}
