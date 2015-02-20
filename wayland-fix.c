#include "wayland-fix.h"

const struct wl_interface *WL_compositor_interface = &wl_compositor_interface;
const struct wl_interface *WL_callback_interface = &wl_callback_interface;
const struct wl_interface *WL_shell_interface = &wl_shell_interface;
const struct wl_interface *WL_shell_surface_interface = &wl_shell_surface_interface;
const struct wl_interface *WL_xdg_surface_interface= &xdg_surface_interface;
const struct wl_interface *WL_seat_interface= &wl_seat_interface;
