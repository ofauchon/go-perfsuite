package core

import (
	"fmt"
	"time"

	uuid "github.com/nu7hatch/gouuid"
)

const (
	STATE_USER_STOPPED = 0
	STATE_USER_RUNNING = 1
	STATE_USER_PAUSED  = 2
)

type Counter struct {
	Start int64
	End   int64
}

type Iuser struct {
	Uuid     string
	Scenario Iscenario
	NRuns    int
	Id       int
	Inj      *Injector
	Counters map[string]Counter

	state int
}

// Create an IUser, with LUA runtime
func NewIuser(pInj *Injector) *Iuser {
	newI := &Iuser{Counters: make(map[string]Counter)}
	newI.Inj = pInj
	newI.state = STATE_USER_STOPPED
	u4, err := uuid.NewV4()
	if err != nil {
		panic("Can't gen uuid")
	}
	newI.Uuid = u4.String()
	return newI
}

func (i *Iuser) TransactionStart(tName string) {
	if _, ok := i.Counters[tName]; !ok {
		tCount := Counter{}
		tCount.Start = time.Now().UnixNano()
		i.Counters[tName] = tCount
	} else {
		fmt.Printf("WARN: Counter '%s' already exisits\n", tName)
	}
}

func (i *Iuser) TransactionStop(tName string, tStatus int) {
	if xx, ok := i.Counters[tName]; ok {
		xx.End = time.Now().UnixNano()
		tms := (xx.End - xx.Start) / int64(time.Millisecond)
		//fmt.Printf("%s : End counter %s: Delta: %d ms\n", i.Uuid, tName, tms)
		i.Inj.Stat.Push(tName, (float64)(tms))
		delete(i.Counters, tName)
	} else {
		fmt.Printf("WARN: Counter '%s' can't end while not started\n", tName)
	}
}

/*
 *  Entry points
 */
func (i *Iuser) LoadScenarioString(pScenario string) {

}

func (i *Iuser) DoInit() {
	//i.CounterStart(i.Uuid + "_DoInit");
	i.Scenario.Init()
	//i.CounterEnd(i.Uuid + "_DoInit");
}

func (i *Iuser) DoRun() {
	i.Inj.wg.Add(1)
	defer i.Inj.wg.Done()
	i.state = STATE_USER_RUNNING
	//fmt.Println("Iuser DoRun()")
	i.Scenario.Start()
	//i.TransactionStop(i.Uuid+"_DoRun", 1)
	fmt.Println(i.Uuid + "Iuser DoRun() End")
}

func (i *Iuser) DoStop() {
	//i.CounterStart(i.Uuid + "_DoStart");

	i.Scenario.Stop()

	//i.CounterEnd(i.Uuid + "_DoStart");
}
