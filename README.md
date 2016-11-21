go-perfsuite : Performance Testing Suite in Go



# Introduction

I started this project with no exact idea how far it would go.

Main goals are: 

- Use Go Language for the engine
- Use LUA for scripting 
- Implement common injection features (RampUp, Rendezvous ...) 
- Use influxdb and graphite for metrics storage and rendering

# Quick start


```
go run run_injector.go -scenario myscript.lua -rampup 1,10 
=> Starts test with 'myscript.lua' scenario, and add 1 vuser per second, for 10 second

Another ramp up example :  -rampup 1,60,0,200,10,60,0,3600

1,60  -> Add 1 vuser/s for 60s    ( 1*60 = 60 more users
0,200 -> Add 0 vuser/s for 200s   ( No vuser increment for 200s)
10,60 -> Add 10 vuser/s for 60s   ( 10*60 = 600 more users ) 
0,3600 -> Add 0 vuser/s for 3600s ( No vuser increment for 3600s)
```

# Profiling 


Build and generate cpu and memory profile: 

```
go build run_injector.go
./run_injector -scenario myscenario.lua -rampup 1,180,0,200 -cpuprofile prof/cpu.prof -memprofile prof/mem.prof
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


