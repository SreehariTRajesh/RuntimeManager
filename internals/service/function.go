package service

import (
	"fmt"
	"runtime-manager/internals/models"
	"runtime-manager/internals/utils"
)

func CreateFunction(request *models.CreateFunctionRequest) (*models.CreateFunctionResponse, error) {
	fn_name := request.FunctionName
	fn_bundle := request.FunctionBundle
	image_name := request.ImageName
	cpu := request.CPU
	memory := request.Memory
	virt_ip := request.VirtualIP
	mac := request.MacAddress
	cid, err := utils.CreateContainerFunction(fn_name, fn_bundle, image_name, cpu, int64(memory), virt_ip, mac)
	if err != nil {
		return nil, fmt.Errorf("error while creating function response: %w", err)
	}
	return &models.CreateFunctionResponse{
		FunctionName: fn_name,
		ContainerId:  cid,
		ContainerIP:  virt_ip,
	}, nil
}

func InvokeFunction(request *models.InvokeFunctionRequest) (*models.InvokeFunctionResponse, error) {
	return nil, nil
}

func MigrateFunction(request *models.MigrateFunctionRequest) (*models.MigrateFunctionResponse, error) {
	return nil, nil
}

func StartMigratedFunction(request *models.StartMigratedFunctionRequest) (*models.StartMigratedFunctionResponse, error) {
	return nil, nil
}

func DeleteFunction(request *models.DeleteFunctionRequest) (*models.DeleteFunctionResponse, error) {
	id := request.ContainerId
	err := utils.DeleteContainerFunction(id)
	if err != nil {
		return nil, fmt.Errorf("error while deleting container: %w", err)
	}
	return &models.DeleteFunctionResponse{
		Result: "deletion successful",
	}, nil
}

func UpdateResources(request *models.UpdateFunctionRequest) (*models.UpdateFunctionResponse, error) {
	return nil, nil
}
