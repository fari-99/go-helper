package gohelper

import (
	"fmt"
	"reflect"

	"golang.org/x/exp/slices"
)

// InArray checks whether needle is in haystack.
// if haystack is an array of struck, then needle need to be a function
// example : https://stackoverflow.com/questions/38654383/how-to-search-for-an-element-in-a-golang-slice
func InArray(needle any, haystack []any) (bool, int, error) {
	haystackValue := reflect.ValueOf(haystack)
	checkData := haystackValue.Index(0)
	checkDataKind := reflect.ValueOf(checkData).Kind()
	needleValue := reflect.ValueOf(needle)
	needleKind := needleValue.Kind()

	// haystack is empty array
	if !checkData.IsValid() {
		return false, -1, nil
	}

	isStruct := false
	if checkDataKind == reflect.Struct && needleKind == reflect.Func {
		isStruct = true
	}

	if !isStruct {
		for i := 0; i < haystackValue.Len(); i++ {
			hayVal := haystackValue.Index(i).Interface()

			if reflect.DeepEqual(hayVal, needle) {
				return true, i, nil
			}
		}
	} else {
		if needleValue.Kind() != reflect.Func {
			return false, -1, fmt.Errorf("must be a function")
		}

		index := slices.IndexFunc(haystack, needle.(func(any) bool))
		if index != -1 {
			return true, index, nil
		}
	}

	return false, -1, nil
}
