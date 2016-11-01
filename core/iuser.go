package core

import (
	"fmt"
	"github.com/yuin/gopher-lua"
	"github.com/cjoudrey/gluahttp"
	"github.com/nu7hatch/gouuid"
	"net/http"
	"time"
)

type Counter struct {
	Start int64
	End   int64
}

type Iuser struct {
	Uuid	 string
	Scenario string
	NRuns    int
	Id       int
	Inj	*Injector
	Counters map[string]Counter
	LuaState *lua.LState
}

func NewIuser(pInj *Injector) *Iuser {
	newI := &Iuser{Counters: make(map[string]Counter), }
	newI.Inj=pInj
	u4, err := uuid.NewV4()
	if err!=nil {
		panic("Can't gen uuid")
	}
	newI.Uuid=u4.String()
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
		fmt.Printf("%s : End counter %s: Delta: %d ms\n", i.Uuid, tName, (xx.End-xx.Start) / int64(time.Millisecond) )
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
//	fmt.Printf("%s : DoInit start\n", i.Uuid)
	if err := i.LuaState.DoString(`rinit()`); err != nil {
		panic(err)
	}
//	fmt.Printf("%s : DoInit end\n", i.Uuid)
}

func (i *Iuser) DoRun() {
//	fmt.Printf("%s : DoRun start\n", i.Uuid)
	i.Inj.wg.Add(1)
	if err := i.LuaState.DoString(`rrun()`); err != nil {
		panic(err)
	}
	i.Inj.wg.Done()
//	fmt.Printf("%s : DoRun end\n", i.Uuid)
}
func (i *Iuser) DoStop() {
//	fmt.Printf("%s : DoStop start\n", i.Uuid)
	if err := i.LuaState.DoString(`rstop()`); err != nil {
		panic(err)
	}
//	fmt.Printf("%s : DoStop end\n", i.Uuid)
}
