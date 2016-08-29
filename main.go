package ohetcd

import (
	"github.com/coreos/etcd/client"
	"log"
	"time"
)

const (
	CH_CLOSE_SIG = 0
)

var (
	// channel map to stop watchers
	chMap       = make(map[*Data]chan int)
	etcdCluster = "http://127.0.0.1:2379"
	kapi        client.KeysAPI
)

func SetEtcd(cluster string) {
	etcdCluster = cluster
}

func getKapi() client.KeysAPI {
	if kapi == nil {
		// init connection
		cfg := client.Config{
			Endpoints: []string{etcdCluster},
			Transport: client.DefaultTransport,
			// set timeout per request to fail fast when the target endpoint is unavailable
			HeaderTimeoutPerRequest: time.Second,
		}
		c, err := client.New(cfg)
		if err != nil {
			log.Fatal(err)
		}
		kapi = client.NewKeysAPI(c)
		return kapi
	} else {
		return kapi
	}
}
