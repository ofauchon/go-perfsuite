package core

import "sync"
import "fmt"
import "time"

type Injector struct {
	Users 			[]*Iuser
	wg       		sync.WaitGroup
	scenario 		string
	ramp 			[]int64
	elapsedSeconds 	int64
	startTime		time.Time
	Stat 			*StatStack
}

func NewInjector() *Injector {
	i := &Injector{}
	i.Stat = NewStatStack(i)
	return (i)
}

func (inj *Injector) Initialize(pScenario string) {
	inj.scenario = pScenario
	/*
	for k := 0; k < pUserCount; k++ {
		u := NewIuser(inj)
		u.Uuid=fmt.Sprintf("uuid_%05d", k)
		u.LoadScenarioString(inj.Scenario)
		u.DoInit()
		inj.Users = append(inj.Users, u)
	}
	*/
}


// Run start an infinite loop for starting, stopping iusers
func (inj *Injector) Run() {
	inj.startTime=time.Now()
	go inj.Stat.DoRun() 

	for {
		inj.UpdateSpeed();
		time.Sleep(time.Millisecond * 5000 )
	}
}

// Starts new iusers if needed, according to ramp
func (inj *Injector) UpdateSpeed() {
	sumVusers:=int64(0)
	cSec := (int64)(time.Now().Sub(inj.startTime)/time.Second)



	for v:=0; v < len(inj.ramp); v+=2 {
		cStepRate:=inj.ramp[v]
		cStepDura:=inj.ramp[v+1]
		//fmt.Printf("UpdateSpeed v=%d  Rate=%d Duration=%d cSec=%d, sumVusers=%d\n", v, cStepRate, cStepDura, cSec, sumVusers)
		if (cSec>cStepDura){
			sumVusers+=cStepRate*cStepDura
			cSec -= cStepDura
		} else {
			sumVusers+=cStepRate*cSec 
			//fmt.Println("We break here")
			break
		}

	}

	delta := sumVusers - int64(len(inj.Users))
	fmt.Printf("Elapsed time(s):%d Computed Requested Vusers: %d, Current Vusers :%d, we need to add:%d\n", cSec, sumVusers, int64(len(inj.Users)), delta)

	for i:=int64(0);i<delta;i++ {
		u := NewIuser(inj)
		u.Uuid=fmt.Sprintf("uuid_%05d", len(inj.Users))
		u.LoadScenarioString(inj.scenario)
		u.DoInit()
		go u.DoRun();
		inj.Users = append(inj.Users, u)
	}

}

func (inj *Injector) SetRamp( pRamp []int64) {
	inj.ramp=pRamp
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

