package morpheus

import (
	"strconv"
)

// This file contains some helper methods for things like
// converting values to one type or another.
// Suprisingly this is not as simple as it should be
// but I probably don't know what I am doing though...
// Hey, let's convert all our IDs to strings though, for real.

func toInt64(i interface{}) int64 {
	return stringToInt64(i.(string))
}

func stringToInt64(s string) int64 {
	v, _ := strconv.ParseInt(s, 10, 64)
	return v
}

func intToString(n int) string {
	return strconv.FormatInt(int64(n), 10)
}

func int64ToString(n int64) string {
	return strconv.FormatInt(n, 10)
}

// func toString(i interface{}) string {
// 	value, ok := i.(string)
// 	if ok != true {
// 		panic(fmt.Sprintf("toString() failed to convert value '%v'", i)) // oh dear
// 	}
// 	return value
// }

// func toBool(i interface{}) bool {
// 	// if i == nil {
// 	// 	return false
// 	// }
// 	value, ok := i.(bool)
// 	if ok != true {
// 		panic(fmt.Sprintf("toString() failed to convert value '%v'", i)) // oh dear
// 	}
// 	return value
// }

// func StringToInt(v string) int {
// 	value, _ := strconv.ParseInt(v, 10, 32)
// 	return value
// }

// func StringToInt64(v string) int64 {
// 	value, _ := strconv.ParseInt(v, 10, 64)
// 	return value
// }

// This should all go away...please.
// Maybe that stuff up there too.

// Bool returns a pointer to the bool value passed in.
func Bool(v bool) *bool {
	return &v
}

// BoolValue returns the value of the bool pointer passed in or
// false if the pointer is nil.
func BoolValue(v *bool) bool {
	if v != nil {
		return *v
	}
	return false
}

// BoolSlice converts a slice of bool values into a slice of
// bool pointers
func BoolSlice(src []bool) []*bool {
	dst := make([]*bool, len(src))
	for i := 0; i < len(src); i++ {
		dst[i] = &(src[i])
	}
	return dst
}

// BoolValueSlice converts a slice of bool pointers into a slice of
// bool values
func BoolValueSlice(src []*bool) []bool {
	dst := make([]bool, len(src))
	for i := 0; i < len(src); i++ {
		if src[i] != nil {
			dst[i] = *(src[i])
		}
	}
	return dst
}

// Int64 returns a pointer to the int64 value passed in.
// func Int64(v int64) *int64 {
// 	return &v
// }

// // Int64Value returns the value of the int64 pointer passed in or
// // 0 if the pointer is nil.
// func Int64Value(v *int64) int64 {
// 	if v != nil {
// 		return *v
// 	}
// 	return 0
// }

// // Int64Slice converts a slice of int64 values into a slice of
// // int64 pointers
// func Int64Slice(src []int64) []*int64 {
// 	dst := make([]*int64, len(src))
// 	for i := 0; i < len(src); i++ {
// 		dst[i] = &(src[i])
// 	}
// 	return dst
// }

// // Int64ValueSlice converts a slice of int64 pointers into a slice of
// // int64 values
// func Int64ValueSlice(src []*int64) []int64 {
// 	dst := make([]int64, len(src))
// 	for i := 0; i < len(src); i++ {
// 		if src[i] != nil {
// 			dst[i] = *(src[i])
// 		}
// 	}
// 	return dst
// }
