package ohetcd

import (
	"log"
	"golang.org/x/net/context"
	"github.com/coreos/etcd/client"
	"reflect"
	"strings"
	"encoding/json"
)

func deepRetrieve(path string, obj interface{}) error {
	resp, err := kapi.Get(context.Background(), path, &client.GetOptions{Recursive: true})
	if err != nil {
		log.Fatalln(err)
	}
	// init struct instance
	var content string
	for _, n := range resp.Node.Nodes {
		if !n.Dir {
			content = n.Value
			log.Println(content)
		}
	}
	err = json.Unmarshal([]byte(content), obj)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(obj)
	val := reflect.ValueOf(obj)
	retrieve(resp.Node, val)
	return nil
}

func retrieve(node *client.Node, val reflect.Value) {
	log.Println(val)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	for _, n := range node.Nodes {
		log.Println(n.Key)
		if n.Dir {
			// f corresponds to an item in slice
			for _, f := range n.Nodes {
				keys := strings.Split(f.Key, "/")
				fieldName := keys[len(keys)-2]
				log.Println(val)
				log.Println(fieldName)
				val = val.FieldByName(fieldName)
				// TODO
				// val is a slice slot
				// we need to parse n.Nodes into a slice
				retrieve(n, val)
			}
		} else {
			log.Println(n.Value)
		}
	}
}
