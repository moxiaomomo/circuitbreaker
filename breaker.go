package circuitbreaker

import (
	"sync"
	"time"

	"github.com/moxiaomomo/circuitbreaker/logger"
)

// NewCirucuitBreaker NewCirucuitBreaker
func NewCirucuitBreaker(timeWin time.Duration, failCnt int, failPercent int) *Circuits {
	// valid time window: 1s~1h
	if timeWin < time.Second {
		timeWin = time.Second
	} else if timeWin > time.Hour {
		timeWin = time.Hour
	}

	// priority: failCnt > failPercent
	return &Circuits{
		Instances:                make(map[string]*Circuit),
		DefaultTimeWindow:        timeWin,
		DefaultThredsholdCount:   failCnt,
		DefaultThredsholdPercent: failPercent,
		mutex: &sync.RWMutex{},
	}
}

// RegisterCommandAsDefault Register command using default settings
func (cb *Circuits) RegisterCommandAsDefault(cmd string) bool {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	if _, ok := cb.Instances[cmd]; ok {
		return false
	}
	cb.Instances[cmd] = &Circuit{
		mutex:             &sync.RWMutex{},
		status:            StatusClosed,
		TimeWindow:        cb.DefaultTimeWindow,
		ThredsholdCount:   cb.DefaultThredsholdCount,
		ThredsholdPercent: cb.DefaultThredsholdPercent,
	}
	return true
}

// RegisterCommand Register command using given settings
func (cb *Circuits) RegisterCommand(cmd string, timeWin time.Duration, failCnt int, failPercent int) bool {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	if _, ok := cb.Instances[cmd]; ok {
		return false
	}
	cb.Instances[cmd] = &Circuit{
		mutex:             &sync.RWMutex{},
		waitExecFinish:    false,
		status:            StatusClosed,
		TimeWindow:        timeWin,
		ThredsholdCount:   failCnt,
		ThredsholdPercent: failPercent,
	}
	return true
}

// Report report command status after execute (suc or fail);
// returns report result, suc or not
func (cb *Circuits) Report(cmd string, isSuc bool) bool {
	cb.mutex.Lock()
	cb.mutex.Unlock()

	if _, ok := cb.Instances[cmd]; !ok {
		return false
	}
	return cb.Instances[cmd].Report(isSuc)
}

// AllowExec is allow to execute command
func (cb *Circuits) AllowExec(cmd string) bool {
	cb.mutex.Lock()
	cb.mutex.Unlock()

	if _, ok := cb.Instances[cmd]; !ok {
		return false
	}
	return cb.Instances[cmd].AllowExec()
}

// Report update a circuit state
func (c *Circuit) Report(isSuc bool) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if isSuc {
		c.sucCount++
	} else {
		c.failCount++
	}

	c.waitExecFinish = false
	c._updateStatus(isSuc)
	return true
}

// AllowExec is allow to execute command
func (c *Circuit) AllowExec() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.status == StatusClosed {
		return true
	} else if c.status == StatusOpen {
		return false
	}

	// in HalfOpen status
	if c.waitExecFinish {
		return false
	}
	c.waitExecFinish = true
	return true
}

func (c *Circuit) updateStatusAsTimeout() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.status == StatusOpen {
		c.status = StatusHalfOpen
		c.waitExecFinish = false
		logger.Debug("Open --> HalfOpen")
	}
}

func (c *Circuit) _updateStatus(isSuc bool) {
	fper := c.failCount * 100 / (c.failCount + c.sucCount)
	switch c.status {
	case StatusClosed:
		if c.failCount >= c.ThredsholdCount || fper >= c.ThredsholdPercent {
			c.status = StatusOpen
			time.AfterFunc(c.TimeWindow, func() {
				c.updateStatusAsTimeout()
			})
			logger.Debug("Closed --> HalfOpen")
		}
	case StatusHalfOpen:
		if isSuc {
			c._reset(StatusClosed)
			logger.Debug("HalfOpen --> Closed")
		} else {
			c.status = StatusOpen
			time.AfterFunc(c.TimeWindow, func() {
				c.updateStatusAsTimeout()
			})
			logger.Debug("HalfOpen --> Open")
		}
	}
}

func (c *Circuit) _reset(status StatusEnum) {
	c.waitExecFinish = false
	c.status = status
	c.failCount = 0
	c.sucCount = 0
}
