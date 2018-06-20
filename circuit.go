package circuitbreaker

import (
	"sync"
	"time"
)

// Circuit Circuit
type Circuit struct {
	// TimeWindow time-window
	TimeWindow time.Duration
	// ThredsholdPercent fail-count's percent (1~100)
	ThredsholdPercent int
	// ThredsholdCount fail-count in s time-window
	ThredsholdCount int

	// wait a command to finish execute,
	// as only one command can be executed in HalfOpen status
	waitExecFinish bool

	// status Closed, Open, HalfOpen
	status    StatusEnum
	failCount int
	sucCount  int
	mutex     *sync.RWMutex
}

// Circuits Circuits
type Circuits struct {
	mutex     *sync.RWMutex
	Instances map[string]*Circuit
	// DefaultTimeWindow time-window
	DefaultTimeWindow time.Duration
	// DefaultThredsholdPercent fail-count's percent (1~100)
	DefaultThredsholdPercent int
	// DefaultThredsholdCount fail-count in s time-window
	DefaultThredsholdCount int
}
