package main

/*
#cgo pkg-config: wayland-server

#include <wayland-server.h>
#include <sys/types.h>
#include <pwd.h>
#include <stdlib.h>
#include <stdio.h>

#include <ffi.h>


*/
import "C"

import (
	"os"
	// "reflect"
	"fmt"
)


var display *C.struct_wl_display
var event_loop *C.struct_wl_event_loop

func main() {

	display = C.wl_display_create()

	if display == nil {
		os.Exit(1)
	}

	if C.wl_display_add_socket(display, nil) != 0 {
		os.Exit(1)
	}


	compositorInit(display)

	shmInit(display)

	seatInit(display)

	shellInit(display)

	xdgShellInit(display)

    event_loop = C.wl_display_get_event_loop(display);

	fmt.Println("Wayland chrome")
	println("start running...")

	C.wl_display_run(display)

	C.wl_display_destroy(display)

	os.Exit(0)

}
