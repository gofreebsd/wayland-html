package main

/*
#cgo pkg-config: wayland-server

#include <wayland-server.h>
*/
import "C"

import (
	"bytes"
	_ "fmt"
	"github.com/fangyuanziti/wayland-html/cfn"
	_ "time"
	"unsafe"
)

type Surface struct {
	pendingBuffer *C.struct_wl_resource
	frame_cbs     []*C.struct_wl_resource
	CommitBuffer  []byte
}

type Compositor struct {
	Pid      int
	Surfaces map[*C.struct_wl_resource]*Surface
}

var new_surface_signal C.struct_wl_signal

var surface_destroy = cfn.CreateFunc(
	func(client *C.struct_wl_client,
		resource *C.struct_wl_resource) {
		println("surface_destroy")
		C.wl_resource_destroy(resource)
	},
)

var attach = cfn.CreateFunc(
	func(client *C.struct_wl_client,
		resource *C.struct_wl_resource,
		buffer *C.struct_wl_resource,
		x C.int32_t,
		y C.int32_t) {

		println("attach")

		surface := compositors[client].Surfaces[resource]
		surface.pendingBuffer = buffer
	},
)

var damage = cfn.CreateFunc(
	get_echo("damage"),
)

var frame = cfn.CreateFunc(
	func(client *C.struct_wl_client,
		resource *C.struct_wl_resource,
		id C.uint32_t) {

		println("frame")

		callback_resource := C.wl_resource_create(client,
			&C.wl_callback_interface,
			1, id)

		C.wl_resource_set_implementation(callback_resource,
			nil, nil, nil)

		surface := compositors[client].Surfaces[resource]

		surface.frame_cbs = append(surface.frame_cbs, callback_resource)

	},
)

func argbToRgba(buf []byte) []byte {

	buffer := bytes.NewBuffer(buf)
	new_buffer := new(bytes.Buffer)

	read_buffer := make([]byte, 4)
	convert_buffer := make([]byte, 4)

	for {
		_, err := buffer.Read(read_buffer)

		if err != nil {
			break
		}

		convert_buffer[0] = read_buffer[2]
		convert_buffer[1] = read_buffer[1]
		convert_buffer[2] = read_buffer[0]
		convert_buffer[3] = read_buffer[3]

		new_buffer.Write(convert_buffer)

	}

	return new_buffer.Bytes()

}

var commit = cfn.CreateFunc(
	func(client *C.struct_wl_client,
		resource *C.struct_wl_resource) {

		compositor := compositors[client]
		surface := compositor.Surfaces[resource]

		// use pending buffer
		println(
			buffers[surface.pendingBuffer].offset,
			buffers[surface.pendingBuffer].width,
			buffers[surface.pendingBuffer].height,
			buffers[surface.pendingBuffer].stride,
			buffers[surface.pendingBuffer].format,
		)

		pendingBuffer := buffers[surface.pendingBuffer]
		size := pendingBuffer.height * pendingBuffer.stride

		var commitBuffer = make([]byte, size)
		copy(commitBuffer, pendingBuffer.pool.ptr[pendingBuffer.offset:])
		surface.CommitBuffer = argbToRgba(commitBuffer)

		// release pending buffer
		C.wl_buffer_send_release(surface.pendingBuffer)

	},
)

var surface_impl = C.struct_wl_surface_interface{
	destroy: cPtr(surface_destroy.CPtr()),
	attach:  cPtr(attach.CPtr()),
	damage:  cPtr(damage.CPtr()),
	frame:   cPtr(frame.CPtr()),
	commit:  cPtr(commit.CPtr()),
}

var create_surface = cfn.CreateFunc(
	func(client *C.struct_wl_client,
		resource *C.struct_wl_resource, id C.uint32_t) {

		surface_res := C.wl_resource_create(client, &C.wl_surface_interface,
			C.wl_resource_get_version(resource), id)

		compositor := compositors[client]

		surface := new(Surface)

		compositor.Surfaces[surface_res] = surface

		C.wl_resource_set_implementation(surface_res,
			(unsafe.Pointer)(&surface_impl),
			(unsafe.Pointer)(surface), nil)

		C.wl_signal_emit(&new_surface_signal, (unsafe.Pointer)(surface_res))
	},
)

var create_region = cfn.CreateFunc(
	func(client *C.struct_wl_client,
		resource *C.struct_wl_resource, id C.uint32_t) {
		println("region")
	},
)

var compositor_impl = C.struct_wl_compositor_interface{
	create_surface: (cPtr)(create_surface.CPtr()),
	create_region:  (cPtr)(create_region.CPtr()),
}

func (c *Compositor) resetFrameCallback() {
	for _, s := range c.Surfaces {
		for _, cb := range s.frame_cbs {
			C.wl_callback_send_done(cb, 1)
			C.wl_resource_destroy(cb)
		}

		s.frame_cbs = []*C.struct_wl_resource{}
	}
}

var compositors = make(map[*C.struct_wl_client]*Compositor)

var bind_compositor = cfn.CreateFunc(
	func(client *C.struct_wl_client, data unsafe.Pointer,
		version C.int, id C.uint32_t) {

		if version >= 3 {
			version = 3
		}

		resource := C.wl_resource_create(client, &C.wl_compositor_interface, version, id)

		timer := NewRepeatTimer()

		compositors[client] = new(Compositor)

		destroy := once_func(
			func(resource *C.struct_wl_resource) {
				delete(compositors, client)
				// delete(compositor_cleans, client)
				timer.stop()
			},
		)

		var pid C.pid_t
		C.wl_client_get_credentials(client, &pid, nil, nil)
		compositors[client].Pid = int(pid)
		compositors[client].Surfaces = make(map[*C.struct_wl_resource]*Surface)

		C.wl_resource_set_implementation(resource,
			unsafe.Pointer(&compositor_impl),
			unsafe.Pointer(compositors[client]),
			cPtr(destroy.CPtr()))

		timer.start(compositors[client])
	},
)

type RepeatTimer struct {
	timer   *C.struct_wl_event_source
	is_stop bool
}

func NewRepeatTimer() *RepeatTimer {
	ret := new(RepeatTimer)
	ret.is_stop = false
	return ret
}

func (t *RepeatTimer) stop() {
	t.is_stop = true
}

func (t *RepeatTimer) start(compositor *Compositor) {
	timer_tick := once_func(func() {
		if !t.is_stop {
			println("timer", compositor)
			compositor.resetFrameCallback()
			t.start(compositor)
		}
	})

	t.timer = C.wl_event_loop_add_timer(event_loop, cPtr(timer_tick.CPtr()), nil)
	C.wl_event_source_timer_update(t.timer, 3*1000)
}

func compositorInit(display *C.struct_wl_display) {

	C.wl_global_create(display,
		&C.wl_compositor_interface,
		3,
		nil,
		cPtr(bind_compositor.CPtr()))

	C.wl_signal_init(&new_surface_signal)
}
