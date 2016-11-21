package core

import "sync"
import "fmt"
import "time"

type Injector struct {
	Users          []*Iuser
	wg             sync.WaitGroup
	scenario       string
	ramp           []int64
	rampDuration   int64
	elapsedSeconds int64
	startTime      time.Time
	Stat           *StatStack
}

func NewInjector() *Injector {
	i := &Injector{}
	i.Stat = NewStatStack(i)
	return (i)
}

func (inj *Injector) Initialize(pScenario string) {
	inj.scenario = pScenario
}

// Run start an infinite loop for starting, stopping iusers
func (inj *Injector) Run() {
	inj.startTime = time.Now()
	go inj.Stat.DoRun()

	for shootAgain := bool(true); shootAgain == true; {

		shootAgain = inj.UpdateSpeed()
		time.Sleep(time.Millisecond * 1000)
	}
	fmt.Printf("*** Test done ***\n")
}

// Starts new iusers if needed, according to ramp
func (inj *Injector) UpdateSpeed() bool {
	sumVusers := int64(0)
	cSec := (int64)(time.Now().Sub(inj.startTime) / time.Second) //#sec after start of the test
	if cSec > inj.rampDuration {
		//fmt.Printf("Run lasted: %d and rampDuration %d \n", cSec, inj.rampDuration)
		return false
	}

	for v := 0; v < len(inj.ramp); v += 2 {
		cStepRate := inj.ramp[v]
		cStepDura := inj.ramp[v+1]
		//fmt.Printf("UpdateSpeed v=%d  Rate=%d Duration=%d cSec=%d, sumVusers=%d\n", v, cStepRate, cStepDura, cSec, sumVusers)

		if cSec > cStepDura {
			sumVusers += cStepRate * cStepDura
			cSec -= cStepDura
		} else {
			sumVusers += cStepRate * cSec
			//fmt.Println("We break here")
			break
		}

	}

	delta := sumVusers - int64(len(inj.Users))
	fmt.Printf("Elapsed time(s):%d Computed Requested Vusers: %d, Current Vusers :%d, we need to add:%d\n", cSec, sumVusers, int64(len(inj.Users)), delta)

	for i := int64(0); i < delta; i++ {
		u := NewIuser(inj)
		u.Uuid = fmt.Sprintf("uuid_%05d", len(inj.Users))
		u.LoadScenarioString(inj.scenario)
		u.DoInit()
		go u.DoRun()
		inj.Users = append(inj.Users, u)
	}

	return true

}

func (inj *Injector) SetRamp(pRamp []int64) {
	inj.ramp = pRamp
	for v := 0; v < len(inj.ramp); v += 2 {
		inj.rampDuration += inj.ramp[v+1]
	}

}

func (inj *Injector) GetState() map[int]int {
	m := make(map[int]int)
	for i := range inj.Users {
		u := inj.Users[i]
		m[u.state]++
	}
	print("Injector: GetState: STOPPED:", m[STATE_USER_STOPPED], " RUNNING:", m[STATE_USER_RUNNING], " PAUSED:", m[STATE_USER_PAUSED], "\n")
	return m
}
