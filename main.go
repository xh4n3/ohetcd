package ohetcd

import (
	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
	"gopkg.in/yaml.v2"
	"log"
	"strings"
	"time"
)

const (
	ETCD_CLUSTER = "http://127.0.0.1:2379"
	CH_CLOSE_SIG = 0
	CH_OPEN_SIG  = 0
)

var kapi client.KeysAPI

type Linkable interface {
	Set()
	Update()
	Save()
	Link()
	Unlink()
}

type Data struct {
	Directory string
	Object    interface{}
}

func (d *Data) Set(dir string, object interface{}) {
	d.Directory = dir
	d.Object = object
}

// get up-to-date value from etcd
func (d *Data) Update() {
	resp, err := kapi.Get(context.Background(), d.Directory, &client.GetOptions{
		Recursive: true,
	})
	if err != nil {
		if strings.Contains(err.Error(), string(client.ErrorCodeKeyNotFound)) {
			// if key not found, init etcd with current variable
			d.Save()
			return
		}
		log.Println(err)
	}
	data := resp.Node.Value
	err = yaml.Unmarshal([]byte(data), d.Object)
	if err != nil {
		log.Println(err)
	}
}

// changes made, save to etcd
func (d *Data) Save() {
	val, err := yaml.Marshal(d.Object)
	if err != nil {
		log.Println(err)
	}
	resp, err := kapi.Set(context.Background(), d.Directory, string(val), &client.SetOptions{})
	if err != nil {
		log.Println(err)
	}
	log.Println(resp)
}

func (d *Data) Link() {
	log.Println("LINK")
	ch := make(chan int, 1)
	go func(d *Data, ch chan int) {
		var i int
		for {
			i = <- ch
			if i == CH_CLOSE_SIG {
				break
			} else {
				d.Update()
			}
		}
	}(d, ch)
}

func (d *Data) Unlink() {
	log.Println("UNLINK")
	close(chMap[d])
}

var (
	chMap map[*Data]chan int
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

func NewData() *Data {
	return &Data{}
}

func syncLoop() {
	for {
		for _, ch := range chMap {
			ch <- CH_OPEN_SIG
		}
		time.Sleep(time.Second)
	}
}
