go-perfsuite : Performance Testing Suite in Go



# Introduction

I started this project with no exact idea how far it would go.

Main goals are: 

- Use Go Language for the engine
- Use LUA for scripting 
- Implement common injection features (RampUp, Rendezvous ...) 
- Use influxdb and graphite for metrics storage and rendering

# Quick start


``
go run run_injector.go  			To run 10 concurent virtual users (default
go run run_injector.go -vuser 20    To run 20 Concurrent vusers
```

# Profiling 


Build and generate cpu and memory profile

``
go build run_injector.go
./run_injector -cpuprofile cpu.prof -memprofile mem.prof
``

``
$ go tool pprof run_injector cpu.prof
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



go tool pprof -alloc_objects run_injector mem.prof
``


