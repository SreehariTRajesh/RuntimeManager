package main

import (
	"runtime-manager/internals/lifecycle"
	_ "runtime-manager/internals/manager"
)

func main() {
	lifecycle.InitializeAll()
	lifecycle.WaitForShutDown()
}
