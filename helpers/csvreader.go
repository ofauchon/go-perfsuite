package helpers

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
)

type Csvsource struct {
	Records [][]string
}

func (c *Csvsource) LoadFile(f string) {
	fmt.Println("[csvsource] Start reading csv:", f)

	if strings.HasSuffix(f, ".zip") {

	}
	dat, err := os.Open(f)
	if err != nil {
		log.Fatal(err)
	}
	r := csv.NewReader(dat)
	c.Records, err = r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[csvsource] Read %d records", len(c.Records))
}
