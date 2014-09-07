package main

/*
#cgo pkg-config: wayland-server

#include <wayland-server.h>
*/
import "C"

type wl_client C.struct_wl_client
type wl_resource C.struct_wl_resource
