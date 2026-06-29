package main

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// FindSubsystems searches for a given process name in the system snapshot
func FindSubsystemProcess(targetName string) {
	//Take a snapshot of all the process in windows
	hProcessSnap, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		fmt.Println("[-] Failed to take snapshot: ", err)
		return
	}
	//Ensure the handleis closed when the function finishes
	defer windows.CloseHandle(hProcessSnap)

	var pe32 windows.ProcessEntry32
	pe32.Size = uint32(unsafe.Sizeof(pe32))

	//Get the first process
	err = windows.Process32First(hProcessSnap, &pe32)
	if err != nil {
		fmt.Println("[-] Failed to retrieve first process: ", err)
		return
	}

	found := false

	for {
		//Covert the null-terminated UTF-16 array (szExeFile) to a Go string
		processName := windows.UTF16ToString(pe32.ExeFile[:])

		if processName == targetName {
			fmt.Printf("[+] Found %s | PID: %d | ParentPID: %d\n", processName, pe32.ProcessID, pe32.ParentProcessID)
			found = true
		}
		//Move to next processs
		err = windows.Process32Next(hProcessSnap, &pe32)
		if err != nil {
			// ERROR_NO_MORE_FILES indicates we have reached the end of the list
			if err == syscall.ERROR_NO_MORE_FILES {
				break
			}
			fmt.Println("[-] Error iterating processes: ", err)
			break
		}
	}

	if !found {
		fmt.Printf("[-] Could not find %s (Are you running as an Administrator?)\n", targetName)
	}
}

func main() {
	fmt.Println("--------- Window Subsystem Discovery --------------")

	subsystems := []string{"smss.exe", "csrss.exe", "lsass.exe"}

	for _, name := range subsystems {
		FindSubsystemProcess(name)
	}
	fmt.Println("\n Press enter to exit .....")
	var input string
	fmt.Scanln(&input)
}
