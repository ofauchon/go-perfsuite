# Loadwizard : Load generator in GO

Loadwizard is an attempt to write a new load generator in Go.

Features:

- Only one language (GO) for both load generator core components and performance scripts 
- Basic injection features at first (constant load, RampUp)
- No onboard metrics and logs processing (everything sent to InfluxDB, Graphite or Elk)

Why Golang ?

First, It's a cool language. But it's also fast and flexible...

Experiments with LUA and JS based load generators convinced me a full-GO solution would be
more efficient.

I choosed the 'GO plugin' approach to isolate loadwizard core and custom performance scripts

# Quick start

```
go get github.com/ofauchon/loadwizard
cd $GOPATH/src/github.com/ofauchon/loadwizard
go build
go install

cd scripts/dataset
gunzip top-1m.csv.gz

cd ../..
go build
./loadwizard -scenario scripts/03_csvsource.go -rampup 1,10,0,100

Rampup : Add 1 user every second for 10 seconds, then add 0 users/sec for 100s

```


# Profiling 


Build and generate cpu and memory profile: 

```
go build
mkdif prof
./loadwizard -scenario scripts/03_csvsource.go -rampup 1,10,0,100 -cpuprofile prof/cpu.prof -memprofile prof/mem.prof
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


 
