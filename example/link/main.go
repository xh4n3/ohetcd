package main

import (
	"ohetcd"
	"log"
	"time"
)

type Service struct {
	Name   string `yaml:"name"`
	Custom *C     `yaml:"-",orm:"private"`
}

type C struct {
	Name string
}

func main() {
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
	data.Link()
	time.Sleep(time.Second)
	go func(data *ohetcd.Data, s *Service) {
		for i := 0; i < 20; i ++ {
			log.Println(s.Name)
			time.Sleep(time.Second)
		}
	}(data, s)
	time.Sleep(10*time.Second)
	data.Unlink()
	for {

	}
}
