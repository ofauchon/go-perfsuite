package core

import	"sync"

type Injector struct {
	NUsers 		int
	Users  		[]*Iuser
	wg			sync.WaitGroup
	Scenario	string
}

func NewInjector() *Injector {
	return (&Injector{})
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
	for i:=0; i< inj.NUsers;i++ {
		u := inj.Users[i]
		u.DoRun()
		inj.wg.Wait();
	}
}

