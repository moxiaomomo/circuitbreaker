# circuitbreaker

This is a circuitbreaker module which is easy to use when you build your microservice.

## How to use it

```golang
// func NewCirucuitBreaker(timeWin time.Duration, failCnt int, failPercent int) *Circuits

// create an instance of circuitbreaker
cbs := NewCirucuitBreaker(time.Second, 150, 20)
// register a command to a circuitbreaker
suc := cbs.RegisterCommandAsDefault(testcmd)

// report a result (true or false) to a circuit with the command-name
cbs.Report(testcmd, false)
cbs.Report(testcmd, true)

// to check if current request could be allowed
execAllow := cbs.AllowExec(testcmd)
```

## Testing

```shell
$ go test
[DEBUG] To run test for command: testcmd
[DEBUG] [spec]total: 1000, allow: 1000
[DEBUG] Closed --> Open
[DEBUG] Open --> HalfOpen
[DEBUG] HalfOpen --> Closed
[DEBUG] Closed --> Open
[DEBUG] [spec]total: 1000, allow: 809
[DEBUG] Open --> HalfOpen
[DEBUG] HalfOpen --> Closed
[DEBUG] Closed --> Open
[DEBUG] Open --> HalfOpen
[DEBUG] HalfOpen --> Closed
[DEBUG] Closed --> Open
[DEBUG] Open --> HalfOpen
[DEBUG] HalfOpen --> Closed
[DEBUG] Closed --> Open
[DEBUG] Open --> HalfOpen
[DEBUG] HalfOpen --> Closed
[DEBUG] Closed --> Open
[DEBUG] Open --> HalfOpen
[DEBUG] HalfOpen --> Closed
[DEBUG] Closed --> Open
[DEBUG] Open --> HalfOpen
[DEBUG] HalfOpen --> Closed
[DEBUG] Closed --> Open
[DEBUG] [spec]total: 1000, allow: 24
[DEBUG] Closed --> Open
[DEBUG] Open --> HalfOpen
[DEBUG] Open --> HalfOpen
[DEBUG] HalfOpen --> Closed
[DEBUG] Closed --> Open
[DEBUG] Open --> HalfOpen
[DEBUG] HalfOpen --> Closed
[DEBUG] Closed --> Open
[DEBUG] Open --> HalfOpen
[DEBUG] HalfOpen --> Closed
[DEBUG] [rand]total: 1000, allow: 517
PASS
ok      github.com/moxiaomomo/circuitbreaker    24.737s
```