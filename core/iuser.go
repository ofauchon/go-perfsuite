package core

import (
	"fmt"
	"encoding/base64"
	"github.com/yuin/gluare"
	"../../gluahttp"
	"github.com/nu7hatch/gouuid"
	"github.com/yuin/gopher-lua"
	"net/http"
	"time"
)


const (
	STATE_USER_STOPPED=0
	STATE_USER_RUNNING=1
	STATE_USER_PAUSED=2
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

	state    int
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
	Lptr := lua.NewState()
	defer Lptr.Close()

	Lptr.PreloadModule("http", gluahttp.NewHttpModule(&http.Client{}).Loader)
	Lptr.PreloadModule("re", gluare.Loader)
	
	Lptr.SetGlobal("k_TransactionStart", Lptr.NewFunction(newI.k_TransactionStart))
	Lptr.SetGlobal("k_TransactionStop", Lptr.NewFunction(newI.k_TransactionStop))
	Lptr.SetGlobal("k_Sleep", Lptr.NewFunction(newI.k_Sleep))
	Lptr.SetGlobal("k_GetId", Lptr.NewFunction(newI.k_GetId))
	Lptr.SetGlobal("k_b64enc", Lptr.NewFunction(newI.k_b64enc))
	newI.LuaState = Lptr
	return newI
}

/*
 *  performance wrapped functions
 */
func (i *Iuser) k_Sleep(L *lua.LState) int{
	i.state=STATE_USER_PAUSED
	ts := (time.Duration)(L.ToInt(1))
	time.Sleep(ts * time.Millisecond)
	i.state=STATE_USER_RUNNING	
	return 1
}

func (i *Iuser) k_GetId(L *lua.LState) int{
	tid:=lua.LString(i.Uuid)
	L.Push(tid)
	return 1
}
func (i *Iuser) k_b64enc(L *lua.LState) int{
	tstr:=L.ToString(1)
	rb64:=base64.StdEncoding.EncodeToString([]byte(tstr))
	L.Push(lua.LString(rb64))
	return 1
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
func (i *Iuser) k_TransactionStart(L *lua.LState) int {
	tName := L.ToString(1)
	i.TransactionStart(tName)
	return 1
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

func (i *Iuser) k_TransactionStop(L *lua.LState) int {
	tName := L.ToString(1)
	tStatus := L.ToInt(2)
	i.TransactionStop(tName,tStatus)
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
	i.state = STATE_USER_RUNNING
	//fmt.Println("Iuser DoRun()")
	i.TransactionStart(i.Uuid + "_DoRun")
	if err := i.LuaState.DoString(`rrun()`); err != nil {
		panic(err)
	}
	i.TransactionStop(i.Uuid + "_DoRun", 1)
	fmt.Println(i.Uuid + "Iuser DoRun() End")
}

func (i *Iuser) DoStop() {
	//i.CounterStart(i.Uuid + "_DoStart");
	if err := i.LuaState.DoString(`rstop()`); err != nil {
		panic(err)
	}
	//i.CounterEnd(i.Uuid + "_DoStart");
}
