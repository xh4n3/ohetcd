package ohetcd

import (
	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
	"gopkg.in/yaml.v2"
	"log"
	"strings"
)

type Node interface {
	Set()
	Update()
	Save()
	Watch()
	Unwatch()
}

type Data struct {
	Directory string
	Object    interface{}
	Deep      bool
}

func NewData() *Data {
	return &Data{}
}

// set etcd path for data
func (d *Data) Set(dir string, object interface{}, deep bool) {
	d.Directory = dir
	d.Object = object
	d.Deep = deep
}

// get up-to-date value from etcd
func (d *Data) Update() {
	if d.Deep {
		deepRetrieve(d.Directory, d.Object)
		return
	}
	resp, err := getKapi().Get(context.Background(), d.Directory, &client.GetOptions{
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
func (d *Data) Save() error {
	if d.Deep {
		return deepSave(d.Directory, d.Object)
	}
	val, err := yaml.Marshal(d.Object)
	if err != nil {
		return err
	}
	_, err = getKapi().Set(context.Background(), d.Directory, string(val), &client.SetOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (d *Data) Watch() {
	watcher := getKapi().Watcher(d.Directory, &client.WatcherOptions{
		Recursive: true,
	})
	ch := make(chan int)
	chMap[d] = ch
	go func(watcher client.Watcher, d *Data, ch chan int) {
		for {
			select {
			case <-ch:
				return
			default:
				resp, err := watcher.Next(context.Background())
				if err != nil {
					log.Println(err)
				} else {
					log.Println(resp)
					d.Update()
				}
			}
		}
	}(watcher, d, ch)
}

func (d *Data) Unwatch() {
	chMap[d] <- CH_CLOSE_SIG
}
