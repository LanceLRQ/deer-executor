package utils

import (
	"reflect"
)

// 判断数组、切片或Map是否存在某个值
func Contains(array interface{}, obj interface{}) bool {
	arrayValue := reflect.ValueOf(array)
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < arrayValue.Len(); i++ {
			if arrayValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if arrayValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}

	return false
}
