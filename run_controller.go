package main

import "core"

func main() {

	ctrl := core.NewControllerInstance()
	ctrl.StartHttpServer()

}
