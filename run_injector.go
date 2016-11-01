package main

import "./core"
import "flag"
import "os"
import "runtime/pprof"
import "log"
import "io/ioutil"
import "runtime"


func readFile(pPath string) string{
	dat, err := ioutil.ReadFile(pPath)
	if err != nil {
		log.Fatal("Can't open file")
	}
	
    return string(dat)
}

func main() {

	var cpuprofile = flag.String("cpuprofile", "", "write cpu profile `file`")
	var memprofile = flag.String("memprofile", "", "write memory profile to `file`")
	var vusers = flag.Int("vusers", 10, "number of virtual users")
	var scenario = flag.String("scenario", "", "path to scenario")


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


	tScen := readFile(*scenario)
	injector := core.NewInjector()
	injector.Initialize(*vusers, tScen )
	injector.Run()



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
