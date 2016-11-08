package core

import (
	"fmt"
	"github.com/cjoudrey/gluahttp"
	"github.com/nu7hatch/gouuid"
	"github.com/yuin/gopher-lua"
	"net/http"
	"time"
)

type Counter struct {
	Start int64
	End   int64
}

type Iuser struct {
	Uuid     string
	Scenario string
	NRuns    int
	Id       int
	Inj      *Injector
	Counters map[string]Counter
	LuaState *lua.LState
}

func NewIuser(pInj *Injector) *Iuser {
	newI := &Iuser{Counters: make(map[string]Counter)}
	newI.Inj = pInj
	u4, err := uuid.NewV4()
	if err != nil {
		panic("Can't gen uuid")
	}
	newI.Uuid = u4.String()
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
func (i *Iuser) CounterStart(tName string) {
	if _, ok := i.Counters[tName]; !ok {
		tCount := Counter{}
		tCount.Start = time.Now().UnixNano()
		i.Counters[tName] = tCount
	} else {
		fmt.Printf("WARN: Counter '%s' already exisits\n", tName)
	}
}
func (i *Iuser) k_CounterStart(L *lua.LState) int {
	tName := L.ToString(1)
	i.CounterStart(tName)
	return 1
}

func (i *Iuser) CounterEnd(tName string) {
	if xx, ok := i.Counters[tName]; ok {
		xx.End = time.Now().UnixNano()
		tms := (xx.End - xx.Start) / int64(time.Millisecond)
		fmt.Printf("%s : End counter %s: Delta: %d ms\n", i.Uuid, tName, tms)
		i.Inj.Stat.Push(tName, (float64)(tms))
	} else {
		fmt.Printf("WARN: Counter '%s' can't end while not started\n", tName)
	}
}
func (i *Iuser) k_CounterEnd(L *lua.LState) int {
	tName := L.ToString(1)
	i.CounterEnd(tName)
	return 1
}

/*
 *  Entry points
 */
func (i *Iuser) LoadScenarioString(pScenario string) {
	i.Scenario = pScenario
	if err := i.LuaState.DoString(pScenario); err != nil {
		panic(err)
	}
}

func (i *Iuser) DoInit() {
	//i.CounterStart(i.Uuid + "_DoInit");
	if err := i.LuaState.DoString(`rinit()`); err != nil {
		panic(err)
	}
	//i.CounterEnd(i.Uuid + "_DoInit");
}

func (i *Iuser) DoRun() {
	i.Inj.wg.Add(1)
	defer i.Inj.wg.Done()
	fmt.Println("Iuser DoRun()")
	i.CounterStart(i.Uuid + "_DoRun")
	if err := i.LuaState.DoString(`rrun()`); err != nil {
		panic(err)
	}
	i.CounterEnd(i.Uuid + "_DoRun")
	fmt.Println("Iuser DoRun() End")
}
func (i *Iuser) DoStop() {
	//i.CounterStart(i.Uuid + "_DoStart");
	if err := i.LuaState.DoString(`rstop()`); err != nil {
		panic(err)
	}
	//i.CounterEnd(i.Uuid + "_DoStart");
}
