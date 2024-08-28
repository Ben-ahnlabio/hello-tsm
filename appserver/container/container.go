package container

import (
	"log"
)

var container *Container

type Container struct {
}

func GetInstnace() *Container {
	if container == nil {
		log.Print("Container is not initialized. Create new container.")

		container = &Container{}
	}
	return container
}
