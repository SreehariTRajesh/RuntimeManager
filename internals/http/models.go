package http_runtime

type CreateFunctionRequest struct {
	FunctionName   string `json:"function_name"`
	FunctionBundle string `json:"function_bundle"`
	CPU            []int  `json:"cpu"`
	Memory         int    `json:"memory"`
}

type CreateFunctionResponse struct {
	FunctionName string `json:"function_name"`
	ContainerId  string `json:"container_id"`
	ContainerIP  string `json:"container_ip"`
	Error        string `json:"error_message"`
}

type MigrateFunctionRequest struct {
	SourceNode      string `json:"source_node"`
	DestinationNode string `json:"destination_node"`
}

type MigrateFunctionResponse struct {
	Message string `json:"message"`
	Error   bool   `json:"error"`
}

type InvokeFunctionRequest struct {
	ContainerIP string         `json:"container_ip"`
	Params      map[string]any `json:"params"`
}

type InvokeFunctionResponse struct {
	Result map[string]any `json:"result"`
	Error  bool           `json:"error"`
}

type DeleteFunctionRequest struct {
	ContainerIds []string `json:"container_ids"`
}

type DeleteFunctionResponse struct {
	Message string `json:"result"`
	Error   bool   `json:"error"`
}
