package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
)

var (
	Concur      = flag.Int("C", 1, "Concurrence")
	NumOfReq    = flag.Int("N", -1, "Number of request")
	Tps         = flag.Int("TPS", 1, "TPS limits")
	Cpu         = flag.Int("CPU", 8, "TPS limits")
	Hosts       = flag.String("H", "http://127.0.0.1:8080", "host:port")
	PrefixDN    = flag.String("prefix", "ds=profile", "prefix Of DN")
	Method      = flag.String("method", "GET", "GET method")
	Body        = flag.String("d", "", "Body of HTTP")
	FilePath    = flag.String("@d", "", "FilePath")
	ListOfDN    []string
	ListData    []string
	NumberOfUID = flag.Int("maxUID", 1000, "Miximum UID")
	APIVersion  = flag.Int("V", 1, "Version")
)

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(*Cpu)

	//Initalize and Declare
	InitRequest(1000, *PrefixDN)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	req, err := http.NewRequest(*Method, *Hosts, nil)
	if err != nil {
		os.Exit(0)
	}
	clnt := &client{
		N:       *NumOfReq,
		C:       *Concur,
		Tps:     *Tps,
		Timeout: 100,
		Request: req,
		Queues:  make(chan *http.Request, 10000),
		Wg:      &sync.WaitGroup{},
		Done:    make(chan int),
	}

	clnt.Run()

	var i int = 0
	isExite := false
	for {
		select {
		case <-c:
			isExite = true
		case <-clnt.Done:
			i++
		}

		if i >= clnt.C || isExite {
			break
		}
	}
	fmt.Println("close : ", i)
	Console()
	Status()
}
