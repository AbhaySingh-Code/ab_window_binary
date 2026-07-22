package main

import (
	"fmt"
	"unsafe"
)

type SampleStruct struct {
	ValueA uint32
	ValueB uint32
}

// C give go a raw memory address, we must tell go to look at exact memory location and read it as if it were a structured go object

func main() {
	data := SampleStruct{ValueA: 42, ValueB: 99}

	// we take memroy address of 'data' and turn it into a raw uintptr
	rawAddress := uintptr(unsafe.Pointer(&data))
	fmt.Printf("[+] Data lives at raw memory address : 0x%X\n", rawAddress)

	reconstructData := (*SampleStruct)(unsafe.Pointer(rawAddress))

	fmt.Printf("[+] Extracted ValueA: %d\n", reconstructData.ValueA)
	fmt.Printf("[+] Extracted ValueB: %d\n", reconstructData.ValueB)
}
