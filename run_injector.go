package main

import "./core"

func main() {
	i := new(core.Injector)
	i.NUsers = 10
	i.Run()
}
