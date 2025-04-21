package models

type CreateFunctionRequest struct {
	FunctionName   string `json:"function_name"`
	FunctionBundle string `json:"function_bundle"`
	ImageName      string `json:"image_name"`
	CPU            []int  `json:"cpu"`
	Memory         int    `json:"memory"`
	VirtualIP      string `json:"virtual_+ip"`
}

type CreateFunctionResponse struct {
	FunctionName string `json:"function_name"`
	ContainerId  string `json:"container_id"`
	ContainerIP  string `json:"container_ip"`
	Error        string `json:"error_message"`
}

type MigrateFunctionRequest struct {
	SourceIP      string `json:"source_ip"`
	DestinationIP string `json:"destination_ip"`
	ContainerId   string `json:"container_id"`
	ImageName     string `json:"image_name"`
}

type MigrateFunctionResponse struct {
	Message string `json:"message"`
	Error   bool   `json:"error"`
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
	Message string `json:"result"`
	Error   string `json:"error"`
}
