package main

/*
#cgo pkg-config: wayland-server

#include <wayland-server.h>

*/
// import "C"

// import (
// 	"unsafe"
// )



// var xdg_shell_impl = C.struct_xdg_shell_interface {
// } 

// func bind_xdg_shell(client *C.struct_wl_client, data unsafe.Pointer,
//     version C.int , id C.uint32_t) {

// 	if version >= 1 {
// 		version = 1
// 	}

// 	resource := C.wl_resource_create(client, &C.xdg_shell_interface, version, id)

//     C.wl_resource_set_implementation(
// 		resource,
// 		(unsafe.Pointer)(&shell_impl),
//         nil,
// 		nil)
// }

// var xdg_shell = create_func(bind_xdg_shell)

// func xdgInit(display  *C.struct_wl_display) {

// 	C.wl_global_create(display,
// 		&C.xdg_shell_interface,
// 		1,
// 		nil,
// 		cPtr(xdg_shell.fn_ptr))
// }
