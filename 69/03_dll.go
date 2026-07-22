package main

import (
	"fmt"
	"syscall"
)

func main() {
	user32 := syscall.NewLazyDLL("user32.dll")
	procMessageBox := user32.NewProc("MessageBoxW")

	fmt.Printf("[+] Successfully located user32.dll at address : %X\n", user32.Handle())
	fmt.Printf("[+] Successfully locates MessageBoxW function pointer!\n")
	procMessageBox.Call(0, 0, 0, 0)
}
