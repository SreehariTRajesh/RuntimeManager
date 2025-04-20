package lifecycle

import (
	"log"
	"sort"
)

var initializables []Initializable

func RegisterInitializable(initializable Initializable) {
	initializables = append(initializables, initializable)
}

func InitializeAll() {
	sort.SliceStable(initializables, func(i, j int) bool {
		return initializables[i].Order() < initializables[j].Order()
	})
	for _, initializable := range initializables {
		err := initializable.Initialize()
		if err != nil {
			log.Println(err)
		}
	}
}
