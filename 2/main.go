package main

import (
	"fmt"
	"syscall"
)

func main() {
	//1. Declare a local variable and print its memory address
	localVariable := 100
	fmt.Printf("Address of local variable is: %p\n", &localVariable)

	//2. Load kernel32.dll and call GetTickCount
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	getTickCount := kernel32.NewProc("GetTickCount")

	//Call returns (r1, r2, err). r1 is the DWORD return value
	r1, _, _ := getTickCount.Call()
	time := uint32(r1)

	fmt.Printf("Value returned for GetTickCount is: %dms", time)

	//3. Pause and wait for Enter to exit
	fmt.Println("\n Press Enter to exit.....")
	var input string
	fmt.Scanln(&input)
}
