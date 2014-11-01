package main

/*
#cgo pkg-config: wayland-server

#include <wayland-server.h>
*/
import "C"

import (
	"fmt"
)

type wl_client C.struct_wl_client
type wl_resource C.struct_wl_resource



var display *C.struct_wl_display
var event_loop *C.struct_wl_event_loop

func wayland() {
	display = C.wl_display_create()

	if display == nil {
		return
	}

	if C.wl_display_add_socket(display, nil) != 0 {
		return
	}

	event_loop = C.wl_display_get_event_loop(display)

	compositorInit(display)

	shmInit(display)

	seatInit(display)

	shellInit(display)

	xdgShellInit(display)


	fmt.Println("Wayland chrome")
	println("start running...")

	C.wl_display_run(display)

	C.wl_display_destroy(display)
}
