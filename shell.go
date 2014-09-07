package main

/*
#cgo pkg-config: wayland-server

#include <wayland-server.h>

*/
import "C"

import (
	"unsafe"
)

type shell_surface_interface C.struct_wl_shell_surface_interface

var pong = create_func(
	get_echo("pong"),
)
var move = create_func(
	get_echo("move"),
)
var shell_resize = create_func(
	get_echo("shell_resize"),
)
var set_toplevel = create_func(
	get_echo("set_toplevel"),
)

var set_transient = create_func(
	get_echo("set_transient"),
)
var set_fullscreen = create_func(
	get_echo("set_fullscreen"),
)
var set_popup = create_func(
	get_echo("set_popup"),
)
var set_maximized = create_func(
	get_echo("set_maximized"),
)
var set_title = create_func(
	get_echo("set_title"),
)
var set_class = create_func(
	get_echo("set_class"),
)

var shell_surface_impl = shell_surface_interface{
	pong:           cPtr(pong.fn_ptr),
	move:           cPtr(move.fn_ptr),
	resize:         cPtr(shell_resize.fn_ptr),
	set_toplevel:   cPtr(set_toplevel.fn_ptr),
	set_transient:  cPtr(set_transient.fn_ptr),
	set_fullscreen: cPtr(set_fullscreen.fn_ptr),
	set_popup:      cPtr(set_popup.fn_ptr),
	set_maximized:  cPtr(set_maximized.fn_ptr),
	set_title:      cPtr(set_title.fn_ptr),
	set_class:      cPtr(set_class.fn_ptr),
}

var get_shell_surface = create_func(
	func(client *C.struct_wl_client,
		resource *C.struct_wl_resource,
		id C.uint32_t,
		surface_resource *C.struct_wl_resource) {
		println("get_shell_surface")
		shell_surface_res := C.wl_resource_create(client,
			&C.wl_shell_surface_interface,
			C.wl_resource_get_version(resource),
			id)
		C.wl_resource_set_implementation(
			shell_surface_res,
			(unsafe.Pointer)(&shell_surface_impl),
			nil,
			nil)
	},
)

var shell_impl = C.struct_wl_shell_interface{
	get_shell_surface: cPtr(get_shell_surface.fn_ptr),
}

var bind_shell = create_func(
	func(client *C.struct_wl_client, data unsafe.Pointer,
		version C.int, id C.uint32_t) {

		if version >= 1 {
			version = 1
		}

		resource := C.wl_resource_create(client, &C.wl_shell_interface, version, id)

		C.wl_resource_set_implementation(
			resource,
			(unsafe.Pointer)(&shell_impl),
			nil,
			nil)
	},
)

func shellInit(display *C.struct_wl_display) {

	C.wl_global_create(display,
		&C.wl_shell_interface,
		1,
		nil,
		cPtr(bind_shell.fn_ptr))
}
