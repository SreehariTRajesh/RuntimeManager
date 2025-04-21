package service

import (
	"fmt"
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
	return nil, nil
}

func MigrateFunction(request *models.MigrateFunctionRequest) (*models.MigrateFunctionResponse, error) {
	return nil, nil
}

func DeleteFunction(request *models.DeleteFunctionRequest) (*models.DeleteFunctionRequest, error) {
	return nil, nil
}
