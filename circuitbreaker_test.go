package circuitbreaker

import (
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/moxiaomomo/circuitbreaker/logger"
)

func doTest(cbs *Circuits, testcmd string, factor int) {
	allowCnt := 0
	totalCnt := 1000
	for i := 1; i <= totalCnt; i++ {
		if cbs.AllowExec(testcmd) {
			allowCnt++
		}
		if i%factor == 0 {
			cbs.Report(testcmd, false)
		} else {
			cbs.Report(testcmd, true)
		}
		time.Sleep(time.Millisecond * 6)
	}
	logger.Debugf("[spec]total: %d, allow: %d\n", totalCnt, allowCnt)
}

func doTestRandom(cbs *Circuits, testcmd string) {
	allowCnt := 0
	totalCnt := 1000
	for i := 1; i <= totalCnt; i++ {
		if cbs.AllowExec(testcmd) {
			allowCnt++
		}
		if rand.Intn(30) > (rand.Intn(50)+1)*(rand.Intn(10)+1) {
			cbs.Report(testcmd, false)
		} else {
			cbs.Report(testcmd, true)
		}
		time.Sleep(time.Millisecond * time.Duration(5+rand.Intn(3)))
	}
	logger.Debugf("[rand]total: %d, allow: %d\n", totalCnt, allowCnt)
}

func TestCircuitBreaker(t *testing.T) {
	logger.SetLogLevel(logger.LOG_DEBUG)
	testcmd := "testcmd"

	cbs := NewCirucuitBreaker(time.Second, 100, 10)
	suc := cbs.RegisterCommandAsDefault(testcmd)
	if !suc {
		os.Exit(1)
	}
	logger.Debugf("To run test for command: %s\n", testcmd)

	// make the 2 percent of the reports is failed
	doTest(cbs, testcmd, 50)
	// make the 10 percent of the reports is failed
	doTest(cbs, testcmd, 10)
	// make the 20 percent of the reports is failed
	doTest(cbs, testcmd, 5)

	cbs = NewCirucuitBreaker(time.Second, 150, 20)
	suc = cbs.RegisterCommandAsDefault(testcmd)
	if !suc {
		os.Exit(1)
	}
	doTestRandom(cbs, testcmd)
}
