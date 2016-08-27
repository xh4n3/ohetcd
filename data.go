package ohetcd

import (
	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
	"gopkg.in/yaml.v2"
	"log"
	"strings"
)

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

func NewData() *Data {
	return &Data{}
}

// set etcd path for data
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
	ch := make(chan int)
	// save ch into global chMap
	chMap[d] = ch
	go func(d *Data, ch chan int) {
		var i int
		for {
			i = <-ch
			if i == CH_CLOSE_SIG {
				break
			} else {
				d.Update()
			}
		}
	}(d, ch)
	log.Printf("Linked to %v\n", d.Directory)
}

func (d *Data) Unlink() {
	chMap[d] <- CH_CLOSE_SIG
	log.Printf("Unlinked from %v\n", d.Directory)
}
