package main

import (
	"log"
	"net/http"
	"runtime-manager/internals/api"
	"runtime-manager/internals/lifecycle"
	_ "runtime-manager/internals/manager"
)

func main() {
	go func() {
		lifecycle.InitializeAll()
		lifecycle.WaitForShutDown()
	}()
	router := api.DefineMuxRouter()
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("error starting http server: %v", err)
	}
}
