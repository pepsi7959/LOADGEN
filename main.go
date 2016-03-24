package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"math/rand"
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
)

func InitRequest(rangeUID int, prefix string) {

	//Generate DN
	for i := 1; i < rangeUID; i++ {
		strUID := fmt.Sprintf("%s,uid=%015d,o=ais,dc=subscriber,dc=C-NTDB", prefix, i)
		ListOfDN = append(ListOfDN, strUID)
	}
	//fmt.Println(ListOfDN[0])

	//TODO: Generate Data
}

func RequestGen(Req *http.Request) *http.Request {
	Req.RequestURI = ListOfDN[rand.Intn(len(ListOfDN))]
	fmt.Println(Req)
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

	//Initalize and Declare
	InitRequest(1000, *PrefixDN)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	req, err := http.NewRequest(*Method, *Hosts, nil)
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
