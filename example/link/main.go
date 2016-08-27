package main

import (
	"ohetcd"
	"log"
	"time"
	"fmt"
)

type Service struct {
	Name   string `yaml:"name"`
	Custom *C     `yaml:"-",orm:"private"`
}

type C struct {
	Name string
}

func main() {
	c := &C{Name: "C"}
	s := &Service{}
	log.SetFlags(log.Lshortfile)
	log.Println(s)
	data := ohetcd.NewData()
	// Register to /service
	data.Set("/service", s)
	log.Println(s)
	// Pull updates if any
	data.Update()
	log.Println(s)
	// Set to a local field
	s.Custom = c
	// Set to a shared field
	s.Name = "HAHA"
	// Update val to etcd
	data.Save()
	log.Println(s)
	data.Link()
	go func(data *ohetcd.Data, s *Service) {
		for i := 0; i < 10; i ++ {
			s.Name = fmt.Sprintf("%v", i)
			data.Save()
			time.Sleep(time.Second)
		}

	}(data, s)
	time.Sleep(time.Second)
	go func(data *ohetcd.Data, s *Service) {
		for i := 0; i < 10; i ++ {
			log.Println(*s)
			time.Sleep(time.Second)
		}
	}(data, s)
	for {

	}
}
