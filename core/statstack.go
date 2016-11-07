package core

import "fmt"
import "time"

type StatStack struct {
	Addr 		string
	Db			string
	User		string
	Pass		string
	
	values []Stat
	count	int
}

type Stat struct{
	name	string
	value	float64
}
	

func NewStatStack() *StatStack {
	return &StatStack{ values: make([]Stat, 0) , count: 0 }
}


func (i *StatStack) Push(pName string, pVal float64) {
	i.values=append(i.values,Stat{name: pName, value: pVal})
}

func (i *StatStack) DoRun() {
	fmt.Println("StatStack Started")
	for {
		if (len(i.values)>0){
			fmt.Printf("StatStack Run: %d in stack\n", len(i.values))
		}
		fmt.Printf("StatStack Run: pause\n")
		time.Sleep(300*time.Millisecond)
	}
}


