package main

import (
	"fmt"
	"strconv"
	"syscall"
	"unsafe"

	"github.com/go-check/check"
)

func openEvent(desiredAccess uint32, inheritHandle bool, name string, proc *syscall.LazyProc) (handle syscall.Handle, err error) {
	namep, _ := syscall.UTF16PtrFromString(name)
	var _p2 uint32
	if inheritHandle {
		_p2 = 1
	}
	r0, _, e1 := proc.Call(uintptr(desiredAccess), uintptr(_p2), uintptr(unsafe.Pointer(namep)))
	handle = syscall.Handle(r0)
	if handle == syscall.InvalidHandle {
		err = e1
	}
	return
}

func pulseEvent(handle syscall.Handle, proc *syscall.LazyProc) (err error) {
	r0, _, _ := proc.Call(uintptr(handle))
	if r0 != 0 {
		err = syscall.Errno(r0)
	}
	return
}

func signalDaemonDump(pid int) {
	modkernel32 := syscall.NewLazyDLL("kernel32.dll")
	procOpenEvent := modkernel32.NewProc("OpenEventW")
	procPulseEvent := modkernel32.NewProc("PulseEvent")

	ev := "Global\\docker-daemon-" + strconv.Itoa(pid)
	h2, _ := openEvent(0x0002, false, ev, procOpenEvent)
	if h2 == 0 {
		return
	}
	pulseEvent(h2, procPulseEvent)
}

func signalDaemonReload(pid int) error {
	return fmt.Errorf("daemon reload not supported")
}

func cleanupExecRoot(c *check.C, execRoot string) {
}
