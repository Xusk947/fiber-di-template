package utils

import "time"

func StringPtr(s string) *string {
	return &s
}

func TimePtr(t time.Time) *time.Time {
	return &t
}

func IntPtr(i int) *int {
	return &i
}

func BoolPtr(b bool) *bool {
	return &b
}

func Int64Ptr(i int64) *int64 {
	return &i
}

func Uint64Ptr(i uint64) *uint64 {
	return &i
}

func UintPtr(i uint) *uint {
	return &i
}

func Float64Ptr(f float64) *float64 {
	return &f
}

func Float32Ptr(f float32) *float32 {
	return &f
}

func Int32Ptr(i int32) *int32 {
	return &i
}

func Int16Ptr(i int16) *int16 {
	return &i
}

func Int8Ptr(i int8) *int8 {
	return &i
}

func Uint32Ptr(i uint32) *uint32 {
	return &i
}

func Uint16Ptr(i uint16) *uint16 {
	return &i
}

func Uint8Ptr(i uint8) *uint8 {
	return &i
}

func UintptrPtr(i uintptr) *uintptr {
	return &i
}

func DurationPtr(d time.Duration) *time.Duration {
	return &d
}

func BytePtr(b byte) *byte {
	return &b
}

func RunePtr(r rune) *rune {
	return &r
}

func Complex64Ptr(c complex64) *complex64 {
	return &c
}

func Complex128Ptr(c complex128) *complex128 {
	return &c
}

func SliceIntPtr(i []int) *[]int {
	return &i
}

func MapPtr[K comparable, V any](m map[K]V) *map[K]V {
	return &m
}
