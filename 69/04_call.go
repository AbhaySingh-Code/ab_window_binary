package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

func main() {
	user32 := syscall.NewLazyDLL("user32.dll")
	procMessageBox := user32.NewProc("MessageBoxW")

	// SetWindowHookExW call method only accepts arguments as uintptr. So we convert Go string to a pointer to a UTF 16 character array
	//Convert Go strings into Windows-compatible UTF-16 pointers
	textPtr, _ := syscall.UTF16PtrFromString("Hello From GO!")
	titlePtr, _ := syscall.UTF16PtrFromString("Alert Box")

	ret, _, _ := procMessageBox.Call(
		0,
		uintptr(unsafe.Pointer(textPtr)),
		uintptr(unsafe.Pointer(titlePtr)),
		0,
	)

	fmt.Printf("[+] MessageBox closed. Return code from windows: %d\n", ret)
}
