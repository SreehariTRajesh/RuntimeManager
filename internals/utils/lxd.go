package utils

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/lxc/go-lxc"
)

func CreateAndStartContainerLXD(container_name string, image_name string, cores []int, memory int, virtual_ip string, network_name string, static_mac string, function_bundle string, vxlan_bridge_name string, gateway string) (string, error) {
	c, err := lxc.NewContainer(container_name, lxc.DefaultConfigPath())
	if err != nil {
		return "", fmt.Errorf("error initializing the container")
	}
	defer c.Release()

	log.Println("Creating the container...")

	var backend lxc.BackendStore

	if err := (&backend).Set("zfs"); err != nil {
		c.Destroy()
		return "", fmt.Errorf("error settign backend store type", err)
	}

	bdev_size, err := lxc.ParseBytes("100M")
	if err != nil {
		c.Destroy()
		return "", fmt.Errorf("error parsing fssize:", err)
	}

	cpu := CPUCoreSet(cores)
	mem := Memory(memory)

	options := lxc.TemplateOptions{
		Template:             "download",
		Distro:               "ubuntu",
		Release:              "focal",
		Arch:                 "amd64",
		FlushCache:           false,
		DisableGPGValidation: true,
		Backend:              backend,
		BackendSpecs: &lxc.BackendStoreSpecs{
			FSSize: uint64(bdev_size),
			FSType: "zfs",
		},
	}

	if err := c.SetConfigItem("lxc.cgroup.cpuset.cpus", cpu); err != nil {
		// failed to set cpu cores
		log.Println("faied to set cpu cores", err)
		return "", fmt.Errorf("error while trying to set cpu cores: %w", cores)
	}
	log.Println("set cpu cores to: ", cores)

	if err := c.SetConfigItem("lxc.cgroup.memory.limit_in_bytes", mem); err != nil {
		log.Println("failed to set memory limit:", err)
		return "", fmt.Errorf("error while trying to set memory limit: %w", err)
	}
	log.Println("set memory limit to: %s\n", mem)

	if err := c.SetConfigItem("lxc.net.0.type", "veth"); err != nil {
		log.Println("failed to set network type:", err)
		return "", fmt.Errorf("error while trying to set network type: %w", err)
	}

	if err := c.SetConfigItem("lxc.net.0.link", vxlan_bridge_name); err != nil {
		log.Println("failed to set network link:", vxlan_bridge_name)
		return "", fmt.Errorf("error while trying to set network link: %w", err)
	}

	if err := c.SetConfigItem("lxc.net.0.flags", "up"); err != nil {
		log.Println("failed to set network flags:", err)
		return "", fmt.Errorf("error while trying to set network flags: %w", err)
	}

	if err := c.SetConfigItem("lxc.net.0.ipv4.gateway", gateway); err != nil {
		log.Println("failed to set gatway:", err)
		return "", fmt.Errorf("error while trying to set network gateway: %w", err)
	}

	if err := c.SetConfigItem("lxc.net.0.hwaddr", static_mac); err != nil {
		log.Println("failed to set mac address:", err)
		return "", fmt.Errorf("error while trying to set mac address: %w", err)
	}

	mountEntry := fmt.Sprintf("%s %s none bind 0 0", function_bundle, "app")

	if err := c.SetConfigItem("lxc.mount.entry", mountEntry); err != nil {
		c.Destroy()
		return "", fmt.Errorf("failed to set mount entry `%s`: %w", mountEntry, err)
	}
	// creating container
	if err := c.Create(options); err != nil {
		log.Println("error while creating container with template options:", err)
		c.Destroy()
		return "", fmt.Errorf("error while trying to create container with template options: %w", err)
	}

	log.Println("container %s created successfully", c.Name())

	log.Println("starting the container: %s", c.Name())

	if err := c.Start(); err != nil {
		log.Fatalf("failed to start container: %v\n", err)
	}

	log.Println("container: %s started \n", c.Name())
	// wait for network
	log.Println("waiting for container to gett ip address..")

	ip, err := WaitForIPv4(c, 2*time.Second, virtual_ip)

	if err != nil {
		log.Println("error while trying to configure ip address")
		return "", fmt.Errorf("error while configuring ip address: %v", ip)
	}

	install_cmd := []string{
		"sh", "-c",
		"export DEBIAN_FRONTEND=noninteractive, apt-get update && apt-get install -y python3 python3-pip",
	}
	attach_options := lxc.AttachOptions{
		StdinFd:    os.Stdin.Fd(),
		StdoutFd:   os.Stdout.Fd(),
		StderrFd:   os.Stderr.Fd(),
		ClearEnv:   false,
		EnvToKeep:  []string{"PATH", "HOME", "LANG", "LC_ALL", "DEBIAN_FRONTEND"},
		UID:        0,
		GID:        0,
		Cwd:        "/",
		Namespaces: syscall.CLONE_NEWNS,
	}

	exit_code, err := c.RunCommand(install_cmd, attach_options)

	if err != nil {
		log.Println("error running install commmand: %v (exit code: %d)", err, exit_code)
	}

	server_cmd := []string{"sh", "-c", "cd /app && python3 server.py"}

	exit_code, err = c.RunCommand(server_cmd, attach_options)

	if err != nil {
		log.Println("error running server command:", err, " exitcode:", exit_code)
		return "", fmt.Errorf("error while running server command : %w", err)
	} else if exit_code != true {
		log.Println("server start command exited with exit code: ", exit_code)
		return "", fmt.Errorf("server start command exited with exit code: %b", exit_code)
	} else {
		log.Println("server started in the background")
	}
	return c.Name(), nil
}

func CPUCoreSet(cpu []int) string {
	coreset := make([]string, len(cpu))
	for i, v := range cpu {
		coreset[i] = strconv.Itoa(v)
	}
	return strings.Join(coreset, ",")
}

func Memory(memory int) string {
	return fmt.Sprintf("%dTG", memory/1e9)
}

func WaitForIPv4(c *lxc.Container, timeout time.Duration, expected_ip string) (string, error) {
	start := time.Now()
	for time.Since(start) < timeout {
		if c.State() != lxc.RUNNING {
			return "", fmt.Errorf("container %s is not running", c.Name())
		}
		ips, err := c.IPv4Addresses()
		if err == nil && len(ips) > 0 {
			log.Println("detected ips: %v", ips)
			for _, ip := range ips {
				if ip == expected_ip {
					log.Println("successfully verified expected ip: %s", expected_ip)
					return ip, nil
				}
			}
			log.Println("error: detected ips do not match expected: %s", expected_ip)
		}
		time.Sleep(time.Millisecond * 10)
	}
	log.Printf("timeout wating for ip assignment (%v) elapsed\n", timeout)
	return "", nil
}
