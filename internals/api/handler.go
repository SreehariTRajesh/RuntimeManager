package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime-manager/internals/models"
	"runtime-manager/internals/service"
)

func CreateFunctionHandler(res http.ResponseWriter, req *http.Request) {
	var create_function_request models.CreateFunctionRequest
	if err := json.NewDecoder(req.Body).Decode(&create_function_request); err != nil {
		http.Error(res, "invalid request body", http.StatusBadRequest)
		return
	}

	response, err := service.CreateFunction(&create_function_request)
	if err != nil {
		error_message := fmt.Sprintf("failed to create function: %v", err)
		http.Error(res, error_message, http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(response)
}

func InvokeFunctionHandler(res http.ResponseWriter, req *http.Request) {
	var invoke_function_request models.InvokeFunctionRequest
	if err := json.NewDecoder(req.Body).Decode(&invoke_function_request); err != nil {
		http.Error(res, "invalid request body", http.StatusBadRequest)
		return
	}
	response, err := service.InvokeFunction(&invoke_function_request)
	if err != nil {
		error_message := fmt.Sprintf("failed to invoke function: %v", err)
		http.Error(res, error_message, http.StatusInternalServerError)
		return
	}
	// response encoding for function invocation errors
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(response)
}

func DeleteFunctionHandler(res http.ResponseWriter, req *http.Request) {
	var delete_function_request models.DeleteFunctionRequest
	if err := json.NewDecoder(req.Body).Decode(&delete_function_request); err != nil {
		http.Error(res, "invalid request body", http.StatusBadRequest)
		return
	}
	response, err := service.DeleteFunction(&delete_function_request)
	if err != nil {
		http.Error(res, "failed to delete function", http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(response)
}

func MigrateFunctionHandler(res http.ResponseWriter, req *http.Request) {
	var migrate_function_request models.MigrateFunctionRequest
	if err := json.NewDecoder(req.Body).Decode(&migrate_function_request); err != nil {
		http.Error(res, "invalid request body", http.StatusBadRequest)
		return
	}
	response, err := service.MigrateFunction(&migrate_function_request)
	if err != nil {
		http.Error(res, "failed to migrate function", http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(response)
}
