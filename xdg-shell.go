package main

/*

#include <wayland-server.h>
#include "xdg-shell-server-protocol.h"

*/
import "C"

import (
	"unsafe"
)

var xdg_surface_destroy = create_func(
	func(client *C.struct_wl_client,
		resource *C.struct_wl_resource) {

		println("xdg_surface_destroy")
		C.wl_resource_destroy(resource)
	},
)

var xdg_surface_set_transient_for = create_func(
	func(client *C.struct_wl_client,
		resource *C.struct_wl_resource,
		parent *C.struct_wl_resource) {

		println("xdg_surface_transient")
	},
)

var xdg_surface_set_margin = create_func(
	func(client *C.struct_wl_client,
		resource *C.struct_wl_resource,
		left_margin C.int32_t,
		right_margin C.int32_t,
		top_margin C.int32_t,
		bottom_margin C.int32_t) {

		println("xdg_surface_margin")
	},
)

var xdg_surface_set_title = create_func(
	func(client *C.struct_wl_client,
		resource *C.struct_wl_resource,
		title *C.char) {

		println("xdg_surface_title")
	},
)

var xdg_surface_set_app_id = create_func(
	func(client *C.struct_wl_client,
		resource *C.struct_wl_resource,
		app_id *C.char) {

		println("xdg_surface_app_id")
	},
)

var xdg_surface_move = create_func(
	func(client *C.struct_wl_client,
		resource *C.struct_wl_resource,
		seat *C.struct_wl_resource,
		serial C.uint32_t) {

		println("xdg_surface_move")
	},
)

var xdg_surface_resize = create_func(
	func(client *C.struct_wl_client,
		resource *C.struct_wl_resource,
		seat *C.struct_wl_resource,
		serial C.uint32_t,
		edges C.uint32_t,
	) {

		println("xdg_surface_resize")
	},
)

var xdg_surface_set_output = create_func(
	func(client *C.struct_wl_client,
		resource *C.struct_wl_resource,
		output *C.struct_wl_resource) {

		println("xdg_surface_set_output")
	},
)

var xdg_surface_request_change_state = create_func(
	func(client *C.struct_wl_client,
		resource *C.struct_wl_resource,
		state_type C.uint32_t,
		value C.uint32_t,
		serial C.uint32_t) {

		println("xdg_surface_request_change_state")
	},
)

var xdg_surface_ack_change_state = create_func(
	func(client *C.struct_wl_client,
		resource *C.struct_wl_resource,
		state_type C.uint32_t,
		value C.uint32_t,
		serial C.uint32_t) {

		println("xdg_surface_ack_change_state")
	},
)

var xdg_surface_set_minized = create_func(
	func(client *C.struct_wl_client,
		resource *C.struct_wl_resource) {

		println("xdg_surface_set_minized")
	},
)

var xdg_surface_impl = C.struct_xdg_surface_interface{
	destroy:              (cPtr)(xdg_surface_destroy.fn_ptr),
	set_transient_for:    (cPtr)(xdg_surface_set_transient_for.fn_ptr),
	set_margin:           (cPtr)(xdg_surface_set_margin.fn_ptr),
	set_title:            (cPtr)(xdg_surface_set_title.fn_ptr),
	set_app_id:           (cPtr)(xdg_surface_set_app_id.fn_ptr),
	move:                 (cPtr)(xdg_surface_set_title.fn_ptr),
	resize:               (cPtr)(xdg_surface_resize.fn_ptr),
	set_output:           (cPtr)(xdg_surface_set_output.fn_ptr),
	request_change_state: (cPtr)(xdg_surface_request_change_state.fn_ptr),
	ack_change_state:     (cPtr)(xdg_surface_ack_change_state.fn_ptr),
	set_minimized:        (cPtr)(xdg_surface_set_minized.fn_ptr),
}

var use_unstable_version = create_func(
	func(client *C.struct_wl_client,
		resource *C.struct_wl_resource,
		version C.int32_t) {

		println("unstable")
	},
)

var get_xdg_surface = create_func(
	func(client *C.struct_wl_client,
		resource *C.struct_wl_resource,
		id C.uint32_t,
		surface *C.struct_wl_resource) {

		surface_res := C.wl_resource_create(client, &C.xdg_surface_interface,
			C.wl_resource_get_version(resource), id)

		println("get_xdg_surface", int(id))

		C.wl_resource_set_implementation(surface_res,
			(unsafe.Pointer)(&xdg_surface_impl),
			nil,
			nil)
	},
)

var get_xdg_popup = create_func(
	func(client *C.struct_wl_client,
		resource *C.struct_wl_resource,
		id C.uint32_t,
		surface *C.struct_wl_resource,
		parent *C.struct_wl_resource,
		seat *C.struct_wl_seat,
		x C.int32_t,
		y C.int32_t,
		flags C.uint32_t) {

		println("get_xdg_popup", int(id))
	},
)

var xdg_pong = create_func(
	func(client *C.struct_wl_client,
		resource *C.struct_wl_resource,
		serial C.uint32_t) {

	},
)

var xdg_shell_impl = C.struct_xdg_shell_interface{
	use_unstable_version: (cPtr)(use_unstable_version.fn_ptr),
	get_xdg_surface:      (cPtr)(get_xdg_surface.fn_ptr),
	get_xdg_popup:        (cPtr)(get_xdg_popup.fn_ptr),
	pong:                 (cPtr)(xdg_pong.fn_ptr),
}

var bind_xdg_shell = create_func(
	func(client *C.struct_wl_client, data unsafe.Pointer,
		version C.int, id C.uint32_t) {

		resource := C.wl_resource_create(client, &C.xdg_shell_interface, version, id)

		C.wl_resource_set_implementation(resource,
			(unsafe.Pointer)(&xdg_shell_impl),
			nil,
			nil)
	},
)

func xdgShellInit(display *C.struct_wl_display) {

	C.wl_global_create(display,
		&C.xdg_shell_interface,
		1,
		nil,
		cPtr(bind_xdg_shell.fn_ptr))
}
