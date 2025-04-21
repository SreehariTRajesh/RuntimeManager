package utils

import (
	"fmt"
	"log"
	"net"
	"runtime-manager/internals/pkg"

	"github.com/vishvananda/netlink"
)

func ConfigureVXLan(vxlan_name string, vxlan_id int, device_name string, remote_ip_str string, remote_port int, bridge_name string) error {
	device, err := netlink.LinkByName(device_name)
	if err != nil {
		return fmt.Errorf("failed to get device %s: %w", device_name, err)
	}
	remote_ip := net.ParseIP(remote_ip_str)

	if remote_ip == nil {
		return fmt.Errorf("invalid ip address: %s", remote_ip_str)
	}
	// Convert to IPv4
	remote_ip = remote_ip.To4()

	if remote_ip == nil {
		return fmt.Errorf("remote IP address is not an IPv4 address: %s", remote_ip_str)
	}

	vx_lan_attrs := netlink.NewLinkAttrs()
	vx_lan_attrs.Name = vxlan_name
	vx_lan := &netlink.Vxlan{
		LinkAttrs:    vx_lan_attrs,
		VxlanId:      vxlan_id,
		Group:        remote_ip,
		SrcAddr:      nil,
		Port:         remote_port,
		Learning:     true,
		L2miss:       true,
		L3miss:       true,
		VtepDevIndex: device.Attrs().Index,
	}

	log.Printf("Adding VXLan Link: %s (VNI: %d, Remote: %s, Port: %d, Dev:%s)", vxlan_name, vxlan_id, remote_ip_str, remote_port, device_name)

	if err := netlink.LinkAdd(vx_lan); err != nil {
		if _, exists_error := netlink.LinkByName(vxlan_name); exists_error == nil {
			log.Printf("VXLAN link %s already exists, proceeding...", vxlan_name)
		} else {
			return fmt.Errorf("failed to add vxlan link %s: %w", vxlan_name, err)
		}
	}

	// setting vxlan link up
	// Get the vxlan name
	vxlan_link, err := netlink.LinkByName(vxlan_name)

	if err != nil {
		return fmt.Errorf("failed to get link %s after creation: %w", vxlan_name, err)
	}

	log.Printf("Setting link %s up", vxlan_name)

	if err := netlink.LinkSetUp(vxlan_link); err != nil {
		return fmt.Errorf("failed to set link %s up: %w", vxlan_name, err)
	}

	bridge, err := netlink.LinkByName(bridge_name)

	if err != nil {
		return fmt.Errorf("failed to get bridge %s: %w", bridge_name, err)
	}

	if bridge.Type() != pkg.BRIDGE {
		return fmt.Errorf("device %s is not a bridge", bridge_name)
	}

	// setting bridge br0 to be master of VXLAN

	if err := netlink.LinkSetMaster(vxlan_link, bridge); err != nil {
		return fmt.Errorf("failed to add link %s to bridge %s: %w", vxlan_name, bridge_name, err)
	}

	log.Printf("Successfully configured VXLAN interface:%s and added to bridge: %s", vxlan_name, bridge_name)
	return nil
}

func DestroyVXLan(vxlan_name string) error {
	link, err := netlink.LinkByName(vxlan_name)

	if err != nil {
		return fmt.Errorf("failed to find link %s for deletion: %w", vxlan_name, err)
	}

	if err := netlink.LinkDel(link); err != nil {
		return fmt.Errorf("failed to delete link %s: %w", vxlan_name, err)
	}

	log.Printf("Successfully destroyed VXLAN link: %s", vxlan_name)
	return nil
}
