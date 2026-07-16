package main

import "C"
import (
	"syscall"

	"golang.org/x/sys/windows"
)

func init() {
	// Spawning a goroutine lets the DLL init finish immediately,
	// avoiding Loader Lock deadlocks.
	go func() {
		windows.MessageBox(
			windows.HWND(0),
			syscall.StringToUTF16Ptr("Hello from inside the Go DLL!"),
			syscall.StringToUTF16Ptr("DLL Message Box"),
			windows.MB_OK|windows.MB_ICONINFORMATION,
		)
	}()
}

// Keep at least one dummy export to ensure cgo compiles properly
// and exports the standard entry points.
//
//export Dummy
func Dummy() {}

func main() {}
