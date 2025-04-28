package service

import (
	"fmt"
	"log"
	"runtime-manager/internals/manager"
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
	mac_address := request.MacAddress
	container_id, err := utils.CreateAndStartContainer(image_name, cpu, memory, virtual_ip, pkg.MACVLAN_NETWORK_NAME, mac_address, function_bundle_file_path)
	if err != nil {
		log.Println("error while creating container:", err)
		return nil, fmt.Errorf("error while creating container: %w", err)
	}
	// add the Container Id to the Container Registry
	manager.ContainerRegistry.Add(container_id)
	return &models.CreateFunctionResponse{
		FunctionName: function_name,
		ContainerId:  container_id,
		ContainerIP:  virtual_ip,
	}, nil
}

func InvokeFunction(request *models.InvokeFunctionRequest) (*models.InvokeFunctionResponse, error) {
	container_ip := request.ContainerIP
	params := request.Params
	function_name := request.FunctionName
	response, err := utils.InvokeFunction(container_ip, params)
	if err != nil {
		return nil, fmt.Errorf("error while invoking function: %s: %w", function_name, err)
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
	cp, err := utils.MigrateContainer(src_ip, dst_ip, container_id, image_name)
	if err != nil {
		return nil, fmt.Errorf("error while migrating container %s, %w", container_id, err)
	}
	return &models.MigrateFunctionResponse{
		Message:        "migration successful",
		CheckPointName: cp,
	}, nil
}

func StartMigratedFunction(request *models.StartMigratedFunctionRequest) (*models.StartMigratedFunctionResponse, error) {
	checkpoint_id := request.CheckpointId
	function_bundle := request.FunctionBundle
	image := request.ImageName
	cpu := request.CPU
	memory := request.Memory
	virtual_ip := request.VirtualIP
	mac_address := request.MacAddress
	err := utils.StartMigratedContainer(image, cpu, memory, virtual_ip, pkg.MACVLAN_NETWORK_NAME, mac_address, function_bundle, checkpoint_id)
	if err != nil {
		return nil, fmt.Errorf("error starting the migrated function: %w", err)
	}
	return &models.StartMigratedFunctionResponse{
		Message: "migrated function started successfully",
	}, nil
}

func DeleteFunction(request *models.DeleteFunctionRequest) (*models.DeleteFunctionResponse, error) {
	container_ids := request.ContainerIds
	//
	for _, id := range container_ids {
		utils.DeleteContainer(id)
	}
	return &models.DeleteFunctionResponse{
		Result: fmt.Sprintf("successfully deleted containers: %v", container_ids),
	}, nil
}

func UpdateResources(request *models.UpdateFunctionRequest) (*models.UpdateFunctionResponse, error) {
	container_id := request.ContainerId
	core_pool := request.CorePool
	memory := request.Memory
	err := utils.UpdateCorePool(core_pool, container_id)
	if err != nil {
		return nil, fmt.Errorf("error while updating corepool: %v", err)
	}
	err = utils.UpdateMemory(int64(memory), container_id)
	if err != nil {
		return nil, fmt.Errorf("error while updating memory: %v", err)
	}
	return &models.UpdateFunctionResponse{
		Message: "resources updated sucessfully",
	}, nil
}
