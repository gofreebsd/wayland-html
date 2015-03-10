package main

/*
#include "wayland-fix.h"

*/
import "C"

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/fangyuanziti/wayland-html/cfn"
	"image"
	"image/color"
	"syscall"
	"unsafe"
)

type Pool struct {
	client    *C.struct_wl_client
	id        int
	destroyFn *cfn.CFn
	fd        int
	ptr       []byte
}

func (pool *Pool) get_destroy_func() unsafe.Pointer {
	if pool.destroyFn == nil {
		pool.destroyFn = cfn.CreateFunc(
			func(resource *C.struct_wl_resource) {
				id := C.wl_resource_get_id(resource)
				delete(pools, int(id))
			},
		)
	}

	return pool.destroyFn.CPtr()
}

var pools = make(map[int]*Pool)

var shm_create_pool = cfn.CreateFunc(func(
	client *C.struct_wl_client,
	resource *C.struct_wl_resource,
	id C.uint32_t,
	fd C.int32_t,
	size C.int32_t) {

	println("creata_pool")

	mmap_ptr, err := syscall.Mmap(
		int(fd), 0,
		int(size),
		syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_SHARED,
	)

	if err != nil {
		fmt.Println(err)
		// TODO: should send error
	}

	pool_res := C.wl_resource_create(
		client,
		&C.wl_shm_pool_interface,
		(C.wl_resource_get_version(resource)),
		id,
	)

	pool := Pool{
		client: client,
		id:     int(id),
		fd:     int(fd),
		ptr:    mmap_ptr,
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

})

type WLBufferImage struct {
	buffer *Buffer
}

func newBufferImage(buffer *Buffer) *WLBufferImage {
	img := &WLBufferImage{
		buffer: buffer,
	}
	log.Info(buffer.format, buffer.stride)
	return img
}

func (img *WLBufferImage) ColorModel() color.Model {
	// hard code for rgba
	return color.RGBAModel
}

func (img *WLBufferImage) Bounds() image.Rectangle {
	min := image.Point{
		X: 0,
		Y: 0,
	}
	max := image.Point{
		X: img.buffer.width,
		Y: img.buffer.height,
	}

	rec := image.Rectangle{
		Min: min,
		Max: max,
	}
	return rec

}

func (img *WLBufferImage) At(x, y int) color.Color {
	pointStride := img.buffer.stride / 4
	pos := axisToPos(x, y, pointStride)
	// c := img.colors[pos]
	bufferOffset := pos * 4
	ptr := img.buffer.pool.ptr[bufferOffset : bufferOffset+4]
	c := color.RGBA{
		R: ptr[2],
		G: ptr[1],
		B: ptr[0],
		A: 255, // hard code for 255 (no transparent)
	}
	return c
}

type Buffer struct {
	offset int
	width  int
	height int
	stride int
	format uint
	pool   *Pool
}

func (b *Buffer) Image() *WLBufferImage {
	img := newBufferImage(b)
	return img
}

var buffers = make(map[*C.struct_wl_resource]*Buffer)

var shm_pool_create_buffer = cfn.CreateFunc(func(
	client *C.struct_wl_client,
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
})

var shm_pool_destroy = cfn.CreateFunc(func(
	client *C.struct_wl_client,
	resource *C.struct_wl_resource) {
	println("destroy")
	// TODO: implement
	C.wl_resource_destroy(resource)
})

var shm_pool_resize = cfn.CreateFunc(func(
	client *C.struct_wl_client,
	resource *C.struct_wl_resource,
	size C.int32_t) {

	pool := (*Pool)(C.wl_resource_get_user_data(resource))

	// ummap old mapped memory
	syscall.Munmap(pool.ptr)

	// map fd with new size
	mmap_ptr, err := syscall.Mmap(
		int(pool.fd), 0,
		int(size),
		syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_SHARED,
	)

	if err != nil {
		fmt.Println(err)
	}

	// keep the reference of memory ptr
	pool.ptr = mmap_ptr
	println("resize")
})

var buffer_destroy = cfn.CreateFunc(func(
	client *C.struct_wl_client,
	resource *C.struct_wl_resource) {
	println("buffer_destroy")
	// TODO: implement
	C.wl_resource_destroy(resource)
})

var shm_impl = C.struct_wl_shm_interface{
	create_pool: (cPtr)(shm_create_pool.CPtr()),
}

var shm_pool_impl = C.struct_wl_shm_pool_interface{
	create_buffer: (cPtr)(shm_pool_create_buffer.CPtr()),
	destroy:       (cPtr)(shm_pool_destroy.CPtr()),
	resize:        (cPtr)(shm_pool_resize.CPtr()),
}

var buffer_impl = C.struct_wl_buffer_interface{
	destroy: (cPtr)(buffer_destroy.CPtr()),
}

func fixVersion(version C.int, implVersion C.int) C.int {
	if version >= implVersion {
		return implVersion
	} else {
		return version
	}
}

var bind_shm = cfn.CreateFunc(func(
	client *C.struct_wl_client,
	data unsafe.Pointer,
	version C.int, id C.uint32_t) {

	version = fixVersion(version, 1)

	println("shm", int(id))

	resource := C.wl_resource_create(client, &C.wl_shm_interface, version, id)
	C.wl_resource_set_implementation(resource,
		(unsafe.Pointer)(&shm_impl),
		nil,
		nil)

	C.wl_shm_send_format(resource, C.WL_SHM_FORMAT_RGBA8888)
	// C.wl_shm_send_format(resource, C.WL_SHM_FORMAT_XRGB8888)
	// C.wl_shm_send_format(resource, C.WL_SHM_FORMAT_ARGB8888)

})

func shmInit(display *C.struct_wl_display) {

	C.wl_global_create(display,
		&C.wl_shm_interface,
		1,
		nil,
		cPtr(bind_shm.CPtr()))
}
