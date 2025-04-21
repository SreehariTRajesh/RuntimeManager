package api

import (
	"encoding/json"
	"net/http"
	"runtime-manager/internals/models"
)

func CreateFunctionHandler(res http.ResponseWriter, req *http.Request) {
	var create_function_request models.CreateFunctionRequest
	if err := json.NewDecoder(req.Body).Decode(&create_function_request); err != nil {
		http.Error(res, "invalid request body", http.StatusBadRequest)
	}

}

func InvokeFunctionHandler(res http.ResponseWriter, req *http.Request) {

}

func DeleteFunctionHandler(res http.ResponseWriter, req *http.Request) {

}

func MigrateFunctionHandler(res http.ResponseWriter, req *http.Request) {

}
