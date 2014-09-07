package main

/*
#cgo pkg-config: wayland-server

#include <wayland-server.h>

*/
import "C"

import (
	"unsafe"
)

var seat_impl C.struct_wl_seat_interface = C.struct_wl_seat_interface{
	get_pointer:  nil,
	get_keyboard: nil,
	get_touch:    nil,
}

func bind_seat(client *C.struct_wl_client, data unsafe.Pointer,
	version C.int, id C.uint32_t) {

	if version >= 2 {
		version = 2
	}

	resource := C.wl_resource_create(client, &C.wl_seat_interface, version, id)

	C.wl_resource_set_implementation(
		resource,
		(unsafe.Pointer)(&seat_impl),
		nil,
		nil)

	if version >= 2 {
		C.wl_seat_send_name(resource, C.CString("default"))
	}

	C.wl_seat_send_capabilities(resource, 0)

}

var seat *CFn

func seatInit(display *C.struct_wl_display) {

	seat = create_func(bind_seat)

	C.wl_global_create(display,
		&C.wl_seat_interface,
		3,
		nil,
		cPtr(seat.fn_ptr))
}
