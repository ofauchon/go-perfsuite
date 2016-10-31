package core

type Injector struct {
	NUsers int
	Users  []*Iuser
}

func (inj *Injector) Run() {

	for k := 0; k < 10; k++ {
		u := NewIuser()
		tScenario :=`


function rinit()
  print "Lua rinit()"
  local http = require("http")
end


function rrun()
  print "Lua Start rrun()"
  k_CounterStart("test1")

  local http = require("http")
  response, error_message = http.request("GET", "http://google.com", {
    query="q=test",
    headers={
        Accept="*/*"
    }
  })
--  print(response.body)

  k_CounterEnd("test1")
  print "Lua End rrun()"
end 


function rstop()
  print "Lua rstop()"
end
`

        u.LoadScenarioString(tScenario)
	inj.Users = append(inj.Users, u)

	u.DoInit()
	u.DoRun()
	u.DoStop()

	}

}

