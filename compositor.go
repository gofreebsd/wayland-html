package main

/*
#cgo pkg-config: wayland-server

#include <wayland-server.h>
*/
import "C"

import (
	"unsafe"
)

var new_surface_signal C.struct_wl_signal

var surface_destroy = create_func(
	func(client *C.struct_wl_client,
		resource *C.struct_wl_resource) {
		println("surface_destroy")
		C.wl_resource_destroy(resource)
	},
)

var attach = create_func(
	func(client *C.struct_wl_client,
		resource *C.struct_wl_resource,
		buffer *C.struct_wl_resource,
		x C.int32_t,
		y C.int32_t) {
		C.wl_buffer_send_release(buffer)
	},
)

var damage = create_func(
	get_echo("damage"),
)

var frame_cbs *C.struct_wl_resource

var frame = create_func(
	func(client *C.struct_wl_client,
		resource *C.struct_wl_resource,
		id C.uint32_t) {

		println("frame")
		callback_resource := C.wl_resource_create(client,
			&C.wl_callback_interface,
			1, id)

		C.wl_resource_set_implementation(callback_resource,
			nil, nil, nil)

		frame_cbs = callback_resource

	},
)

var commit = create_func(
	func(client *C.struct_wl_client,
		resource *C.struct_wl_resource) {
		println("commit")
		C.wl_callback_send_done(frame_cbs, 1)
		C.wl_resource_destroy(frame_cbs)
	},
)

var surface_impl = C.struct_wl_surface_interface{
	destroy: cPtr(surface_destroy.fn_ptr),
	attach:  cPtr(attach.fn_ptr),
	damage:  cPtr(damage.fn_ptr),
	frame:   cPtr(frame.fn_ptr),
	commit:  cPtr(commit.fn_ptr),
}

var create_surface = create_func(
	func(client *C.struct_wl_client,
		resource *C.struct_wl_resource, id C.uint32_t) {
		var pid C.pid_t

		C.wl_client_get_credentials(client, &pid, nil, nil)

		surface_res := C.wl_resource_create(client, &C.wl_surface_interface,
			C.wl_resource_get_version(resource), id)

		C.wl_resource_set_implementation(surface_res,
			(unsafe.Pointer)(&surface_impl),
			nil, nil)

		println("surface", pid)
		C.wl_signal_emit(&new_surface_signal, (unsafe.Pointer)(surface_res))
	},
)

var create_region = create_func(
	func(client *C.struct_wl_client,
		resource *C.struct_wl_resource, id C.uint32_t) {
		println("region")
	},
)

var compositor_impl = C.struct_wl_compositor_interface{
	create_surface: (cPtr)(create_surface.fn_ptr),
	create_region:  (cPtr)(create_region.fn_ptr),
}

type Compositor struct {
}

var compositors = make(map[*C.struct_wl_client]*Compositor)

var compositor_cleans = make(map[*C.struct_wl_client]*CFn)

var bind_compositor = create_func(
	func(client *C.struct_wl_client, data unsafe.Pointer,
		version C.int, id C.uint32_t) {

		if version >= 3 {
			version = 3
		}

		resource := C.wl_resource_create(client, &C.wl_compositor_interface, version, id)

		compositors[client] = new(Compositor)
		compositor_cleans[client] = create_func(
			func(resource *C.struct_wl_resource) {
				delete(compositors, client)
				delete(compositor_cleans, client)
			},
		)

		C.wl_resource_set_implementation(resource,
			(unsafe.Pointer)(&compositor_impl),
			unsafe.Pointer(compositors[client]),
			cPtr(compositor_cleans[client].fn_ptr))
	},
)

func compositorInit(display *C.struct_wl_display) {

	C.wl_global_create(display,
		&C.wl_compositor_interface,
		3,
		nil,
		cPtr(bind_compositor.fn_ptr))

	C.wl_signal_init(&new_surface_signal)
}
