package main

import (
	"crypto/tls"
	"net/http"
	"time"
)

func (c *client) Dispatcher(id, numOfReq int) {
	tick := time.Tick(time.Duration(1e6/(c.Tps)) * time.Microsecond)
	for i := 0; i < numOfReq; i++ {
		<-tick
		c.Queues <- &http.Request{}
	}

	c.Done <- 1
}

func (c *client) WorkerV2() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := &http.Client{Transport: tr}

	for {
		<-c.Queues
		c.DoRequest(client)
	}
}

func (c *client) InvokeWorkerV2() {
	for i := 0; i < c.C; i++ {
		go c.Dispatcher(i, *NumOfReq / *Concur)
		go c.WorkerV2()
	}
}
