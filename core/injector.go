package core

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"plugin"
	"strings"
	"sync"
	"time"
)

type Injector struct {
	Users          []*Iuser
	wg             sync.WaitGroup
	plugin         *plugin.Plugin
	ramp           []int64
	rampDuration   int64
	elapsedSeconds int64
	startTime      time.Time
	Stat           *StatStack
}

func NewInjector() *Injector {
	i := &Injector{}
	//i.Stat = NewStatStack(i)
	return (i)
}

func (inj *Injector) Initialize(scriptFile string) {

	if !strings.HasSuffix(strings.ToLower(scriptFile), ".go") {
		fmt.Printf("Script file (%s) must end with .go extension\n", scriptFile)
		os.Exit(1)
	}

	soFile := scriptFile[0:len(scriptFile)-3] + ".so"
	/*
		_, err := os.Stat(soFile)
		fmt.Println("No " + soFile + ", trying to build it ")
	*/
	fmt.Println("Building " + scriptFile)
	cmd := exec.Command(`go`, `build`, `-o`, soFile, `-buildmode=plugin`, scriptFile)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Loading script plugin " + soFile)
	inj.plugin, err = plugin.Open(soFile)
	if err != nil {
		fmt.Printf("Can't load associated .so file (%s), did you build it ?\n", soFile)
		fmt.Println(err)
		os.Exit(1)
	}

}

// Run start an infinite loop for starting, stopping iusers
func (inj *Injector) Run() {
	inj.startTime = time.Now()
	//go inj.Stat.DoRun()

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

		// Instanciate plugin instance for every VU
		s1, err := inj.plugin.Lookup("NewScenario")
		if err != nil {
			fmt.Printf("Damned, no 'scenario' in .so script")
			fmt.Println(err)
			os.Exit(1)
		}
		s2, ok := s1.(func(*Iuser) Iscenario)
		//fmt.Printf("result: %T %v %v\n", s2, s2, ok)

		if !ok {
			fmt.Println("unexpected type from module symbol")
			fmt.Println(err)
			os.Exit(1)
		}
		u.Scenario = s2(u)

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
