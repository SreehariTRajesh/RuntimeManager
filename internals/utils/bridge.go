package utils

import (
	"fmt"
	"log"

	"github.com/vishvananda/netlink"
)

func ConfigureBridge(bridge_name string) error {
	link_attrs := netlink.NewLinkAttrs()
	link_attrs.Name = bridge_name
	bridge := &netlink.Bridge{LinkAttrs: link_attrs}

	if err := netlink.LinkAdd(bridge); err != nil {
		return fmt.Errorf("failed to create bridge: %w", err)
	}
	log.Printf("Bridge %s created successfully", bridge_name)

	if err := netlink.LinkSetUp(bridge); err != nil {
		_ = netlink.LinkDel(bridge)
		return fmt.Errorf("failed to set up bridge %s after creation: %w", bridge_name, err)
	}
	log.Printf("Bridge %s set up successfully", bridge_name)

	return nil
}

func DestroyBridge(bridge_name string) error {
	bridge, err := netlink.LinkByName(bridge_name)
	if err != nil {
		return fmt.Errorf("failed to find device with name %s: %w", bridge_name, err)
	}
	_ = netlink.LinkDel(bridge)
	return nil
}
