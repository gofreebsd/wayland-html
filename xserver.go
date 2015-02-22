package main

// #include <wayland-server.h>
// #include "csignal.h"
import "C"

import (
	"fmt"
	_ "github.com/fangyuanziti/wayland-html/cfn"
	"github.com/nightlyone/lockfile"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"
	// "golang.org/x/sys/unix"
)

var LOCK_FMT string = "/tmp/.X%d-lock"

func lock(num int) (lockfile.Lockfile, error) {
	fileName := fmt.Sprintf(LOCK_FMT, num)
	lock, err := lockfile.New(fileName)
	if err != nil {
		return lock, err
	}

	err = lock.TryLock()

	// Error handling is essential, as we only try to get the lock.
	if err != nil {
		return lock, err
	}

	return lock, nil
}

func TryLock() (int, lockfile.Lockfile) {
	i := 0
	for {
		lockFile, err := lock(i)
		if err == nil {
			return i, lockFile
		}
		i = i + 1
	}
}

func initXwm() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGUSR1)
	go func() {
		<-sigc
	}()
}

func listen(path string) (*net.UnixListener, error) {
	_, err := os.Stat(path)
	if err == nil {
		os.Remove(path)
	}

	listener, err := net.Listen("unix", path)
	if err != nil {
		return nil, err
	}

	unixListener, _ := listener.(*net.UnixListener)
	if err != nil {
		return nil, err
	}
	return unixListener, nil
}

var SOCKET_FMT string = "/tmp/.X11-unix/X"

func forkXWayland() uintptr {
	ret, _, err := syscall.Syscall(syscall.SYS_FORK, 0, 0, 0)
	if err != 0 {
		return ret
	}
	return ret
}

func xserverInit(display *C.struct_wl_display) {
	// Fetch a valid lock file and DISPLAY number
	displayNum, _ := TryLock()
	numStr := strconv.Itoa(displayNum)
	// Set DISPLAY number
	displayName := ":" + numStr
	os.Setenv("DISPLAY", displayName)

	// init xwm
	initXwm()

	// init DISPLAY unix socket
	unixListener, _ := listen(SOCKET_FMT + numStr)
	abstructListener, _ := listen("@" + SOCKET_FMT + numStr)
	// do not close the listener(Close them will remove the file in filesystem).

	unixFile, _ := unixListener.File()
	abstructFile, _ := abstructListener.File()
	unixFd := unixFile.Fd()
	abstructFd := abstructFile.Fd()

	// Open a socket for the Wayland connection from Xwayland.
	wls, _ := syscall.Socketpair(syscall.AF_UNIX,
		syscall.SOCK_STREAM, 0)

	wms, _ := syscall.Socketpair(syscall.AF_UNIX,
		syscall.SOCK_STREAM, 0)

	client := C.wl_client_create(display, (C.int)(wls[0]))
	println(client)

	pid := forkXWayland()
	if pid == 0 { // child
		// unix.Close(wls[0])
		// unix.Close(wls[1])
		C.signal_ignore((C.int)(syscall.SIGUSR1))
		os.Setenv("WAYLAND_SOCKET", strconv.Itoa(wls[1]))
		args := []string{
			"Xwayland",
			"-rootless",
			"-terminate",
			"-listen", strconv.Itoa((int)(unixFd)),
			"-listen", strconv.Itoa((int)(abstructFd)),
			"-wm", strconv.Itoa(wms[1]),
		}

		binary, lookErr := exec.LookPath("Xwayland")
		if lookErr != nil {
			panic(lookErr)
		}

		env := os.Environ()
		execErr := syscall.Exec(binary, args, env)
		if execErr != nil {
			panic(lookErr)
		}

		println(wls[1], wms[1], unixFd, abstructFd)
	} else { // parent

	}

}
