package core

type Iscenario interface {
	InitOnce()
	Init()
	Run()
	Stop()
}
