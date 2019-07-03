package main

import "github.com/ofauchon/loadwizard/core"
import "flag"
import "os"
import "runtime/pprof"
import "log"

import "runtime"
import "strings"
import "strconv"

func parseRampUp(s string) ([]int64, error) {
	ramp := strings.Split(s, ",")
	rampi := make([]int64, len(ramp))
	for i := 0; i < len(ramp); i++ {
		var err error
		rampi[i], err = strconv.ParseInt(ramp[i], 10, 64)
		if err != nil {
			return []int64{}, err
		}
	}
	return rampi, nil
}

func main() {

	// Command line flags
	var cpuprofile = flag.String("cpuprofile", "", "write cpu profile `file`")
	var memprofile = flag.String("memprofile", "", "write memory profile to `file`")
	var rampup = flag.String("rampup", "", "Ramp-up profile : ex 2,10,0,180,1,20,0,360")
	var scenarioFile = flag.String("scenario", "", "path to go scenario (.go file)")

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	var ramp []int64
	var err error
	if *rampup != "" {
		ramp, err = parseRampUp(*rampup)
		if err != nil {
			log.Fatal("Error parsing rampup parameter '%s'", rampup)
		}
	}

	// Our code
	injector := core.NewInjector()
	injector.Initialize(*scenarioFile)
	injector.SetRamp(ramp)

	injector.Run()

	// Profiling
	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
		f.Close()
	}

}
