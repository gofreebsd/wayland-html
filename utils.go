package main

import (
	"reflect"
	"unsafe"
	"github.com/fangyuanziti/wayland-html/cfn"
)

func cPtr(goPtr unsafe.Pointer) *[0]byte {
	return (*[0]byte)(goPtr)
}

func get_echo(tag string) func() {
	return func() {
		println(tag)
	}
}

var once_funcs = make(map[*interface{}]*cfn.CFn)

func once_func(f interface{}) *cfn.CFn {

	var wrapperFn interface{}

	wrapperFn = reflect.MakeFunc(
		reflect.TypeOf(f),
		func(args []reflect.Value) (results []reflect.Value) {
			results = reflect.ValueOf(f).Call(args)

			if &wrapperFn != nil {
				delete(once_funcs, &wrapperFn)
			}

			return
		},
	).Interface()

	cfn := cfn.CreateFunc(wrapperFn)

	once_funcs[&wrapperFn] = cfn

	return cfn
}
