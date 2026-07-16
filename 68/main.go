package main

import (
	"fmt"
	"log"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

func main() {
	// 1. Double check: Ensure this exact path exists on your disk!
	dPath := "C:\\Temp\\dllmain.dll"
	pId := uintptr(22576)

	kernel32 := windows.NewLazyDLL("kernel32.dll")
	pHandle, err := windows.OpenProcess(windows.PROCESS_CREATE_THREAD|windows.PROCESS_VM_OPERATION|windows.PROCESS_VM_WRITE|windows.PROCESS_VM_READ|windows.PROCESS_QUERY_INFORMATION, false, uint32(pId))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Process Opened")

	VirtualAllocEx := kernel32.NewProc("VirtualAllocEx")
	vAlloc, _, err := VirtualAllocEx.Call(uintptr(pHandle), 0, uintptr(len(dPath)+1), windows.MEM_RESERVE|windows.MEM_COMMIT, windows.PAGE_EXECUTE_READWRITE)
	fmt.Println("Memory allocated")

	bPtrDpath, err := windows.BytePtrFromString(dPath)
	if err != nil {
		log.Fatal(err)
	}
	Zero := uintptr(0)
	err = windows.WriteProcessMemory(pHandle, vAlloc, bPtrDpath, uintptr(len(dPath)+1), &Zero)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("DLL path written")

	LoadLibAddr, err := syscall.GetProcAddress(syscall.Handle(kernel32.Handle()), "LoadLibraryA")
	if err != nil {
		log.Fatal(err)
	}

	tHandle, _, _ := kernel32.NewProc("CreateRemoteThread").Call(uintptr(pHandle), 0, 0, LoadLibAddr, vAlloc, 0, 0)
	fmt.Println("DLL Injected, waiting for thread to exit...")

	// --- DEBUG ENGINE ---

	// Wait for the remote thread to finish executing LoadLibraryA
	WaitForSingleObject := kernel32.NewProc("WaitForSingleObject")
	// 0xFFFFFFFF is INFINITE
	WaitForSingleObject.Call(tHandle, uintptr(0xFFFFFFFF))

	// Get the exit code of the thread (which is the HMODULE returned by LoadLibraryA)
	GetExitCodeThread := kernel32.NewProc("GetExitCodeThread")
	var exitCode uint32
	ret, _, _ := GetExitCodeThread.Call(tHandle, uintptr(unsafe.Pointer(&exitCode)))

	if ret == 0 {
		fmt.Println("Debug Error: Failed to call GetExitCodeThread itself.")
	} else if exitCode == 0 {
		fmt.Println("\n--- DIAGNOSIS ---")
		fmt.Println("Result: FAILURE")
		fmt.Println("Reason: The remote thread ran, but LoadLibraryA returned 0 (NULL).")
		fmt.Println("This means either:")
		fmt.Println(" 1. The target process could not find 'dllmain.dll' at the path provided.")
		fmt.Println(" 2. The target process does not have read permissions for your Desktop folder.")
		fmt.Println(" 3. There is an architecture mismatch (e.g., target is 32-bit, DLL is 64-bit).")
	} else {
		fmt.Printf("\n--- DIAGNOSIS ---\nResult: SUCCESS\nDLL loaded in target memory space at base address: 0x%X\n", exitCode)
	}

	syscall.CloseHandle(syscall.Handle(tHandle))
}
