package utils

import (
	"github.com/gogf/gf/v2/text/gstr"
	"reflect"
)

//IsEmpty 空字符串判断
func IsEmpty(str string) bool {
	if len(gstr.TrimAll(str)) == 0 {
		return true
	}
	return false
}

//InArray 判断是否在切片集合中
func InArray[T any](t T, list []T) bool {
	for _, item := range list {
		if reflect.DeepEqual(t, item) {
			return true
		}
	}
	return false
}
