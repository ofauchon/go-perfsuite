package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/ofauchon/loadwizard/core"
	"github.com/ofauchon/loadwizard/helpers"
)

type Scenario struct {
	user   *core.Iuser
	client http.Client
}

func (s Scenario) InitOnce() {
	fmt.Println("InitOnce Start\n")
	if _, ok := s.user.Inj.Repository["top1mcsv"]; !ok {
		fmt.Println("No CSV available in repository, and will generate it.")
	}
	c := &helpers.Csvsource{}
	c.LoadFile("./scripts/dataset/top-1m.csv")
	s.user.Inj.Repository["top1mcsv"] = c.Records
	fmt.Println("InitOnce End\n")
}

func (s Scenario) Init() {

	s.client = http.Client{
		Timeout: time.Duration(1 * time.Second),
	}

}

func (s Scenario) Run() {

	urls := s.user.Inj.Repository["top1mcsv"].([][]string)
	for {
		site := urls[rand.Intn(len(urls))]
		url := "http://" + site[1] + "/"

		fmt.Println("Start " + url)
		resp, err := s.client.Get(url)

		if err != nil {
			fmt.Println("url is down !!! " + url)
		} else {
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			/*io.Copy(ioutil.Discard, resp.Body)*/

			if err != nil {
				fmt.Println("Can't read all HTTP Response")
			}
			fmt.Printf("HTTP returns Code %d Body Size:%d\n", resp.StatusCode, len(body))
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
