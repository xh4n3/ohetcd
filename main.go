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
)

var kapi client.KeysAPI

type Linkable interface {
	Set()
	Update()
	Link(dir string)
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

func (d *Data) Update() {
	resp, err := kapi.Get(context.Background(), d.Directory, &client.GetOptions{
		Recursive: true,
	})
	if err != nil {
		if strings.Contains(err.Error(), string(client.ErrorCodeKeyNotFound)) {
			d.Save()
			d.Update()
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

func (d *Data) Link(dir string) {
	log.Println("LINK")
}

func (d *Data) Unlink() {
	log.Println("UNLINK")
}

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
}

func NewData() *Data {
	return &Data{}
}
