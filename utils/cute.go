package utils

import (
	"reflect"
)

//obj must be `pointer`of some struct
func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	var s string
	var data = make(map[string]interface{})
	for i := 0; i < t.Elem().NumField(); i++ {
		tag := t.Elem().Field(i).Tag.Get("tomap")
		//get := tag.Get("tomap")
		if tag == "" {
			s = t.Elem().Field(i).Name
		} else {
			s = tag
		}
		data[s] = v.Elem().Field(i).Interface()
	}
	return data
}
