package main

import (
	"fmt"
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/composite"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	// "github.com/BurntSushi/xgbutil/ewmh"
	// "github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xwindow"
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
	win := XWin{
		surfaceId: 0,
	}

	XWins[ev.Window] = &win
}

func destroyNotify(X *xgbutil.XUtil, event xgb.Event) {
	ev := event.(xproto.DestroyNotifyEvent)
	delete(XWins, ev.Window)
	xproto.DestroyWindow(X.Conn(), ev.Window)
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

	flagRequests := 4
	atomNames := []string{
		"WL_SURFACE_ID",
		"WM_DELETE_WINDOW",
		"WM_PROTOCOLS",
		"WM_S0",
	}

	cookies := make([]xproto.InternAtomCookie, len(atomNames))

	for i := 0; i < flagRequests; i++ {
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
		fmt.Printf("Event: %s\n", ev)

		xproto.MapWindow(X.Conn(), ev.Window)
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
			fmt.Printf("Event: %s\n", ev)
		case xproto.PropertyNotifyEvent:
			fmt.Printf("Event: %s\n", ev)
		case xproto.ClientMessageEvent:
			fmt.Printf("Event: %s\n", ev)
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
