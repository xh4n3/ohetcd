package ohetcd

import (
	"github.com/coreos/etcd/client"
	"log"
	"time"
)

const (
	ETCD_CLUSTER = "http://127.0.0.1:2379"
	CH_CLOSE_SIG = 0
	CH_OPEN_SIG  = 1
)

var (
	// global channel map
	heartBeatChMap = make(map[*Data]chan int)
	// watch map
	watchChMap = make(map[*Data]chan int)
	kapi       client.KeysAPI
)

func init() {
	// init connection
	cfg := client.Config{
		Endpoints: []string{ETCD_CLUSTER},
		Transport: client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}
	c, err := client.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	kapi = client.NewKeysAPI(c)

	// Start syncLoop
	go syncLoop()
}

func syncLoop() {
	for {
		for _, ch := range heartBeatChMap {
			ch <- CH_OPEN_SIG
		}
		time.Sleep(time.Second)
	}
}
