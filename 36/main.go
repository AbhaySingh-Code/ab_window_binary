package main

import (
	"fmt"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

type ThreadArgs struct {
	ID    uint32
	Score float64
	Name  [32]byte
}

func main() {
	kernel32 := windows.NewLazySystemDLL("kernel32.dll")
	createThread := kernel32.NewProc("CreateThread")
	args := ThreadArgs{
		ID:    2026,
		Score: 100.0,
	}
	copy(args.Name[:], "GoWinSysThread")

	fmt.Println("[Main] Launching Windows thread via x/sys/windows...")

	threadHandle, _, err := createThread.Call(
		0,
		0,
		windows.NewCallback(threadProc),
		uintptr(unsafe.Pointer(&args)),
		0,
		0,
	)

	if threadHandle == 0 {
		panic(fmt.Sprintf("Failed to create thread. Windows Error: %v", err))
	}

	hThread := windows.Handle(threadHandle)
	defer windows.CloseHandle(hThread)

	time.Sleep(2 * time.Second)
	fmt.Println("[Main] Exiting.")
}

func threadProc(lpParameter uintptr) uintptr {
	args := (*ThreadArgs)(unsafe.Pointer(lpParameter))
	name := windows.ByteSliceToString(args.Name[:])
	fmt.Printf("[Thread] Hello from the API Thread!\n")
	fmt.Printf("[Thread] Received Args -> ID: %d, Score: %.1f, Name: %s\n", args.ID, args.Score, name)
	return 0
}
