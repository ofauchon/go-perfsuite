package core

import (
	"fmt"
	"github.com/yuin/gopher-lua"
	"github.com/cjoudrey/gluahttp"
	"net/http"
	"time"
)

type Counter struct {
	Start int64
	End   int64
}

type Iuser struct {
	Scenario string
	NRuns    int
	Id       int
	Counters map[string]Counter
	LuaState *lua.LState
}

func NewIuser() *Iuser {
	newI := &Iuser{Counters: make(map[string]Counter)}
	Lptr := lua.NewState()
	defer Lptr.Close()

	Lptr.PreloadModule("http", gluahttp.NewHttpModule(&http.Client{}).Loader)
	Lptr.SetGlobal("k_CounterStart", Lptr.NewFunction(newI.k_CounterStart))
	Lptr.SetGlobal("k_CounterEnd", Lptr.NewFunction(newI.k_CounterEnd))
	newI.LuaState = Lptr
	return newI
}

/*
 *  performance wrapped functions
 */
func (i *Iuser) k_CounterStart(L *lua.LState) int {
	tName := L.ToString(1)
	fmt.Printf("DBG: 1 %s\n", tName)

	if _, ok := i.Counters[tName]; !ok {
		tCount := Counter{}
		tCount.Start = time.Now().UnixNano()
		i.Counters[tName] = tCount
	} else {
		fmt.Printf("WARN: Counter '%s' already exisits\n", tName)
	}
	return 1
}

func (i *Iuser) k_CounterEnd(L *lua.LState) int {
	tName := L.ToString(1)

	if xx, ok := i.Counters[tName]; ok {
		xx.End = time.Now().UnixNano()
		fmt.Printf("End counter %s: Delta: %d\n", tName, xx.End-xx.Start)
	} else {
		fmt.Printf("WARN: Counter '%s' can't end while not started\n", tName)
	}
	return 1

}

/*
 *  Entry points
 */

func (i *Iuser) LoadScenarioString(pScenario string) {
	i.Scenario=pScenario
	//fmt.Printf("XXXX %s XXXX\n", i.Scenario)
	if err := i.LuaState.DoString(pScenario); err != nil {
		panic(err)
	}
}

func (i *Iuser) DoInit() {
	if err := i.LuaState.DoString(`rinit()`); err != nil {
		panic(err)
	}
}

func (i *Iuser) DoRun() {
	if err := i.LuaState.DoString(`rrun()`); err != nil {
		panic(err)
	}
}
func (i *Iuser) DoStop() {
	if err := i.LuaState.DoString(`rstop()`); err != nil {
		panic(err)
	}
}
