package api

import (
	"github.com/gorilla/mux"
)

func DefineMuxRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/create", CreateFunctionHandler).Methods("POST")
	router.HandleFunc("/invoke", InvokeFunctionHandler).Methods("POST")
	router.HandleFunc("/delete", DeleteFunctionHandler).Methods("DELETE")
	router.HandleFunc("/migrate", MigrateFunctionHandler).Methods("POST")
	return router
}
