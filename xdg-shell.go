package main

/*
#include "wayland-fix.h"
*/
import "C"

import (
	"github.com/fangyuanziti/wayland-html/cfn"
	"unsafe"
)

var xdg_surface_destroy = cfn.CreateFunc(
	func(client *C.struct_wl_client,
		resource *C.struct_wl_resource) {

		println("xdg_surface_destroy")
		C.wl_resource_destroy(resource)
	},
)

var xdg_surface_set_transient_for = cfn.CreateFunc(
	func(client *C.struct_wl_client,
		resource *C.struct_wl_resource,
		parent *C.struct_wl_resource) {

		println("xdg_surface_transient")
	},
)

var xdg_surface_set_margin = cfn.CreateFunc(
	func(client *C.struct_wl_client,
		resource *C.struct_wl_resource,
		left_margin C.int32_t,
		right_margin C.int32_t,
		top_margin C.int32_t,
		bottom_margin C.int32_t) {

		println("xdg_surface_margin")
	},
)

var xdg_surface_set_title = cfn.CreateFunc(
	func(client *C.struct_wl_client,
		resource *C.struct_wl_resource,
		title *C.char) {

		println("xdg_surface_title")
	},
)

var xdg_surface_set_app_id = cfn.CreateFunc(
	func(client *C.struct_wl_client,
		resource *C.struct_wl_resource,
		app_id *C.char) {

		println("xdg_surface_app_id")
	},
)

var xdg_surface_move = cfn.CreateFunc(
	func(client *C.struct_wl_client,
		resource *C.struct_wl_resource,
		seat *C.struct_wl_resource,
		serial C.uint32_t) {

		println("xdg_surface_move")
	},
)

var xdg_surface_resize = cfn.CreateFunc(
	func(client *C.struct_wl_client,
		resource *C.struct_wl_resource,
		seat *C.struct_wl_resource,
		serial C.uint32_t,
		edges C.uint32_t,
	) {

		println("xdg_surface_resize")
	},
)

var xdg_surface_set_output = cfn.CreateFunc(
	func(client *C.struct_wl_client,
		resource *C.struct_wl_resource,
		output *C.struct_wl_resource) {

		println("xdg_surface_set_output")
	},
)

var xdg_surface_request_change_state = cfn.CreateFunc(
	func(client *C.struct_wl_client,
		resource *C.struct_wl_resource,
		state_type C.uint32_t,
		value C.uint32_t,
		serial C.uint32_t) {

		println("xdg_surface_request_change_state")
	},
)

var xdg_surface_ack_change_state = cfn.CreateFunc(
	func(client *C.struct_wl_client,
		resource *C.struct_wl_resource,
		state_type C.uint32_t,
		value C.uint32_t,
		serial C.uint32_t) {

		println("xdg_surface_ack_change_state")
	},
)

var xdg_surface_set_minized = cfn.CreateFunc(
	func(client *C.struct_wl_client,
		resource *C.struct_wl_resource) {

		println("xdg_surface_set_minized")
	},
)

var xdg_surface_impl = C.struct_xdg_surface_interface{
	destroy:              (cPtr)(xdg_surface_destroy.CPtr()),
	set_transient_for:    (cPtr)(xdg_surface_set_transient_for.CPtr()),
	set_margin:           (cPtr)(xdg_surface_set_margin.CPtr()),
	set_title:            (cPtr)(xdg_surface_set_title.CPtr()),
	set_app_id:           (cPtr)(xdg_surface_set_app_id.CPtr()),
	move:                 (cPtr)(xdg_surface_move.CPtr()),
	resize:               (cPtr)(xdg_surface_resize.CPtr()),
	set_output:           (cPtr)(xdg_surface_set_output.CPtr()),
	request_change_state: (cPtr)(xdg_surface_request_change_state.CPtr()),
	ack_change_state:     (cPtr)(xdg_surface_ack_change_state.CPtr()),
	set_minimized:        (cPtr)(xdg_surface_set_minized.CPtr()),
}

var use_unstable_version = cfn.CreateFunc(
	func(client *C.struct_wl_client,
		resource *C.struct_wl_resource,
		version C.int32_t) {

		println("unstable")
	},
)

var get_xdg_surface = cfn.CreateFunc(
	func(client *C.struct_wl_client,
		resource *C.struct_wl_resource,
		id C.uint32_t,
		surface *C.struct_wl_resource) {

		surface_res := C.wl_resource_create(client, C.WL_xdg_surface_interface,
			C.wl_resource_get_version(resource), id)

		println("get_xdg_surface", int(id))

		C.wl_resource_set_implementation(surface_res,
			(unsafe.Pointer)(&xdg_surface_impl),
			nil,
			nil)
	},
)

var get_xdg_popup = cfn.CreateFunc(
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

var xdg_pong = cfn.CreateFunc(
	func(client *C.struct_wl_client,
		resource *C.struct_wl_resource,
		serial C.uint32_t) {

	},
)

var xdg_shell_impl = C.struct_xdg_shell_interface{
	use_unstable_version: (cPtr)(use_unstable_version.CPtr()),
	get_xdg_surface:      (cPtr)(get_xdg_surface.CPtr()),
	get_xdg_popup:        (cPtr)(get_xdg_popup.CPtr()),
	pong:                 (cPtr)(xdg_pong.CPtr()),
}

var bind_xdg_shell = cfn.CreateFunc(
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
		cPtr(bind_xdg_shell.CPtr()))
}
