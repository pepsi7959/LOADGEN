package main

import (
	"crypto/tls"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type client struct {
	N       int
	C       int
	Tps     int
	Timeout int
	Request *http.Request
	Queues  chan *http.Request
	Done    chan int
	Wg      *sync.WaitGroup
}

func InitRequest(rangeUID int, prefix string) {

	//Generate DN
	for i := 1; i < rangeUID; i++ {
		strUID := fmt.Sprintf("%s,uid=%015d,o=ais,dc=subscriber,dc=C-NTDB", prefix, i)
		ListOfDN = append(ListOfDN, strUID)
	}

	for j := 1; j < 10; j++ {
		strUID := fmt.Sprintf("{\"ds3politeDegree\":\"polite\"}")
		ListData = append(ListData, strUID)
	}
	//fmt.Println(ListOfDN[0])
	//TODO: Generate Data
}

func RequestGen(Req *http.Request) *http.Request {
	Req.RequestURI = ListOfDN[rand.Intn(len(ListOfDN))]
	//fmt.Println(Req)
	return Req
}

func (c *client) DoRequest(clnt *http.Client) {
	NewReq := new(http.Request)
	NewReq = RequestGen(c.Request)
	resp, err := clnt.Do(NewReq)
	if err == nil {
		fmt.Println(resp)
		resp.Body.Close()
	} else {
		//	fmt.Println("sent fail")
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
		//fmt.Println(id, "=> Send...", i)
		c.DoRequest(client)
	}

	c.Done <- 1
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
	if *APIVersion == 1 {
		c.InvokeWorker()
	} else if *APIVersion == 2 {
		c.InvokeWorkerV2()
	} else {
		fmt.Println("Unimplemented version : %d", APIVersion)
	}
}
