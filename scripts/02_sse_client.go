package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"astuart.co/go-sse"
	"github.com/ofauchon/loadwizard/core"
)

type Scenario struct {
	user *core.Iuser
}

func (s Scenario) Init() {
	/*fmt.Println("Init user " + s.user.Uuid)*/
}

func (s Scenario) Start() {
	for {

		uri := "http://samshull.com/server-sent/events.php"
		event := make(chan *sse.Event, 100)
		go sse.Notify(uri, event)

		for i := 0; i < 10; {
			select {
			case e := <-event:
				res, err := ioutil.ReadAll(e.Data)

				if err != nil {
					fmt.Printf("error reading response body: %s\n", err.Error())
					return
				}

				fmt.Printf("SSE Message:%s\n", res)
			}

			time.Sleep(2 * time.Second)
		}

	}
}

func (s Scenario) Stop() {
	fmt.Println("Stop " + s.user.Uuid)
}

func NewScenario(u *core.Iuser) core.Iscenario {
	r := new(Scenario)
	r.user = u
	return r
}
