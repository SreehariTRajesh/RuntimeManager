package lifecycle

import (
	"log"
	"os"
	"os/signal"
	"sort"
	"syscall"
)

var cleanables []Cleanable

func RegisterCleanable(cleanable Cleanable) {
	cleanables = append(cleanables, cleanable)
}

func WaitForShutDown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	log.Println("received signal:", sig)
	sort.SliceStable(initializables, func(i, j int) bool {
		return initializables[i].Order() > initializables[j].Order()
	})
	for _, cleanable := range cleanables {
		cleanable.Cleanup()
	}
	log.Println("cleanup complete, exiting")
	os.Exit(0)
}
