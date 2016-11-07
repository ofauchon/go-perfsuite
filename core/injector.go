package core

import	"sync"
import	"fmt"

type Injector struct {
	NUsers 		int
	Users  		[]*Iuser
	wg			sync.WaitGroup
	Scenario	string

	Stat 		*StatStack
}

func NewInjector() *Injector {
	i:=Injector{}
	i.Stat= NewStatStack()
	return (&i)
}

func (inj *Injector) Initialize(pUserCount int, pScenario string) {
	inj.NUsers=pUserCount
	inj.Scenario=pScenario
	for k := 0; k < inj.NUsers; k++ {
		u := NewIuser(inj)
		u.LoadScenarioString(inj.Scenario)
		u.DoInit()
		inj.Users = append(inj.Users, u)
	}
	//fmt.Printf("Injector : #users:%d Scenario:%s\n", inj.NUsers, inj.Scenario)
}

func (inj *Injector) Run() {
	fmt.Println("Injector Run()")

	go inj.Stat.DoRun()

	for i:=0; i< inj.NUsers;i++ {
		u := inj.Users[i]
		go u.DoRun()
	}
	inj.wg.Wait();
}

