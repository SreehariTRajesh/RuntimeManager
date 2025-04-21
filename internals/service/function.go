package service

import (
	"fmt"
	"log"
	"runtime-manager/internals/models"
	"runtime-manager/internals/pkg"
	"runtime-manager/internals/utils"
)

func CreateFunction(request *models.CreateFunctionRequest) (*models.CreateFunctionResponse, error) {
	function_name := request.FunctionName
	function_bundle_file_path := request.FunctionBundle
	image_name := request.ImageName
	cpu := request.CPU
	memory := request.Memory
	virtual_ip := request.VirtualIP
	container_id, err := utils.CreateAndStartContainer(image_name, cpu, memory, virtual_ip, pkg.MACVLAN_NETWORK_NAME, function_bundle_file_path)
	if err != nil {
		log.Println("error while creating container:", err)
		return &models.CreateFunctionResponse{
			FunctionName: "",
			ContainerId:  "",
			ContainerIP:  "",
			Error:        fmt.Sprintf("error while creating container: %v", err),
		}, fmt.Errorf("error while creating container: %w", err)
	}
	return &models.CreateFunctionResponse{
		FunctionName: function_name,
		ContainerId:  container_id,
		ContainerIP:  virtual_ip,
		Error:        "",
	}, nil
}

func InvokeFunction(request *models.InvokeFunctionRequest) (*models.InvokeFunctionResponse, error) {
	container_ip := request.ContainerIP
	params := request.Params
	function_name := request.FunctionName
	response, err := utils.InvokeFunction(container_ip, params)
	if err != nil {
		return &models.InvokeFunctionResponse{
			Result: nil,
			Error:  fmt.Sprintf("error while invoking function: %s: %v", function_name, err),
		}, fmt.Errorf("error while invoking function: %s", function_name)
	}
	return &models.InvokeFunctionResponse{
		Result: response,
		Error:  "",
	}, nil
}

func MigrateFunction(request *models.MigrateFunctionRequest) (*models.MigrateFunctionResponse, error) {
	container_id := request.ContainerId
	src_ip := request.SourceIP
	dst_ip := request.DestinationIP
	image_name := request.ImageName
	checkpoint_dir := fmt.Sprintf(pkg.DEFAULT_CHECKPOINT_DIR, container_id)
	err := utils.MigrateContainer(src_ip, dst_ip, container_id, checkpoint_dir, image_name)
	if err != nil {
		return &models.MigrateFunctionResponse{
			Message: fmt.Sprintf("error while migrating container %s, %v", container_id, err),
			Error:   true,
		}, fmt.Errorf("error while migrating container %s, %w", container_id, err)
	}
	return &models.MigrateFunctionResponse{
		Message: "Migration successful",
		Error:   false,
	}, nil
}

func DeleteFunction(request *models.DeleteFunctionRequest) (*models.DeleteFunctionRequest, error) {
	return nil, nil
}
