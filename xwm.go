package main

import (
	"fmt"
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/composite"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/icccm"
	"github.com/BurntSushi/xgbutil/xwindow"
	log "github.com/Sirupsen/logrus"
	"net"
	"os"
)

func fromFd(fd uintptr) *xgbutil.XUtil {

	file := os.NewFile(fd, "wm")
	netConn, _ := net.FileConn(file)

	xgbConn, err := xgb.NewConnNet(netConn)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	X, err := xgbutil.NewConnXgb(xgbConn)

	if err != nil {
		panic(err)
	}
	return X
}

type XWin struct {
	surfaceId int
}

var XWins = make(map[xproto.Window]*XWin)

func createNotify(X *xgbutil.XUtil, event xgb.Event) {
	ev := event.(xproto.CreateNotifyEvent)
	// values[0] = XCB_EVENT_MASK_PROPERTY_CHANGE;
	// xcb_change_window_attributes(wm->conn, id, XCB_CW_EVENT_MASK, values);
	win := xwindow.New(X, ev.Window)
	win.Listen(xproto.EventMaskPropertyChange)

	XWins[ev.Window] = &XWin{
		surfaceId: 0,
	}
}

func destroyNotify(X *xgbutil.XUtil, event xgb.Event) {
	ev := event.(xproto.DestroyNotifyEvent)
	xproto.DestroyWindow(X.Conn(), ev.Window)
	delete(XWins, ev.Window)
}
func configureNotify(X *xgbutil.XUtil, event xgb.Event) {
	ev := event.(xproto.ConfigureNotifyEvent)
	log.Info("Event: %s\n", ev)
}
func clientMessage(X *xgbutil.XUtil, event xgb.Event) {
	ev := event.(xproto.ClientMessageEvent)
	cookie := xproto.GetAtomName(X.Conn(), ev.Type)
	reply, _ := cookie.Reply()
	log.WithFields(log.Fields{
		"name": reply.Name,
		"data": ev.Data,
	}).Info("client Message:", ev)
}

func configureRequest(X *xgbutil.XUtil, event xgb.Event) {
	ev := event.(xproto.ConfigureRequestEvent)
	win := xwindow.New(X, ev.Window)

	win.Configure(
		xproto.ConfigWindowX|
			xproto.ConfigWindowY|
			xproto.ConfigWindowWidth|
			xproto.ConfigWindowHeight,
		(int)(ev.X), (int)(ev.Y),
		(int)(ev.Width), (int)(ev.Height),
		ev.Sibling,
		ev.StackMode,
	)
	log.Info("configure request:", ev)
}
func unmapNotify(X *xgbutil.XUtil, event xgb.Event) {
	ev := event.(xproto.UnmapNotifyEvent)
	win := xwindow.New(X, ev.Window)
	log.Info("Unmap:", win)
	win.Unmap()
}

func xwmInit(fd uintptr) {

	X := fromFd(fd)
	defer X.Conn().Close()

	root := xwindow.New(X, X.RootWin())

	if _, err := root.Geometry(); err != nil {
		panic(err)
	}

	// if names, _ := ewmh.DesktopNamesGet(X); len(names) > 0 {
	// 	println(names)
	// }

	composite.Init(X.Conn())

	atomNames := []string{
		"WL_SURFACE_ID",
		"WM_DELETE_WINDOW",
		"WM_PROTOCOLS",
		"WM_S0",
		"WM_NORMAL_HINTS",
		"WM_TAKE_FOCUS",
		"WM_STATE",
		"WM_CLIENT_MACHINE",
		"_NET_WM_CM_S0",
		"_NET_WM_NAME",
		"_NET_WM_PID",
		"_NET_WM_ICON",
		"_NET_WM_STATE",
		"_NET_WM_STATE_FULLSCREEN",
		"_NET_WM_USER_TIME",
		"_NET_WM_ICON_NAME",
		"_NET_WM_WINDOW_TYPE",
		"_NET_WM_WINDOW_TYPE_DESKTOP",
		"_NET_WM_WINDOW_TYPE_DOCK",
		"_NET_WM_WINDOW_TYPE_TOOLBAR",
		"_NET_WM_WINDOW_TYPE_MENU",
		"_NET_WM_WINDOW_TYPE_UTILITY",
		"_NET_WM_WINDOW_TYPE_SPLASH",
		"_NET_WM_WINDOW_TYPE_DIALOG",
		"_NET_WM_WINDOW_TYPE_DROPDOWN_MENU",
		"_NET_WM_WINDOW_TYPE_POPUP_MENU",
		"_NET_WM_WINDOW_TYPE_TOOLTIP",
		"_NET_WM_WINDOW_TYPE_NOTIFICATION",
		"_NET_WM_WINDOW_TYPE_COMBO",
		"_NET_WM_WINDOW_TYPE_DND",
		"_NET_WM_WINDOW_TYPE_NORMAL",
		"_NET_WM_MOVERESIZE",
		"_NET_SUPPORTING_WM_CHECK",
		"_NET_SUPPORTED",
		"_MOTIF_WM_HINTS",
		"CLIPBOARD",
		"CLIPBOARD_MANAGER",
		"TARGETS",
		"UTF8_STRING",
		"_WL_SELECTION",
		"INCR",
		"TIMESTAMP",
		"MULTIPLE",
		"UTF8_STRING",
		"COMPOUND_TEXT",
		"TEXT",
		"STRING",
		"text/plain;charset=utf-8",
		"text/plain",
		"XdndSelection",
		"XdndAware",
		"XdndEnter",
		"XdndLeave",
		"XdndDrop",
		"XdndStatus",
		"XdndFinished",
		"XdndTypeList",
		"XdndActionCopy",
	}

	cookies := make([]xproto.InternAtomCookie, len(atomNames))

	for i := 0; i < len(atomNames); i++ {
		cookies[i] = xproto.InternAtom(X.Conn(),
			false, uint16(len(atomNames[i])), atomNames[i])
	}

	/* Try to select for substructure redirect. */
	evMasks := xproto.EventMaskPropertyChange |
		xproto.EventMaskFocusChange |
		xproto.EventMaskButtonPress |
		xproto.EventMaskButtonRelease |
		xproto.EventMaskStructureNotify |
		xproto.EventMaskSubstructureNotify |
		xproto.EventMaskSubstructureRedirect

	root.Listen(evMasks)

	composite.RedirectSubwindows(X.Conn(), root.Id,
		composite.RedirectManual)

	// change config
	// data := []byte{1}
	// propAtom, _ := xprop.Atm(X, "_NET_SUPPORTED")
	// var format byte = 32
	// xproto.ChangePropertyChecked(X.Conn(), xproto.PropModeReplace,
	// 	root.Id, propAtom, xproto.AtomAtom, format,
	// 	uint32(len(data)/(int(format)/8)), data)

	ewmh.SupportedSet(X, []string{
		"_NET_SUPPORTED",
		"_NET_CLIENT_LIST",
		"_NET_NUMBER_OF_DESKTOPS",
		"_NET_DESKTOP_GEOMETRY",
		"_NET_CURRENT_DESKTOP",
		"_NET_VISIBLE_DESKTOPS",
		"_NET_DESKTOP_NAMES",
		"_NET_ACTIVE_WINDOW",
		"_NET_SUPPORTING_WM_CHECK",
		"_NET_CLOSE_WINDOW",
		"_NET_MOVERESIZE_WINDOW",
		"_NET_RESTACK_WINDOW",
		"_NET_WM_NAME",
		"_NET_WM_DESKTOP",
		"_NET_WM_WINDOW_TYPE",
		"_NET_WM_WINDOW_TYPE_DESKTOP",
		"_NET_WM_WINDOW_TYPE_DOCK",
		"_NET_WM_WINDOW_TYPE_TOOLBAR",
		"_NET_WM_WINDOW_TYPE_MENU",
		"_NET_WM_WINDOW_TYPE_UTILITY",
		"_NET_WM_WINDOW_TYPE_SPLASH",
		"_NET_WM_WINDOW_TYPE_DIALOG",
		"_NET_WM_WINDOW_TYPE_DROPDOWN_MENU",
		"_NET_WM_WINDOW_TYPE_POPUP_MENU",
		"_NET_WM_WINDOW_TYPE_TOOLTIP",
		"_NET_WM_WINDOW_TYPE_NOTIFICATION",
		"_NET_WM_WINDOW_TYPE_COMBO",
		"_NET_WM_WINDOW_TYPE_DND",
		"_NET_WM_WINDOW_TYPE_NORMAL",
		"_NET_WM_STATE",
		"_NET_WM_STATE_STICKY",
		"_NET_WM_STATE_MAXIMIZED_VERT",
		"_NET_WM_STATE_MAXIMIZED_HORZ",
		"_NET_WM_STATE_SKIP_TASKBAR",
		"_NET_WM_STATE_SKIP_PAGER",
		"_NET_WM_STATE_HIDDEN",
		"_NET_WM_STATE_FULLSCREEN",
		"_NET_WM_STATE_ABOVE",
		"_NET_WM_STATE_BELOW",
		"_NET_WM_STATE_DEMANDS_ATTENTION",
		"_NET_WM_STATE_FOCUSED",
		"_NET_WM_ALLOWED_ACTIONS",
		"_NET_WM_ACTION_MOVE",
		"_NET_WM_ACTION_RESIZE",
		"_NET_WM_ACTION_MINIMIZE",
		"_NET_WM_ACTION_STICK",
		"_NET_WM_ACTION_MAXIMIZE_HORZ",
		"_NET_WM_ACTION_MAXIMIZE_VERT",
		"_NET_WM_ACTION_FULLSCREEN",
		"_NET_WM_ACTION_CHANGE_DESKTOP",
		"_NET_WM_ACTION_CLOSE",
		"_NET_WM_ACTION_ABOVE",
		"_NET_AM_ACTION_BELOW",
		"_NET_WM_STRUT_PARTIAL",
		"_NET_WM_ICON",
		"_NET_FRAME_EXTENTS",
		"WM_TRANSIENT_FOR",
	})

	// create wm window
	win, _ := xwindow.Create(X, root.Id)

	xproto.MapWindow(X.Conn(), win.Id)

	atomValues := make([]xproto.Atom, len(atomNames))
	for num, cookie := range cookies {
		reply, _ := cookie.Reply()
		atomValues[num] = reply.Atom
	}

	/* take WM_S0 selection last, which
	 * signals to Xwayland that we're done with setup. */
	xproto.SetSelectionOwner(X.Conn(), win.Id,
		atomValues[3],
		xproto.TimeCurrentTime,
	)

	mapResquest := func(event xgb.Event) {
		ev, _ := event.(xproto.MapRequestEvent)
		log.Info("MapRequest: ", ev)

		xproto.MapWindow(X.Conn(), ev.Window)
		icccm.WmStateSet(X, ev.Window, &icccm.WmState{
			State: icccm.StateNormal,
			Icon:  ev.Window,
		})
		// TODO: set _net_wm_state
	}

	for {
		x := X.Conn()
		ev, xerr := x.WaitForEvent()
		if ev == nil && xerr == nil {
			fmt.Println("Both event and error are nil. Exiting...")
			return
		}
		// if ev != nil {
		// 	fmt.Printf("Event: %s\n", ev)
		// }
		if xerr != nil {
			fmt.Printf("Error: %s\n", xerr)
		}

		switch ev.(type) {
		case xproto.CreateNotifyEvent:
			createNotify(X, ev)
			fmt.Printf("Event: %s\n", ev)
		case xproto.DestroyNotifyEvent:
			destroyNotify(X, ev)
			fmt.Printf("Event: %s\n", ev)
		case xproto.MapRequestEvent:
			mapResquest(ev)
		case xproto.ConfigureNotifyEvent:
			configureNotify(X, ev)
		case xproto.PropertyNotifyEvent:
			fmt.Printf("Event: %s\n", ev)
		case xproto.ClientMessageEvent:
			clientMessage(X, ev)
		case xproto.ConfigureRequestEvent:
			configureRequest(X, ev)
		case xproto.UnmapNotifyEvent:
			unmapNotify(X, ev)
		default:
			log.Info("Event:", ev)
		}

	}

	// xevent.MapRequestFun(
	// 	func(X *xgbutil.XUtil, e xevent.MapRequestEvent) {
	// 		println(e.Window)
	// 	}).Connect(X, root.Id)

	// xevent.ConfigureNotifyFun(
	// 	func(X *xgbutil.XUtil, e xevent.ConfigureNotifyEvent) {
	// 		fmt.Printf("(%d, %d) %dx%d\n", e.X, e.Y, e.Width, e.Height)
	// 	}).Connect(X, root.Id)

	// xevent.FocusInFun(
	// 	func(X *xgbutil.XUtil, e xevent.FocusInEvent) {
	// 		fmt.Printf("(%v, %v)\n", e.Mode, e.Detail)
	// 	}).Connect(X, root.Id)

	// xevent.Main(X)
}
