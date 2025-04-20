package main

import (
	"runtime-manager/internals/lifecycle"
	_ "runtime-manager/internals/macvlan"
	_ "runtime-manager/internals/vxlan"
)

func main() {
	lifecycle.InitializeAll()
	lifecycle.WaitForShutDown()
}
