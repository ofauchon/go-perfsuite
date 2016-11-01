package core

import	"sync"

type Injector struct {
	NUsers int
	Users  []*Iuser
	wg	sync.WaitGroup
}

func (inj *Injector) Run() {

	for k := 0; k < inj.NUsers; k++ {
		u := NewIuser(inj)


// ------ START LUA
	tScenario :=`


function rinit()
  local http = require("http")
end


function rrun()
  k_CounterStart("test1")

  local http = require("http")
  response, error_message = http.request("GET", "http://www.oflabs.com", {
    query="q=test",
    headers={
        Accept="*/*"
    }
  })
--  print(response.body)

  k_CounterEnd("test1")
end 


function rstop()
end
`
// ------ END LUA

        u.LoadScenarioString(tScenario)
	inj.Users = append(inj.Users, u)

	u.DoInit()
	go u.DoRun()
	//u.DoStop()

	inj.wg.Wait();

	}

}

