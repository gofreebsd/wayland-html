package main

/*
#cgo pkg-config: wayland-server

#include <wayland-server.h>

*/
import "C"

import (
	"fmt"
	"syscall"
	"unsafe"
)

type Pool struct {
	client    *C.struct_wl_client
	id        int
	destroyFn *CFn
	fd        int
	ptr       []byte
}

func (pool *Pool) get_destroy_func() unsafe.Pointer {
	if pool.destroyFn == nil {
		pool.destroyFn = create_func(
			func(resource *C.struct_wl_resource) {
				id := C.wl_resource_get_id(resource)
				delete(pools, int(id))
			},
		)
	}

	return pool.destroyFn.fn_ptr
}

var pools = make(map[int]*Pool)

var create_pool = create_func(
	func(client *C.struct_wl_client,
		resource *C.struct_wl_resource,
		id C.uint32_t,
		fd C.int32_t,
		size C.int32_t) {

		println("creata_pool")

		pool_res := C.wl_resource_create(
			client,
			&C.wl_shm_pool_interface,
			(C.wl_resource_get_version(resource)),
			id,
		)

		mmap_ptr, err := syscall.Mmap(
			int(fd), 0,
			int(size),
			syscall.PROT_READ|syscall.PROT_WRITE,
			syscall.MAP_SHARED,
		)

		if err != nil {
			fmt.Println(err)
		}

		pool := Pool{client: client, id: int(id),
			fd: int(fd), ptr: mmap_ptr,
		}

		pools[int(id)] = &pool

		// C.wl_resource_set_implementation(pool_res,
		// 	(unsafe.Pointer)(&shm_pool_impl),
		// 	unsafe.Pointer(&pool),
		// 	cPtr(pool.get_destroy_func()))
		C.wl_resource_set_implementation(pool_res,
			(unsafe.Pointer)(&shm_pool_impl),
			unsafe.Pointer(&pool),
			nil)

	},
)

type Buffer struct {
	offset int
	width  int
	height int
	stride int
	format uint
	pool   *Pool
}

var buffers = make(map[*C.struct_wl_resource]*Buffer)

var create_buffer = create_func(
	func(client *C.struct_wl_client,
		resource *C.struct_wl_resource,
		id C.uint32_t,
		offset C.int32_t,
		width C.int32_t,
		height C.int32_t,
		stride C.int32_t,
		format C.uint32_t) {

		buffer := C.wl_resource_create(
			client,
			&C.wl_buffer_interface,
			(C.wl_resource_get_version(resource)),
			id,
		)

		pool := (*Pool)(C.wl_resource_get_user_data(resource))

		buffer_data := Buffer{
			offset: int(offset),
			width:  int(width),
			height: int(height),
			stride: int(stride),
			format: uint(offset),
			pool:   pool,
		}

		buffers[buffer] = &buffer_data

		C.wl_resource_set_implementation(buffer,
			(unsafe.Pointer)(&buffer_impl),
			unsafe.Pointer(&buffer_data),
			nil)
	},
)

var destroy = create_func(
	func(client *C.struct_wl_client,
		resource *C.struct_wl_resource) {
		println("destroy")
		C.wl_resource_destroy(resource)
	},
)

var resize = create_func(
	func(client *C.struct_wl_client,
		resource *C.struct_wl_resource,
		size C.int32_t) {
		println("resize")
	},
)

var buffer_destroy = create_func(
	func(client *C.struct_wl_client,
		resource *C.struct_wl_resource) {
		println("buffer_destroy")
		C.wl_resource_destroy(resource)
	},
)

var shm_impl = C.struct_wl_shm_interface{
	create_pool: (cPtr)(create_pool.fn_ptr),
}

var shm_pool_impl = C.struct_wl_shm_pool_interface{
	create_buffer: (cPtr)(create_buffer.fn_ptr),
	destroy:       (cPtr)(destroy.fn_ptr),
	resize:        (cPtr)(resize.fn_ptr),
}

var buffer_impl = C.struct_wl_buffer_interface{
	destroy: (cPtr)(buffer_destroy.fn_ptr),
}

var bind_shm = create_func(

	func(client *C.struct_wl_client, data unsafe.Pointer,
		version C.int, id C.uint32_t) {

		if version >= 1 {
			version = 1
		}

		println("shm", int(id))

		resource := C.wl_resource_create(client, &C.wl_shm_interface, version, id)
		C.wl_resource_set_implementation(resource,
			(unsafe.Pointer)(&shm_impl),
			nil,
			nil)

		C.wl_shm_send_format(resource, C.WL_SHM_FORMAT_RGBA8888)
		C.wl_shm_send_format(resource, C.WL_SHM_FORMAT_XRGB8888)
		C.wl_shm_send_format(resource, C.WL_SHM_FORMAT_ARGB8888)

	},
)

func shmInit(display *C.struct_wl_display) {

	C.wl_global_create(display,
		&C.wl_shm_interface,
		1,
		nil,
		cPtr(bind_shm.fn_ptr))
}
