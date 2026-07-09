package main

/*
#include <windows.h>

static HANDLE OpenThreadC(DWORD tid) {
    return OpenThread(THREAD_GET_CONTEXT | THREAD_SUSPEND_RESUME, FALSE, tid);
}

static int GetCtx(HANDLE h, CONTEXT* ctx) {
    ctx->ContextFlags = CONTEXT_FULL;
    SuspendThread(h);
    int ok = GetThreadContext(h, ctx);
    ResumeThread(h);
    CloseHandle(h);
    return ok;
}
*/
import "C"

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	tidCh := make(chan uint32)

	wg.Add(1)
	go func() {
		defer wg.Done()
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		tidCh <- uint32(C.GetCurrentThreadId())
		time.Sleep(5 * time.Second)
	}()

	tid := <-tidCh
	fmt.Printf("Worker thread ID: %d\n", tid)

	h := C.OpenThreadC(C.DWORD(tid))
	if h == nil {
		panic("OpenThread failed")
	}

	var ctx C.CONTEXT
	if C.GetCtx(h, &ctx) == 0 {
		panic("GetThreadContext failed")
	}

	fmt.Printf("RIP: 0x%X\n", ctx.Rip)
	fmt.Printf("RSP: 0x%X\n", ctx.Rsp)
	fmt.Printf("RAX: 0x%X\n", ctx.Rax)

	wg.Wait()
}
