/* 
 * Copyright © 2008-2013 Kristian Høgsberg
 * Copyright © 2013      Rafael Antognolli
 * Copyright © 2013      Jasper St. Pierre
 * Copyright © 2010-2013 Intel Corporation
 * 
 * Permission to use, copy, modify, distribute, and sell this
 * software and its documentation for any purpose is hereby granted
 * without fee, provided that the above copyright notice appear in
 * all copies and that both that copyright notice and this permission
 * notice appear in supporting documentation, and that the name of
 * the copyright holders not be used in advertising or publicity
 * pertaining to distribution of the software without specific,
 * written prior permission.  The copyright holders make no
 * representations about the suitability of this software for any
 * purpose.  It is provided "as is" without express or implied
 * warranty.
 * 
 * THE COPYRIGHT HOLDERS DISCLAIM ALL WARRANTIES WITH REGARD TO THIS
 * SOFTWARE, INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND
 * FITNESS, IN NO EVENT SHALL THE COPYRIGHT HOLDERS BE LIABLE FOR ANY
 * SPECIAL, INDIRECT OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN
 * AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION,
 * ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF
 * THIS SOFTWARE.
 */

#ifndef XDG_SHELL_SERVER_PROTOCOL_H
#define XDG_SHELL_SERVER_PROTOCOL_H

#ifdef  __cplusplus
extern "C" {
#endif

#include <stdint.h>
#include <stddef.h>
#include "wayland-util.h"

struct wl_client;
struct wl_resource;

struct xdg_shell;
struct xdg_surface;
struct xdg_popup;

extern const struct wl_interface xdg_shell_interface;
extern const struct wl_interface xdg_surface_interface;
extern const struct wl_interface xdg_popup_interface;

#ifndef XDG_SHELL_VERSION_ENUM
#define XDG_SHELL_VERSION_ENUM
/**
 * xdg_shell_version - latest protocol version
 * @XDG_SHELL_VERSION_CURRENT: Always the latest version
 *
 * The 'current' member of this enum gives the version of the protocol.
 * Implementations can compare this to the version they implement using
 * static_assert to ensure the protocol and implementation versions match.
 */
enum xdg_shell_version {
	XDG_SHELL_VERSION_CURRENT = 3,
};
#endif /* XDG_SHELL_VERSION_ENUM */

/**
 * xdg_shell - create desktop-style surfaces
 * @use_unstable_version: enable use of this unstable version
 * @get_xdg_surface: create a shell surface from a surface
 * @get_xdg_popup: create a shell surface from a surface
 * @pong: respond to a ping event
 *
 * This interface is implemented by servers that provide desktop-style
 * user interfaces.
 *
 * It allows clients to associate a xdg_surface with a basic surface.
 */
struct xdg_shell_interface {
	/**
	 * use_unstable_version - enable use of this unstable version
	 * @version: (none)
	 *
	 * Negotiate the unstable version of the interface. This
	 * mechanism is in place to ensure client and server agree on the
	 * unstable versions of the protocol that they speak or exit
	 * cleanly if they don't agree. This request will go away once the
	 * xdg-shell protocol is stable.
	 */
	void (*use_unstable_version)(struct wl_client *client,
				     struct wl_resource *resource,
				     int32_t version);
	/**
	 * get_xdg_surface - create a shell surface from a surface
	 * @id: (none)
	 * @surface: (none)
	 *
	 * Create a shell surface for an existing surface.
	 *
	 * Only one shell or popup surface can be associated with a given
	 * surface.
	 */
	void (*get_xdg_surface)(struct wl_client *client,
				struct wl_resource *resource,
				uint32_t id,
				struct wl_resource *surface);
	/**
	 * get_xdg_popup - create a shell surface from a surface
	 * @id: (none)
	 * @surface: (none)
	 * @parent: (none)
	 * @seat: the wl_seat whose pointer is used
	 * @serial: serial of the implicit grab on the pointer
	 * @x: (none)
	 * @y: (none)
	 * @flags: (none)
	 *
	 * Create a popup surface for an existing surface.
	 *
	 * Only one shell or popup surface can be associated with a given
	 * surface.
	 */
	void (*get_xdg_popup)(struct wl_client *client,
			      struct wl_resource *resource,
			      uint32_t id,
			      struct wl_resource *surface,
			      struct wl_resource *parent,
			      struct wl_resource *seat,
			      uint32_t serial,
			      int32_t x,
			      int32_t y,
			      uint32_t flags);
	/**
	 * pong - respond to a ping event
	 * @serial: serial of the ping event
	 *
	 * A client must respond to a ping event with a pong request or
	 * the client may be deemed unresponsive.
	 */
	void (*pong)(struct wl_client *client,
		     struct wl_resource *resource,
		     uint32_t serial);
};

#define XDG_SHELL_PING	0

#define XDG_SHELL_PING_SINCE_VERSION	1

static inline void
xdg_shell_send_ping(struct wl_resource *resource_, uint32_t serial)
{
	wl_resource_post_event(resource_, XDG_SHELL_PING, serial);
}

#ifndef XDG_SURFACE_RESIZE_EDGE_ENUM
#define XDG_SURFACE_RESIZE_EDGE_ENUM
/**
 * xdg_surface_resize_edge - edge values for resizing
 * @XDG_SURFACE_RESIZE_EDGE_NONE: (none)
 * @XDG_SURFACE_RESIZE_EDGE_TOP: (none)
 * @XDG_SURFACE_RESIZE_EDGE_BOTTOM: (none)
 * @XDG_SURFACE_RESIZE_EDGE_LEFT: (none)
 * @XDG_SURFACE_RESIZE_EDGE_TOP_LEFT: (none)
 * @XDG_SURFACE_RESIZE_EDGE_BOTTOM_LEFT: (none)
 * @XDG_SURFACE_RESIZE_EDGE_RIGHT: (none)
 * @XDG_SURFACE_RESIZE_EDGE_TOP_RIGHT: (none)
 * @XDG_SURFACE_RESIZE_EDGE_BOTTOM_RIGHT: (none)
 *
 * These values are used to indicate which edge of a surface is being
 * dragged in a resize operation. The server may use this information to
 * adapt its behavior, e.g. choose an appropriate cursor image.
 */
enum xdg_surface_resize_edge {
	XDG_SURFACE_RESIZE_EDGE_NONE = 0,
	XDG_SURFACE_RESIZE_EDGE_TOP = 1,
	XDG_SURFACE_RESIZE_EDGE_BOTTOM = 2,
	XDG_SURFACE_RESIZE_EDGE_LEFT = 4,
	XDG_SURFACE_RESIZE_EDGE_TOP_LEFT = 5,
	XDG_SURFACE_RESIZE_EDGE_BOTTOM_LEFT = 6,
	XDG_SURFACE_RESIZE_EDGE_RIGHT = 8,
	XDG_SURFACE_RESIZE_EDGE_TOP_RIGHT = 9,
	XDG_SURFACE_RESIZE_EDGE_BOTTOM_RIGHT = 10,
};
#endif /* XDG_SURFACE_RESIZE_EDGE_ENUM */

#ifndef XDG_SURFACE_STATE_ENUM
#define XDG_SURFACE_STATE_ENUM
/**
 * xdg_surface_state - types of state on the surface
 * @XDG_SURFACE_STATE_MAXIMIZED: the surface is maximized
 * @XDG_SURFACE_STATE_FULLSCREEN: the surface is fullscreen
 *
 * The different state values used on the surface. This is designed for
 * state values like maximized, fullscreen. It is paired with the
 * request_change_state event to ensure that both the client and the
 * compositor setting the state can be synchronized.
 *
 * States set in this way are double-buffered. They will get applied on the
 * next commit.
 *
 * Desktop environments may extend this enum by taking up a range of values
 * and documenting the range they chose in this description. They are not
 * required to document the values for the range that they chose. Ideally,
 * any good extensions from a desktop environment should make its way into
 * standardization into this enum.
 *
 * The current reserved ranges are:
 *
 * 0x0000 - 0x0FFF: xdg-shell core values, documented below. 0x1000 -
 * 0x1FFF: GNOME
 */
enum xdg_surface_state {
	XDG_SURFACE_STATE_MAXIMIZED = 1,
	XDG_SURFACE_STATE_FULLSCREEN = 2,
};
#endif /* XDG_SURFACE_STATE_ENUM */

/**
 * xdg_surface - desktop-style metadata interface
 * @destroy: remove xdg_surface interface
 * @set_transient_for: surface is a child of another surface
 * @set_margin: set the visible frame boundaries
 * @set_title: set surface title
 * @set_app_id: set surface class
 * @move: start an interactive move
 * @resize: start an interactive resize
 * @set_output: set the default output used by this surface
 * @request_change_state: client requests to change a surface's state
 * @ack_change_state: ack a change_state event
 * @set_minimized: minimize the surface
 *
 * An interface that may be implemented by a wl_surface, for
 * implementations that provide a desktop-style user interface.
 *
 * It provides requests to treat surfaces like windows, allowing to set
 * properties like maximized, fullscreen, minimized, and to move and resize
 * them, and associate metadata like title and app id.
 *
 * On the server side the object is automatically destroyed when the
 * related wl_surface is destroyed. On client side, xdg_surface.destroy()
 * must be called before destroying the wl_surface object.
 */
struct xdg_surface_interface {
	/**
	 * destroy - remove xdg_surface interface
	 *
	 * The xdg_surface interface is removed from the wl_surface
	 * object that was turned into a xdg_surface with
	 * xdg_shell.get_xdg_surface request. The xdg_surface properties,
	 * like maximized and fullscreen, are lost. The wl_surface loses
	 * its role as a xdg_surface. The wl_surface is unmapped.
	 */
	void (*destroy)(struct wl_client *client,
			struct wl_resource *resource);
	/**
	 * set_transient_for - surface is a child of another surface
	 * @parent: (none)
	 *
	 * Setting a surface as transient of another means that it is
	 * child of another surface.
	 *
	 * Child surfaces are stacked above their parents, and will be
	 * unmapped if the parent is unmapped too. They should not appear
	 * on task bars and alt+tab.
	 */
	void (*set_transient_for)(struct wl_client *client,
				  struct wl_resource *resource,
				  struct wl_resource *parent);
	/**
	 * set_margin - set the visible frame boundaries
	 * @left_margin: (none)
	 * @right_margin: (none)
	 * @top_margin: (none)
	 * @bottom_margin: (none)
	 *
	 * This tells the compositor what the visible size of the window
	 * should be, so it can use it to determine what borders to use for
	 * constrainment and alignment.
	 *
	 * CSD often has invisible areas for decoration purposes, like drop
	 * shadows. These "shadow" drawings need to be subtracted out of
	 * the normal boundaries of the window when computing where to
	 * place windows (e.g. to set this window so it's centered on top
	 * of another, or to put it to the left or right of the screen.)
	 *
	 * This value should change as little as possible at runtime, to
	 * prevent flicker.
	 *
	 * This value is also ignored when the window is maximized or
	 * fullscreen, and assumed to be 0.
	 *
	 * If never called, this value is assumed to be 0.
	 */
	void (*set_margin)(struct wl_client *client,
			   struct wl_resource *resource,
			   int32_t left_margin,
			   int32_t right_margin,
			   int32_t top_margin,
			   int32_t bottom_margin);
	/**
	 * set_title - set surface title
	 * @title: (none)
	 *
	 * Set a short title for the surface.
	 *
	 * This string may be used to identify the surface in a task bar,
	 * window list, or other user interface elements provided by the
	 * compositor.
	 *
	 * The string must be encoded in UTF-8.
	 */
	void (*set_title)(struct wl_client *client,
			  struct wl_resource *resource,
			  const char *title);
	/**
	 * set_app_id - set surface class
	 * @app_id: (none)
	 *
	 * Set an id for the surface.
	 *
	 * The app id identifies the general class of applications to which
	 * the surface belongs.
	 *
	 * It should be the ID that appears in the new desktop entry
	 * specification, the interface name.
	 */
	void (*set_app_id)(struct wl_client *client,
			   struct wl_resource *resource,
			   const char *app_id);
	/**
	 * move - start an interactive move
	 * @seat: the wl_seat whose pointer is used
	 * @serial: serial of the implicit grab on the pointer
	 *
	 * Start a pointer-driven move of the surface.
	 *
	 * This request must be used in response to a button press event.
	 * The server may ignore move requests depending on the state of
	 * the surface (e.g. fullscreen or maximized).
	 */
	void (*move)(struct wl_client *client,
		     struct wl_resource *resource,
		     struct wl_resource *seat,
		     uint32_t serial);
	/**
	 * resize - start an interactive resize
	 * @seat: the wl_seat whose pointer is used
	 * @serial: serial of the implicit grab on the pointer
	 * @edges: which edge or corner is being dragged
	 *
	 * Start a pointer-driven resizing of the surface.
	 *
	 * This request must be used in response to a button press event.
	 * The server may ignore resize requests depending on the state of
	 * the surface (e.g. fullscreen or maximized).
	 */
	void (*resize)(struct wl_client *client,
		       struct wl_resource *resource,
		       struct wl_resource *seat,
		       uint32_t serial,
		       uint32_t edges);
	/**
	 * set_output - set the default output used by this surface
	 * @output: (none)
	 *
	 * Set the default output used by this surface when it is first
	 * mapped.
	 *
	 * If this value is NULL (default), it's up to the compositor to
	 * choose which display will be used to map this surface.
	 *
	 * When fullscreen or maximized state are set on this surface, and
	 * it wasn't mapped yet, the output set with this method will be
	 * used. Otherwise, the output where the surface is currently
	 * mapped will be used.
	 */
	void (*set_output)(struct wl_client *client,
			   struct wl_resource *resource,
			   struct wl_resource *output);
	/**
	 * request_change_state - client requests to change a surface's
	 *	state
	 * @state_type: the state to set
	 * @value: the value to change the state to
	 * @serial: an event serial
	 *
	 * This asks the compositor to change the state. If the
	 * compositor wants to change the state, it will send a
	 * change_state event with the same state_type, value, and serial,
	 * and the event flow continues as if it it was initiated by the
	 * compositor.
	 *
	 * If the compositor does not want to change the state, it will
	 * send a change_state to the client with the old value of the
	 * state.
	 */
	void (*request_change_state)(struct wl_client *client,
				     struct wl_resource *resource,
				     uint32_t state_type,
				     uint32_t value,
				     uint32_t serial);
	/**
	 * ack_change_state - ack a change_state event
	 * @state_type: the state to set
	 * @value: the value to change the state to
	 * @serial: a serial to pass to change_state
	 *
	 * When a change_state event is received, a client should then
	 * ack it using the ack_change_state request to ensure that the
	 * compositor knows the client has seen the event.
	 *
	 * By this point, the state is confirmed, and the next attach
	 * should contain the buffer drawn for the new state value.
	 *
	 * The values here need to be the same as the values in the
	 * cooresponding change_state event.
	 */
	void (*ack_change_state)(struct wl_client *client,
				 struct wl_resource *resource,
				 uint32_t state_type,
				 uint32_t value,
				 uint32_t serial);
	/**
	 * set_minimized - minimize the surface
	 *
	 * Minimize the surface.
	 */
	void (*set_minimized)(struct wl_client *client,
			      struct wl_resource *resource);
};

#define XDG_SURFACE_CONFIGURE	0
#define XDG_SURFACE_CHANGE_STATE	1
#define XDG_SURFACE_ACTIVATED	2
#define XDG_SURFACE_DEACTIVATED	3
#define XDG_SURFACE_CLOSE	4

#define XDG_SURFACE_CONFIGURE_SINCE_VERSION	1
#define XDG_SURFACE_CHANGE_STATE_SINCE_VERSION	1
#define XDG_SURFACE_ACTIVATED_SINCE_VERSION	1
#define XDG_SURFACE_DEACTIVATED_SINCE_VERSION	1
#define XDG_SURFACE_CLOSE_SINCE_VERSION	1

static inline void
xdg_surface_send_configure(struct wl_resource *resource_, int32_t width, int32_t height)
{
	wl_resource_post_event(resource_, XDG_SURFACE_CONFIGURE, width, height);
}

static inline void
xdg_surface_send_change_state(struct wl_resource *resource_, uint32_t state_type, uint32_t value, uint32_t serial)
{
	wl_resource_post_event(resource_, XDG_SURFACE_CHANGE_STATE, state_type, value, serial);
}

static inline void
xdg_surface_send_activated(struct wl_resource *resource_)
{
	wl_resource_post_event(resource_, XDG_SURFACE_ACTIVATED);
}

static inline void
xdg_surface_send_deactivated(struct wl_resource *resource_)
{
	wl_resource_post_event(resource_, XDG_SURFACE_DEACTIVATED);
}

static inline void
xdg_surface_send_close(struct wl_resource *resource_)
{
	wl_resource_post_event(resource_, XDG_SURFACE_CLOSE);
}

/**
 * xdg_popup - desktop-style metadata interface
 * @destroy: remove xdg_surface interface
 *
 * An interface that may be implemented by a wl_surface, for
 * implementations that provide a desktop-style popups/menus. A popup
 * surface is a transient surface with an added pointer grab.
 *
 * An existing implicit grab will be changed to owner-events mode, and the
 * popup grab will continue after the implicit grab ends (i.e. releasing
 * the mouse button does not cause the popup to be unmapped).
 *
 * The popup grab continues until the window is destroyed or a mouse button
 * is pressed in any other clients window. A click in any of the clients
 * surfaces is reported as normal, however, clicks in other clients
 * surfaces will be discarded and trigger the callback.
 *
 * The x and y arguments specify the locations of the upper left corner of
 * the surface relative to the upper left corner of the parent surface, in
 * surface local coordinates.
 *
 * xdg_popup surfaces are always transient for another surface.
 */
struct xdg_popup_interface {
	/**
	 * destroy - remove xdg_surface interface
	 *
	 * The xdg_surface interface is removed from the wl_surface
	 * object that was turned into a xdg_surface with
	 * xdg_shell.get_xdg_surface request. The xdg_surface properties,
	 * like maximized and fullscreen, are lost. The wl_surface loses
	 * its role as a xdg_surface. The wl_surface is unmapped.
	 */
	void (*destroy)(struct wl_client *client,
			struct wl_resource *resource);
};

#define XDG_POPUP_POPUP_DONE	0

#define XDG_POPUP_POPUP_DONE_SINCE_VERSION	1

static inline void
xdg_popup_send_popup_done(struct wl_resource *resource_, uint32_t serial)
{
	wl_resource_post_event(resource_, XDG_POPUP_POPUP_DONE, serial);
}

#ifdef  __cplusplus
}
#endif

#endif
