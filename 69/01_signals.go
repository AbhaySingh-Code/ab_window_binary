package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"
)

func main() {
	fmt.Println("[+] Application started. Press Cntl + C to stop me...")

	//Creating a channel that listens for operating system signals
	sigChan := make(chan os.Signal, 1)

	//Nofify our channel when an Interrupt (Ctrl + C) occurs
	signal.Notify(sigChan, os.Interrupt)

	//Create goroutine to wait for the signal
	go func() {
		<-sigChan // This blocks (waits) until Ctrl + C is pressed
		fmt.Println("\n[+] Intercepted Ctrl + C! Performing cleanup here")
		os.Exit(0)
	}()

	for {
		time.Sleep(1 * time.Millisecond)
	}
}
