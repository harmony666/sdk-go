package main

import (
	"flag"
	"sync"
)

var txsMap = sync.Map{}

type transaction struct {
	requestId string //only as request id, maybe not txId
	startTime int64
	endTime   int64
	success   bool
}

var (
	threadNum int
	loopNum   int
	execTimes int
	method    string
	interval  int
)

func init() {
	flag.IntVar(&threadNum, "threadNum", 10000, "total thread number, default 10000")
	flag.IntVar(&loopNum, "loopNum", 1, "total loop number per thread, default 1")
	flag.IntVar(&execTimes, "execTimes", 5, "exec times, default 5")
	flag.IntVar(&interval, "interval", 0, "interval between sendRequest, default 0 millisecond")
	flag.StringVar(&method, "method", "increaseBalance", "contract method to execute")
}
