package http_runtime

type CreateFunctionRequest struct {
	FunctionName   string `json:"function_name"`
	FunctionBundle string `json:"function_bundle"`
	CPU            []int  `json:"cpu"`
	Memory         int    `json:"memory"`
}

type CreateFunctionResponse struct {
	
}

type MigrateFunctionRequest struct {
}

type MigrateFunctionResponse struct {
}

type InvokeFunctionRequest struct {
}

type InvokeFunctionResponse struct {
}

type DeleteFunctionRequest struct {
}

type DeleteFunctionResponse struct {
}
