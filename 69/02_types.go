package main

import (
	"fmt"
)

// type USDollars float64
// type Euros float64

type HANDLE uintptr
type HHOK HANDLE
type WPARAM uintptr

// On a 64 bit machine go automaticaly makes uintptr 64 bits and on a 32 bit machine go automatically makes it 32 bits with uintptr

func main() {
	var rawHookAddress uintptr = 0x7FFF81A2
	var myHook HHOK = HHOK(rawHookAddress)
	fmt.Printf("Raw Type: %T, value: %X\n", rawHookAddress, rawHookAddress)
	fmt.Printf("Custom Hook Type: %T, Value: %X\n", myHook, myHook)
}
