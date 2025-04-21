package utils

import (
	"fmt"
	"log"
	"runtime-manager/internals/pkg"

	"github.com/vishvananda/netlink"
)

func CreateMacVLANNetworkInterface(interface_name string, ip_addr string, parent_device_name string) error {
	parent_link, err := netlink.LinkByName(parent_device_name)

	if err != nil {
		return fmt.Errorf("failed to find parent link for device %s: %w", parent_device_name, err)
	}
	log.Println("found parent link for device ", parent_device_name)

	macvlan_interface := &netlink.Macvlan{
		LinkAttrs: netlink.LinkAttrs{
			Name:        interface_name,
			ParentIndex: parent_link.Attrs().Index,
		},
		Mode: netlink.MACVLAN_MODE_BRIDGE,
	}
	err = netlink.LinkAdd(macvlan_interface)
	if err != nil {
		if err.Error() != pkg.FILE_EXISTS {
			return fmt.Errorf("error adding macvlan link %s: %w", interface_name, err)
		}
	}
	log.Printf("successfully added macvlan link %s", interface_name)
	// fetch link by name
	macvlan_link, err := netlink.LinkByName(interface_name)
	if err != nil {
		return fmt.Errorf("error fetching macvlan link by name %s: %w", interface_name, err)
	}

	addr, err := netlink.ParseAddr(ip_addr)
	if err != nil {
		DestroyMacVlanNetworkInterface(macvlan_link)
		return fmt.Errorf("error parsing ip address %s: %w", ip_addr, err)
	}

	err = netlink.AddrAdd(macvlan_link, addr)
	if err != nil {
		DestroyMacVlanNetworkInterface(macvlan_link)
		return fmt.Errorf("error adding address to link %s: %w", ip_addr, err)
	}

	err = netlink.LinkSetUp(macvlan_link)
	if err != nil {
		DestroyMacVlanNetworkInterface(macvlan_link)
		return fmt.Errorf("error setting up the macvlan link %s: %w", interface_name, err)
	}
	log.Println("link set up successfully")
	return nil
}

func DestroyMacVlanNetworkInterface(macvlan_link netlink.Link) error {
	err := netlink.LinkDel(macvlan_link)
	if err != nil {
		return fmt.Errorf("error destroying macvlan link %s: %w", macvlan_link.Attrs().Name, err)
	}
	return nil
}
