package models

type CreateFunctionRequest struct {
	FunctionName   string `json:"function_name"`
	FunctionBundle string `json:"function_bundle"`
	ImageName      string `json:"image_name"`
	CPU            []int  `json:"cpu"`
	Memory         int    `json:"memory"`
	VirtualIP      string `json:"virtual_ip"`
	MacAddress     string `json:"mac_address"`
}

type CreateFunctionResponse struct {
	FunctionName string `json:"function_name"`
	ContainerId  string `json:"container_id"`
	ContainerIP  string `json:"container_ip"`
}

type MigrateFunctionRequest struct {
	SourceIP      string `json:"source_ip"`
	DestinationIP string `json:"destination_ip"`
	ContainerId   string `json:"container_id"`
	ImageName     string `json:"image_name"`
}

type StartMigratedFunctionRequest struct {
	ContainerId    string `json:"container_id"`
	CheckPointName string `json:"checkpoint_name"`
}

type StartMigratedFunctionResponse struct {
	Message string `json:"message"`
}

type UpdateFunctionRequest struct {
	ContainerId string `json:"container_id"`
	CorePool    []int  `json:"core_pool"`
	Memory      int    `json:"memory"`
}

type UpdateFunctionResponse struct {
	Message string `json:"message"`
}

type MigrateFunctionResponse struct {
	Message        string `json:"message"`
	CheckPointName string `json:"checkpoint_name"`
}

type InvokeFunctionRequest struct {
	FunctionName string         `json:"function_name"`
	ContainerIP  string         `json:"container_ip"`
	Params       map[string]any `json:"params"`
}

type InvokeFunctionResponse struct {
	Result map[string]any `json:"result"`
	Error  string         `json:"error"`
}

type DeleteFunctionRequest struct {
	ContainerIds []string `json:"container_ids"`
}

type DeleteFunctionResponse struct {
	Result string `json:"result"`
}
