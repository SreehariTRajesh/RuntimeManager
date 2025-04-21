package api

import (
	"github.com/gorilla/mux"
)

func DefineMuxRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/create", CreateFunctionHandler).Methods("POST")
	router.HandleFunc("/invoke", CreateFunctionHandler).Methods("POST")
	router.HandleFunc("/delete", DeleteFunctionHandler).Methods("DELETE")
	router.HandleFunc("/migrate", CreateFunctionHandler).Methods("POST")
	return router
}
