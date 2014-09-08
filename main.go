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
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"strconv"
)

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

	compositorInit(display)

	shmInit(display)

	seatInit(display)

	shellInit(display)

	xdgShellInit(display)

	event_loop = C.wl_display_get_event_loop(display)

	fmt.Println("Wayland chrome")
	println("start running...")

	C.wl_display_run(display)

	C.wl_display_destroy(display)
}

func main() {

	go wayland()

	server := martini.Classic()

	server.Use(render.Renderer())

	server.Get("/", func(r render.Render) {
		r.HTML(200, "index", strconv.Itoa(len(compositors)))
	})

	server.Run()
}
