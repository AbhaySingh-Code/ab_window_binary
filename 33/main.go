package main

import (
	"fmt"
	"os"
	"unsafe"

	"golang.org/x/sys/windows"
)

func main() {

	var si windows.StartupInfo
	var pi windows.ProcessInformation
	si.Cb = uint32(unsafe.Sizeof(si))
	commandLine := "C:\\Windows\\System32\\notepad.exe"

	fmt.Printf("[+] Attempting to create process: %s\n", commandLine)

	commandLinePtr, err := windows.UTF16PtrFromString(commandLine)
	if err != nil {
		fmt.Printf("[-] Failed to convert string to UTF16: %v\n", err)
		return
	}

	err = windows.CreateProcess(
		nil,
		commandLinePtr,
		nil,
		nil,
		false,
		windows.CREATE_SUSPENDED,
		nil,
		nil,
		&si,
		&pi,
	)

	if err != nil {
		fmt.Printf("[-] Create process failed. Error : %v\n", err)
		return
	}

	fmt.Print("[+] Proces Created successfully!")
	fmt.Printf("[+] Process ID: %d\n", pi.ProcessId)
	fmt.Printf("[+] Thread ID: %d\n", pi.ThreadId)
	fmt.Printf("[+] Notepad is currently frozen in memory. Check task manager!")

	fmt.Println("\n Press enter to resume ..........")
	var input string
	fmt.Fscanln(os.Stdin, &input)

	_, err = windows.ResumeThread(pi.Thread)
	if err != nil {
		fmt.Printf("[-] Failed to resume thread: %v\n", err)
	} else {
		fmt.Println("[+] Process resumed")
	}

	windows.CloseHandle(pi.Process)
	windows.CloseHandle(pi.Thread)
}
