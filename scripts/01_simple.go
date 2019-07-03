package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

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
		/*fmt.Println("Start " + s.user.Uuid)*/
		resp, err := http.Get("http://www.oflabs.com")

		if err != nil {
			fmt.Println("url is down !!!")
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		/*io.Copy(ioutil.Discard, resp.Body)*/

		if err != nil {
			fmt.Println("Can't read all HTTP Response")
		}
		fmt.Printf("HTTP returns Code %d Body Size:%d\n", resp.StatusCode, len(body))

		// Pause 5s
		duration := time.Second * time.Duration(rand.Intn(5))
		time.Sleep(duration)
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
