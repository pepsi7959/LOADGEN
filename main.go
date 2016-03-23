package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

type client struct {
	N       int
	C       int
	Tps     int
	Timeout int
	Request *http.Request
}

var (
	Concur   = flag.Int("C", 1, "Concurrence")
	NumOfReq = flag.Int("N", -1, "Number of request")
	Tps      = flag.Int("TPS", 1, "TPS limits")
	Cpu      = flag.Int("CPU", 8, "TPS limits")
)

func (c *client) DoRequest(clnt *http.Client) {
	for i := 0; i < 100; i++ {

	}

	var size int64
	var code int
	NewReq := new(http.Request)
	NewReq = c.Request
	resp, err := clnt.Do(NewReq)
	if err == nil {
		size = resp.ContentLength
		code = resp.StatusCode
		io.Copy(ioutil.Discard, resp.Body)
		fmt.Println(resp)
		resp.Body.Close()
	} else {
		fmt.Println("sent fail")
	}

}

func (c *client) Worker(id, num int) {
	tick := time.Tick(time.Duration(1e6/(c.Tps)) * time.Microsecond)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := &http.Client{Transport: tr}

	for i := 0; i < num || c.N == -1; i++ {
		<-tick
		fmt.Println(id, "=> Send...", i)
		c.DoRequest(client)
	}
}
func (c *client) InvokeWorker() {
	for i := 0; i < c.C; i++ {
		go c.Worker(i, c.N/c.C)
	}
}
func Status() {
	fmt.Println("Number of Request:")
	fmt.Println("Number of Concurrent :")
	fmt.Println("TPS :")
}
func Console() { fmt.Println("====== LOAD GENERATOR ======") }
func (c *client) Run() {
	c.InvokeWorker()
}
func main() {
	flag.Parse()
	runtime.GOMAXPROCS(*Cpu)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	req, err := http.NewRequest("GET", "http://www.google.co.th", nil)
	if err != nil {
		os.Exit(0)
	}
	(&client{
		N:       *NumOfReq,
		C:       *Concur,
		Tps:     *Tps,
		Timeout: 100,
		Request: req,
	}).Run()

	<-c
	Console()
	Status()
}
