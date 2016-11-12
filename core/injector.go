package core

import "sync"
import "fmt"
import "time"

type Injector struct {
	Users    []*Iuser
	wg       sync.WaitGroup
	Scenario string
	Stat *StatStack
}

func NewInjector() *Injector {
	i := &Injector{}
	i.Stat = NewStatStack(i)
	return (i)
}

func (inj *Injector) Initialize(pUserCount int, pScenario string) {
	inj.Scenario = pScenario
	for k := 0; k < pUserCount; k++ {
		u := NewIuser(inj)
		u.Uuid=fmt.Sprintf("uuid_%05d", k)
		u.LoadScenarioString(inj.Scenario)
		u.DoInit()
		inj.Users = append(inj.Users, u)
	}
}


func (inj *Injector) RunRamp(){

}

func (inj *Injector) Run() {

	go inj.Stat.DoRun() 

	for i := range(inj.Users) {
		u := inj.Users[i]
		//fmt.Println("Starting new user")
		go u.DoRun()
	}
	time.Sleep(time.Second)
	inj.wg.Wait()
	fmt.Println("After wait")
}


func (inj *Injector) GetState() map[int]int{
	m := make(map[int]int)
	for i := range(inj.Users){
		u := inj.Users[i]
		m[u.state]++
	}
	print("Injector: GetState: STOPPED:", m[STATE_USER_STOPPED], " RUNNING:", m[STATE_USER_RUNNING], " PAUSED:", m[STATE_USER_PAUSED],"\n")
	return m
}

