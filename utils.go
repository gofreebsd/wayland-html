package main

import "unsafe"

func cPtr(goPtr unsafe.Pointer) *[0]byte {
	return (*[0]byte)(goPtr)
}

func get_echo(tag string) func() {
	return func() {
		println(tag)
	}
}
