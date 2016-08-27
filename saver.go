package ohetcd

import (
	"log"
	"strings"
	"reflect"
	"fmt"
	"gopkg.in/yaml.v2"
)

func deepSave(path string, obj interface{}) error {
	log.Printf("PATH %v <- %v\n", path, obj)
	objVal := reflect.ValueOf(obj)

	if objVal.Kind() == reflect.Ptr {
		objVal = objVal.Elem()
	}

	for i := 0; i < objVal.NumField(); i ++ {
		if objVal.Field(i).Kind() == reflect.Slice {
			subPath := strings.Join([]string{path, objVal.Type().Field(i).Name}, "/")
			for j := 0; j < objVal.Field(i).Len(); j ++ {
				itemPath := strings.Join([]string{subPath, fmt.Sprint(j)}, "/")
				log.Println(itemPath)
				deepSave(itemPath, objVal.Field(i).Index(j).Interface())
			}
			// empty the slice field
			objVal.Field(i).SetLen(0)
		}
	}
	marshalAndSave(path, objVal.Interface())
	return nil
}

func marshalAndSave(path string, obj interface{}) error {
	data, err := yaml.Marshal(obj)
	if err != nil {
		log.Println(err)
	}
	log.Printf("%v <- %v", path, string(data))
	return nil
}