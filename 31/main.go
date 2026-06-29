package main

import (
	"fmt"
	"os"
	"unsafe"

	"golang.org/x/sys/windows"
)

func main() {
	//1. Setup the structures
	var si windows.StartupInfo
	var pi windows.ProcessInformation

	si.Cb = uint32(unsafe.Sizeof(si))

	//Path to executable
	commandLine := "C:\\Windows\\System32\\notepad.exe"

	fmt.Printf("[+] Attempting to create process: %s\n", commandLine)

	//Convert the string to a UTF-16 pointer (*uint16) for windows api
	commandLinePtr, err := windows.UTF16PtrFromString(commandLine)
	if err != nil {
		fmt.Printf("[-]Failed to convert string to UTF16: %v\n", err)
		return
	}

	//2. Call create process
	//We use CREATE_SUSPENDED to start notepad frozen
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
		fmt.Printf("[-] Create process failed. Error: %v\n", err)
		return
	}

	//3. Inspect what we created
	fmt.Printf("[+]Process created sucessfully!")
	fmt.Printf("[+]ProcessID: %d\n", pi.ProcessId)
	fmt.Printf("[+]Thread ID: %d\n", pi.ThreadId)
	fmt.Println("[!] Notepad is currently frozen in memory. Check task manager!")

	fmt.Println("\nPress Enter to resume the process.............")
	var input string
	fmt.Fscanln(os.Stdin, &input)

	//4. Resume the process
	//We resume the 'Primary Thread' handle returned in pi
	_, err = windows.ResumeThread(pi.Thread)
	if err != nil {
		fmt.Printf("[-]Failed to resume thread: %v\n", err)
	} else {
		fmt.Println("[+]Process Resumed.")
	}

	//5. Cleanup
	// Defering or manually closing handles to prevent resource leaks
	windows.CloseHandle(pi.Process)
	windows.CloseHandle(pi.Thread)
}
