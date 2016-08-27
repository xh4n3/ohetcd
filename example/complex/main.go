package main

import (
	"github.com/kr/pretty"
	"log"
	"ohetcd"
)

type Service struct {
	Name   string `yaml:"name"`
	Custom *C     `yaml:"-",orm:"private"`
	Jobs   []*Job `yaml:"jobs"`
}

type Job struct {
	Name string `yaml:"name"`
}

type C struct {
	Name string
}

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	// local data
	c := &C{Name: "C"}
	j1 := &Job{Name: "Job1"}
	j2 := &Job{Name: "Job2"}
	// target data
	s := &Service{}
	s.Name = "Service"
	s.Jobs = []*Job{j1, j2}
	s.Custom = c
	pretty.Println(s)
	data := ohetcd.NewData()
	// Register to /service
	data.Set("/service", s)
	// Update val to etcd
	data.Save()
}
