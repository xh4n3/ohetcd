package main

import (
	"fmt"
	"github.com/kr/pretty"
	"log"
	"reflect"
	"kubernetes/pkg/util/json"
	"strings"
)

type Service struct {
	Name string  `yaml:"name"`
	Ints []*Item `yaml:"ints"`
}

type Item struct {
	Name string `yaml:"name"`
	Book string `yaml:"Book"`
}

func main() {
	log.SetFlags(log.Lshortfile)
	s := &Service{Name: "s", Ints: []*Item{
		{
			Name: "i1",
			Book: "b1",
		},
		{
			Name: "i2",
			Book: "b2",
		},
	}}
	// Original
	pretty.Println(s)
	fmt.Printf("%p\n", s.Ints[0])
	// Add a new item
	addItem(s)
	pretty.Println(s)
	// Address of 1st Item is unchanged
	fmt.Printf("%p\n", s.Ints[0])
	fmt.Printf("%p\n", s.Ints[2])
	modifyItem(s)
	// Remove an item
	removeItem(s)
	pretty.Println(s)
	// Address of 1st Item is unchanged
	fmt.Printf("%p\n", s.Ints[0])
	fmt.Printf("%p\n", s.Ints[1])
}

func addItem(s *Service) {
	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	slice := val.FieldByName("Ints")
	sliceType := slice.Type()
	// get item's type in slice
	itemType := sliceType.Elem()
	if itemType.Kind() == reflect.Ptr {
		itemType = itemType.Elem()
	}
	newItem := reflect.New(itemType)
	newItem.Elem().FieldByName("Name").SetString("New")
	slice.Set(reflect.Append(slice, newItem))
}

func modifyItem(s *Service) {
	jsonData, err := json.Marshal(s.Ints[0])
	if err != nil {
		log.Fatalln(err)
	}
	changed := strings.Replace(string(jsonData), "b1", "book1", 1)
	target := reflect.ValueOf(s).Elem().FieldByName("Ints").Index(0)
	err = json.Unmarshal([]byte(changed), target.Interface())
	if err != nil {
		log.Fatalln(err)
	}
	pretty.Println(target)
}

func removeItem(s *Service) {
	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	slice := val.FieldByName("Ints")
	// remove 2nd item
	i := 1
	resultSlice := reflect.AppendSlice(slice.Slice(0, i), slice.Slice(i+1, slice.Len()))
	slice.Set(resultSlice)
}

