package main

import (
	"reflect"
)

func isArray(a interface{}) bool {
	var v reflect.Value
	v = reflect.ValueOf(a)
	

	var k reflect.Kind
	k = v.Kind()
	
	
	if (k == reflect.Array) {
	 return true
	}
	return false
}
