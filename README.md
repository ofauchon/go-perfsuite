##Loadwizard : Load injector for Gophers##

# Introduction

Loadwizard is an attempt to write my own load generator in Go.

Design choices

- Use Golang for the load generator
- Use Golang for performance scripts (go script built built as .so plugins)
- Implement common basic injection features (RampUp, Rendezvous ...) 
- Send metrics and logs (InfluxDB, Graphite, Elk)

Why Golang ?

It's a cool, fast, flexible & powerfull language.
It evident the core injection part of the project would be in go.
After making  experiments with LUA and JS engines, I found out it would be more efficient and fast to write injection scripts in go, too.

I choosed the 'GO plugin' approach to isolate loadwizard core and custom performance scripts

# Quick start

```
go run run_injector.go -scenario scripts/simple.go -rampup 25,100,0,200
```

 -scenario scripts/simple.go    ( Scenario to run  )
 -rampup 25,100,0,200           ( Rampup profile )


# Configuration 

  * Rampup

```
-rampup 2,60,0,200,10,60,0,3600

2,60  -> Add 2 vuser/s for 60s    ( 2*60 = 120 users after 60s)
0,200 -> Add 0 vuser/s for 200s   ( No more vuser  for 200s)
10,60 -> Add 10 vuser/s for 60s   ( 10*60 = 600 more users ) 
0,3600 -> Add 0 vuser/s for 3600s ( No more vuser for 3600s)
```

# Profiling 


Build and generate cpu and memory profile: 

```
go build run_injector.go
./run_injector -scenario myscenario.lua -rampup 1,180,0,200 -cpuprofile prof/cpu.prof -memprofile prof/mem.prof
./run_injector.go -scenario scripts/simple.go -rampup 25,100,0,200 -cpuprofile prof/cpu.prof -memprofile prof/mem.prof
```

Analyze profile : 

```
$ go tool pprof run_injector prof/cpu.prof
(pprof) top
400ms of 520ms total (76.92%)
Showing top 10 nodes out of 222 (cum >= 10ms)
      flat  flat%   sum%        cum   cum%
      90ms 17.31% 17.31%      100ms 19.23%  runtime.scanobject
      50ms  9.62% 26.92%       50ms  9.62%  runtime.mach_semaphore_signal
      50ms  9.62% 36.54%       50ms  9.62%  runtime.memclr
      50ms  9.62% 46.15%       50ms  9.62%  runtime.memmove
      50ms  9.62% 55.77%       50ms  9.62%  runtime.usleep
      40ms  7.69% 63.46%       40ms  7.69%  runtime.cgocall
      30ms  5.77% 69.23%       60ms 11.54%  runtime.(*mcentral).grow
      20ms  3.85% 73.08%       30ms  5.77%  syscall.Syscall
      10ms  1.92% 75.00%       10ms  1.92%  github.com/yuin/gopher-lua.(*LTable).RawGet
      10ms  1.92% 76.92%       10ms  1.92%  nanotime



go tool pprof -alloc_objects run_injector prof/mem.prof
Entering interactive mode (type "help" for commands)
(pprof) top5
1678903 of 4017670 total (41.79%)
Dropped 107 nodes (cum <= 20088)
Showing top 5 nodes out of 150 (cum >= 481762)
      flat  flat%   sum%        cum   cum%
    436910 10.87% 10.87%     436910 10.87%  strings.genSplit
    399781  9.95% 20.83%    1761377 43.84%  encoding/asn1.parseField
    360450  8.97% 29.80%     360450  8.97%  reflect.(*structType).Field
    285152  7.10% 36.89%     285152  7.10%  reflect.unsafe_NewArray
    196610  4.89% 41.79%     481762 11.99%  reflect.MakeSlice
(pprof) 

```


 
