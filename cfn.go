package main

/*
#cgo pkg-config: libffi

#include <sys/types.h>
#include <pwd.h>
#include <stdlib.h>
#include <stdio.h>

#include <ffi.h>

#include "cfn.h"
*/
import "C"

import (
	"reflect"
	"runtime"
	"unsafe"
)

func get_c_slice(c_ptr unsafe.Pointer, num int, slice_ptr unsafe.Pointer) {

	sliceHeader := (*reflect.SliceHeader)((slice_ptr))
	sliceHeader.Cap = num
	sliceHeader.Len = num
	sliceHeader.Data = uintptr(c_ptr)

}

//export cfn_go_callback
func cfn_go_callback(p unsafe.Pointer, args unsafe.Pointer, ret unsafe.Pointer) {

	fn_data := (*FuncData)(p)

	fn_data.Call(args, ret)
}

type FuncData struct {
	fn      reflect.Value
	fn_type reflect.Type

	cif C.ffi_cif
}

func (func_data *FuncData) Arg(index int) reflect.Value {
	return reflect.ValueOf(1)
}

func (func_data *FuncData) InTypes() []reflect.Type {
	numIn := func_data.NumIn()

	types := make([]reflect.Type, numIn)

	for i := 0; i < numIn; i++ {
		types[i] = func_data.fn_type.In(i)
	}

	return types
}

func (func_data *FuncData) OutTypes() []reflect.Type {
	numOut := func_data.NumOut()

	types := make([]reflect.Type, numOut)

	for i := 0; i < numOut; i++ {
		types[i] = func_data.fn_type.Out(i)
	}

	return types
}

func (func_data *FuncData) NumIn() int {
	return func_data.fn_type.NumIn()
}

func (func_data *FuncData) NumOut() int {
	return func_data.fn_type.NumOut()
}

func (func_data *FuncData) Call(args unsafe.Pointer, ret unsafe.Pointer) {

	numIn := func_data.fn.Type().NumIn()

	var goArgs []*uint32
	get_c_slice(args, numIn, unsafe.Pointer(&goArgs))

	getArg := func(index int) reflect.Value {
		return reflect.NewAt(func_data.fn_type.In(index), unsafe.Pointer(goArgs[index])).Elem()
	}

	funArgs := make([]reflect.Value, numIn)

	for i := 0; i < numIn; i++ {
		funArgs[i] = getArg(i)
	}

	rets := func_data.fn.Call(funArgs)

	numOut := func_data.fn_type.NumOut()

	if numOut == 1 {
		reflect.NewAt(func_data.fn_type.Out(0), ret).Elem().Set(rets[0])
	}
}

type CFn struct {
	closure unsafe.Pointer
	cif     C.ffi_cif

	fn_ptr  unsafe.Pointer
	fn_data *FuncData
}

func ffi_type(t reflect.Type) *C.ffi_type {

	switch t.Kind() {
	case reflect.Ptr:
		return &C.ffi_type_pointer
	case reflect.Int:
		return &C.ffi_type_sint
	case reflect.Bool:
		return &C.ffi_type_sint
	case reflect.Int8:
		return &C.ffi_type_sint8
	case reflect.Int16:
		return &C.ffi_type_sint16
	case reflect.Int32:
		return &C.ffi_type_sint32
	case reflect.Int64:
		return &C.ffi_type_sint64
	case reflect.Uint8:
		return &C.ffi_type_uint8
	case reflect.Uint16:
		return &C.ffi_type_uint16
	case reflect.Uint32:
		return &C.ffi_type_uint32
	case reflect.Uint64:
		return &C.ffi_type_uint64
	case reflect.Float32:
		return &C.ffi_type_float
	case reflect.Float64:
		return &C.ffi_type_double
	default:
		return &C.ffi_type_pointer
	}
}

func get_args_define(types []reflect.Type, argsNum int) **C.struct__ffi_type {
	if argsNum <= 0 {
		return nil
	}

	numIn := argsNum
	inTypes := types

	cargs := make([]*C.ffi_type, numIn)

	for i, t := range inTypes {
		cargs[i] = ffi_type(t)
	}

	return &cargs[0]

}

func create_func(f interface{}) *CFn {

	cfn := new(CFn)

	closure := C.ffi_closure_alloc(
		C.size_t(unsafe.Sizeof(C.ffi_closure{})),
		&(cfn.fn_ptr))

	if closure == nil {
		return nil
	}

	cfn.closure = closure

	runtime.SetFinalizer(cfn, free_cfn)

	fn_data := new(FuncData)

	fn_data.fn = reflect.ValueOf(f)
	fn_data.fn_type = fn_data.fn.Type()

	cfn.fn_data = fn_data

	numIn := fn_data.NumIn()

	args := get_args_define(fn_data.InTypes(), numIn)

	var ret_type *C.ffi_type

	numOut := fn_data.NumOut()
	if numOut == 1 {
		ret_type = ffi_type(fn_data.OutTypes()[0])
		// fmt.Println((fn_data.OutTypes()[0]).Kind())
	} else {
		// TODO: please handle numOut > 1
		ret_type = &C.ffi_type_void
	}

	if C.ffi_prep_cif(&(fn_data.cif),
		C.FFI_DEFAULT_ABI,
		C.uint(numIn),
		ret_type,
		args) != C.FFI_OK {
		return nil
	}

	if C.ffi_prep_closure_loc(
		(*C.ffi_closure)(closure),
		&(fn_data.cif),
		(*[0]byte)(C.binding),
		unsafe.Pointer(fn_data),
		cfn.fn_ptr) != C.FFI_OK {
		return nil
	}

	return cfn

}


func free_cfn(cfn *CFn) {
	C.ffi_closure_free(cfn.closure)
}

