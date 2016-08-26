package main

import (
	"ohetcd"
	"log"
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
}
