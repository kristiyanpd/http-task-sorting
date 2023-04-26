package main

import (
	"task/pkg/router"
)

func main() {
	// Init and Start Router
	taskRouter := router.Init()
	taskRouter.Start()
}
