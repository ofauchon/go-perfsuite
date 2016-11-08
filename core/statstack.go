// Package core provide main components of go-perfsuite project
package core

import "fmt"
import "log"
import "time"
import "sync"

import "github.com/influxdata/influxdb/client/v2"

// pause_ms is the delay in millis between two cycles in DoRun() routine
const pause_ms int = 500
const influx_db string = "perfsuite"
const influx_user string = "username"
const influx_pass string = "password"

var mutex = &sync.Mutex{}

type StatStack struct {
	Addr string
	Db   string
	User string
	Pass string

	values []Stat
	count  int
	inj		Injector
}

type Stat struct {
	name  string
	value float64
}

// NewStatStack returns instance of StatStack structure
func NewStatStack(pInj Injector) *StatStack {
	return &StatStack{values: make([]Stat, 0), count: 0, inj: pInj}
}

// Push adds new statistic to the StatStack buffer
// The new statistic is the couple (name,value)
func (i *StatStack) Push(pName string, pVal float64) {
	i.values = append(i.values, Stat{name: pName, value: pVal})
}

func (i *StatStack) FlushInflux() {

	cli, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://localhost:8086",
		Username: influx_user,
		Password: influx_pass,
	})

	if err != nil {
		log.Fatalln("Can't create HTTPClient err:", err)
	}

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  influx_db,
		Precision: "us",
	})

	if err != nil {
		log.Fatalln("Can't create BatchPoints: ", err)
	}

	mutex.Lock(); 
	for z := range i.values {
		stat := i.values[z]

		tags := map[string]string{"metric": stat.name}

		fields := map[string]interface{}{"value": stat.value}

		//fmt.Printf("StatStack : FlushInflux: Add point (%s,%f)\n", stat.name, stat.value)
		point, err := client.NewPoint("stats", tags, fields)
		if err != nil {
			log.Fatalln("Can't create Point: ", err)
		}
		bp.AddPoint(point)
	}
	mutex.Unlock(); 

	err = cli.Write(bp)

	if err != nil {
		log.Fatalln("Can't write batchpoints to InfluxDB: ", err)
	}

}

// DoRun is the main loop that will pop the stats from the buffer and
// send them to the selected backend
func (i *StatStack) DoRun() {
	fmt.Println("StatStack DoRun")
	for {
		if len(i.values) > 0 {
			fmt.Printf("StatStack DoRun: %d in stack\n", len(i.values))
		}
		fmt.Printf("StatStack DoRun: paused for %d ms\n", pause_ms)
		i.FlushInflux()
		time.Sleep(time.Duration(pause_ms) * time.Millisecond)
	}
}
