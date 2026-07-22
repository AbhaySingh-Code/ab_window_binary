package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"unsafe"
)

// Updated Windows API constants for the mouse
const (
	WH_MOUSE_LL    = 14     // Low-level mouse hook identifier
	WM_LBUTTONDOWN = 0x0201 // Left mouse button click down
	WM_RBUTTONDOWN = 0x0204 // Right mouse button click down
)

type (
	HANDLE  uintptr
	HHOOK   HANDLE
	WPARAM  uintptr
	LPARAM  uintptr
	LRESULT uintptr
)

// POINT defines X and Y screen coordinates
type POINT struct {
	X int32
	Y int32
}

// MSLLHOOKSTRUCT contains information about a low-level mouse input event
type MSLLHOOKSTRUCT struct {
	Pt          POINT // Where the cursor is on the screen (X, Y)
	MouseData   uint32
	Flags       uint32
	Time        uint32
	DwExtraInfo uintptr
}

var (
	user32                  = syscall.NewLazyDLL("user32.dll")
	procSetWindowsHookExW   = user32.NewProc("SetWindowsHookExW")
	procCallNextHookEx      = user32.NewProc("CallNextHookEx")
	procUnhookWindowsHookEx = user32.NewProc("UnhookWindowsHookEx")
	procGetMessageW         = user32.NewProc("GetMessageW")

	hhk HHOOK
)

// The callback function triggered on every mouse movement or click
func mouseHookCallback(nCode int32, wParam WPARAM, lParam LPARAM) LRESULT {
	if nCode >= 0 {
		// Twist open the capsule using the Mouse structure layout instead
		mouseData := (*MSLLHOOKSTRUCT)(unsafe.Pointer(lParam))

		// Check what kind of mouse event happened
		if wParam == WM_LBUTTONDOWN {
			fmt.Printf("[+] Left Click Detected! Coordinates: X=%d, Y=%d\n", mouseData.Pt.X, mouseData.Pt.Y)
		} else if wParam == WM_RBUTTONDOWN {
			fmt.Printf("[+] Right Click Detected! Coordinates: X=%d, Y=%d\n", mouseData.Pt.X, mouseData.Pt.Y)
		}
	}

	// Pour the event back into the pipe so the computer can actually click
	ret, _, _ := procCallNextHookEx.Call(uintptr(hhk), uintptr(nCode), uintptr(wParam), uintptr(lParam))
	return LRESULT(ret)
}

func main() {
	fmt.Println("[*] Setting up Global Low-Level Mouse Hook...")

	callbackPtr := syscall.NewCallback(mouseHookCallback)

	// Install the hook using WH_MOUSE_LL
	ret, _, err := procSetWindowsHookExW.Call(
		WH_MOUSE_LL,
		callbackPtr,
		0,
		0,
	)

	if ret == 0 {
		fmt.Printf("[-] Failed to set hook. Error: %v\n", err)
		return
	}
	hhk = HHOOK(ret)
	fmt.Println("[+] Hook installed successfully. Press Ctrl+C to exit.")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Println("\n[*] Removing hook and exiting...")
		procUnhookWindowsHookEx.Call(uintptr(hhk))
		os.Exit(0)
	}()

	// Keep the program standing guard
	var msg [11]uint32
	for {
		ret, _, _ := procGetMessageW.Call(uintptr(unsafe.Pointer(&msg[0])), 0, 0, 0)
		if int32(ret) <= 0 {
			break
		}
	}
}
