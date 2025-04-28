package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime-manager/internals/pkg"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
)

func CreateAndStartContainer(image_name string, cores []int, memory int, virtual_ip string, network_name string, static_mac string, function_bundle string) (string, error) {
	// functionality to start a container
	/*
			@param: cpu_cores
			@param: memory
			@param: virtual_ip
			@param: function_bundle
			# default port assigned 80if err != nil {
			return nil, fmt.Errorf("error creating docker client: %w", err)
		}
		defer cli.Close()
	*/
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return "", fmt.Errorf("error creating docker client: %w", err)
	}
	// Get the network name from the network_manager
	//network_name := macvlan.GetNetworkname()
	cpuSet := getCoreSet(cores)

	resource_config := container.Resources{
		Memory:     int64(memory),
		CpusetCpus: cpuSet,
	}

	network_config := &network.EndpointSettings{
		IPAMConfig: &network.EndpointIPAMConfig{
			IPv4Address: virtual_ip,
		},
		MacAddress: static_mac,
	}

	exposed_ports := nat.PortSet{
		nat.Port(pkg.DEFAULT_CONTAINER_PORT): struct{}{},
	}

	networking_config := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			network_name: network_config,
		},
	}

	host_config := &container.HostConfig{
		Resources: resource_config,
		Binds: []string{
			fmt.Sprintf("%s:/app", function_bundle),
		},
	}

	container_config := &container.Config{
		Image:        image_name,
		Cmd:          []string{"sh", "-c", "chmod +x /app/startup.sh && /app/startup.sh"},
		ExposedPorts: exposed_ports,
	}

	// create the container
	res, err := cli.ContainerCreate(ctx, container_config, host_config, networking_config, nil, "")

	if err != nil {
		return "", fmt.Errorf("error while creating container: %w", err)
	}

	//start the container
	if err := cli.ContainerStart(ctx, res.ID, container.StartOptions{}); err != nil {
		return "", fmt.Errorf("error starting the container: %w", err)
	}
	return res.ID, nil
}

func InvokeFunction(container_ip string, params map[string]any) (map[string]any, error) {
	// functionality to execute container functions
	return MakeHttpRequest(container_ip, 80, params)
}

// set of util functions;
func getCoreSet(cores []int) string {
	core_array := make([]string, len(cores))
	for i, core := range cores {
		core_array[i] = strconv.Itoa(core)
	}
	return strings.Join(core_array, ",")
}

// assign a name for the container -> name must be unique
func AssignName() string {
	unique_id := uuid.New().String()
	return unique_id[:16]
}

func MakeHttpRequest(ip string, port int, params map[string]any) (map[string]any, error) {
	host := fmt.Sprintf("%s:%d", ip, port)
	path := "/invoke"
	scheme := "http"
	baseURL := &url.URL{
		Scheme: scheme,
		Host:   host,
		Path:   path,
	}
	http_url := baseURL.String()
	request_body, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("error while marshalling the request body %v: %w", params, err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", http_url, bytes.NewBuffer(request_body))

	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("error making http request: %w", err)
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, fmt.Errorf("error while reading response: %w", err)
	}

	var response map[string]any
	err = json.Unmarshal(body, &response)

	if err != nil {
		return nil, fmt.Errorf("error while umarshalling json response: %w", err)
	}

	return response, nil
}

func MigrateContainer(source_ip string, dest_ip string, container_id string, image_name string) (string, error) {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv)

	if err != nil {
		return "", fmt.Errorf("error creating a docker client: %w", err)
	}

	checkPointOpts := types.CheckpointCreateOptions{
		CheckpointID: fmt.Sprintf("cp-%s", container_id),
		Exit:         true,
	}

	if err := cli.CheckpointCreate(ctx, container_id, checkPointOpts); err != nil {
		return "", fmt.Errorf("error creating checkpoint for container %s: %w", container_id, err)
	}
	//
	log.Println("checkpoint created for container: ", container_id)
	// checkpoint path

	err = TransferCheckpointFiles(container_id, pkg.DEFAULT_CHECKPOINT_DIR_PARENT, dest_ip)

	if err != nil {
		return "", fmt.Errorf("error while transferring files to remote host: %w", err)
	}

	log.Println("config files transferred to destination node:", dest_ip)

	// remove container;
	if err := cli.ContainerRemove(ctx, container_id, container.RemoveOptions{}); err != nil {
		return "", fmt.Errorf("error removing the container: %w", err)
	}

	return fmt.Sprintf("cp-%s", container_id), nil
}

func TransferContainerFiles(container_id string, dst_checkpoint_dir string, dst_ip string) error {
	src_checkpoint_dir := fmt.Sprintf(pkg.DEFAULT_DOCKER_CONTAINER_PATH, container_id)
	cmd := fmt.Sprintf("scp -r %s rajesh@%s:%s", src_checkpoint_dir, dst_ip, dst_checkpoint_dir)
	fmt.Println("Executing command:", cmd)
	if err := ExecCommand(cmd); err != nil {
		return fmt.Errorf("error transferring checkpoint files: %w", err)
	}
	return nil
}

func TransferCheckpointFiles(src_checkpoint_dir string, dst_checkpoint_dir string, dest_ip string) error {
	cmd := fmt.Sprintf("scp -r %s rajesh@%s:%s", src_checkpoint_dir, dest_ip, dst_checkpoint_dir)
	fmt.Println("Executing command:", cmd)
	if err := ExecCommand(cmd); err != nil {
		return fmt.Errorf("error transferring checkpoint files: %w", err)
	}
	return nil
}

func TransferConfigFiles(src_config_file string, dst_config_dir string, dest_ip string) error {
	cmd := fmt.Sprintf("scp -r %s rajesh@%s:%s", src_config_file, dest_ip, dst_config_dir)
	fmt.Println("Executing command:", cmd)
	if err := ExecCommand(cmd); err != nil {
		return fmt.Errorf("error transferring config files: %w", err)
	}
	return nil
}

func ExecCommand(command string) error {
	out, err := exec.Command("bash", "-c", command).CombinedOutput()
	if err != nil {
		return fmt.Errorf("command execution failed: %s, output: %s", err, string(out))
	}
	log.Println("command output:", string(out))
	return nil
}

func CopyCheckpointToDockerDir(container_id string, checkpoint_dir string, checkpoint_name string) error {
	checkpoint_dst := fmt.Sprintf(pkg.DEFAULT_DOCKER_CHECKPOINT_PATH, container_id, checkpoint_name)
	err := CopyDirectory(checkpoint_dir, checkpoint_dst)
	if err != nil {
		return fmt.Errorf("error copying checkpoint directory: %w", err)
	}
	return nil
}

func CopyContainerFilesToDockerDir(container_id string) error {
	container_dir := pkg.DEFAULT_DOCKER_CONTAINER_PARENT
	err := CopyDirectory(container_dir, fmt.Sprintf(pkg.DEFAULT_CHECKPOINT_DIR, container_id))
	if err != nil {
		return fmt.Errorf("error copying the container files to docker directory: %w", err)
	}
	return nil
}

func CopyDirectory(src_dir string, dst_dir string) error {
	log.Printf("copying from %s to %s", src_dir, dst_dir)
	return filepath.WalkDir(src_dir, func(path string, dir fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error traversing path %s: %w", src_dir, err)
		}
		rel_path, err := filepath.Rel(src_dir, path)
		if err != nil {
			return fmt.Errorf("error while extracting relative path: %w", err)
		}
		dst_path := filepath.Join(dst_dir, rel_path)

		if dir.IsDir() {
			return os.MkdirAll(dst_path, 0755)
		} else {
			return CopyFile(path, dst_path)
		}
	})
}

func CopyFile(src_file string, dst_file string) error {
	src, err := os.Open(src_file)
	if err != nil {
		return fmt.Errorf("error while opening file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(dst_file)
	if err != nil {
		return fmt.Errorf("error while creating the file: %w", err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return fmt.Errorf("error while copying file from src to dst: %w", err)
	}
	return dst.Sync()
}

func StartMigratedContainer(container_id string, checkpoint_id string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return fmt.Errorf("error creating docker client for remote host: %w", err)
	}
	// check if image exists
	err = CopyContainerFilesToDockerDir(container_id)
	if err != nil {
		return fmt.Errorf("error copying checkpoint to docker directory: %w", err)
	}

	fmt.Printf("container with id: %s\n", container_id)

	// once checkpoint is copied, start the container from the host machine hosting t
	// transfer check point files from

	log.Printf("container %s created successfully", container_id)

	if err := cli.ContainerStart(ctx, container_id, container.StartOptions{
		CheckpointID: checkpoint_id,
	}); err != nil {
		return fmt.Errorf("error starting the container on remote host: %w", err)
	}
	return nil
}

func GetNetworkConfigsForMigratedContainer(network_name string, settings map[string]*network.EndpointSettings, mac_address string) *network.NetworkingConfig {
	virtual_ip := settings[network_name].IPAMConfig.IPv4Address
	fmt.Println("MAC ADDRESS:", mac_address)
	network_config := &network.EndpointSettings{
		IPAMConfig: &network.EndpointIPAMConfig{
			IPv4Address: virtual_ip,
		},
		MacAddress: mac_address,
	}

	networking_config := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			network_name: network_config,
		},
	}

	return networking_config
}

func DeleteContainer(container_id string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return fmt.Errorf("error creating docker client: %w", err)
	}
	// delete container
	if err := cli.ContainerStop(ctx, container_id, container.StopOptions{}); err != nil {
		return fmt.Errorf("error stopping the container: %w", err)
	}

	if err := cli.ContainerRemove(ctx, container_id, container.RemoveOptions{}); err != nil {
		return fmt.Errorf("error removing the container: %w", err)
	}
	return nil
}

func GetContainerConfigs(container_id string) (*types.ContainerJSON, error) {
	file_path := fmt.Sprintf("%s/config.v2.json", pkg.DEFAULT_CHECKPOINT_DIR_PARENT)
	config_data, err := os.ReadFile(file_path)
	if err != nil {
		return nil, fmt.Errorf("error reading config data from file: %w", err)
	}
	var container_config types.ContainerJSON
	err = json.Unmarshal(config_data, &container_config)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling data to types.ContainerJSON: %w", err)
	}
	fmt.Println(container_config)
	return &container_config, nil
}

/*
func MigrateContainerV2(src_ip, dest_ip string, container_id string) error {
	ctx := context.Background()
	source_docker_host := fmt.Sprintf("tcp://%s:2375", src_ip)
	destination_docker_host := fmt.Sprintf("tcp://%s:2375", dest_ip)
	checkpoint_name := uuid.New().String()[:8]
	source_cli, err := client.NewClientWithOpts(client.WithHost(source_docker_host), client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("error creating docker client for source: %s: %w", source_docker_host, err)
	}
	defer source_cli.Close()
	destination_cli, err := client.NewClientWithOpts(client.WithHost(destination_docker_host), client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("error creating docker client for source: %s: %w", source_docker_host, err)
	}
	defer destination_cli.Close()
	err = source_cli.CheckpointCreate(ctx, container_id, checkpoint.CreateOptions{
		CheckpointID: checkpoint_name,
	})
	if err != nil {
		return fmt.Errorf("failed to create checkpoint: %s for container %s: %v", checkpoint_name)
	}
	// transfer checkpoint files
}
*/
