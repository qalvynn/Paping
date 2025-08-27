package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"syscall"
	"time"
	"unsafe"
)

// Enable colours directly via winapi instead of using shitty external libs
func enableColours() {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	getConsoleMode := kernel32.NewProc("GetConsoleMode")
	setConsoleMode := kernel32.NewProc("SetConsoleMode")
	stdout := uintptr(syscall.Stdout)
	var mode uint32
	getConsoleMode.Call(stdout, uintptr(unsafe.Pointer(&mode)))
	const ENABLE_VIRTUAL_TERMINAL_PROCESSING = 0x0004
	mode |= ENABLE_VIRTUAL_TERMINAL_PROCESSING
	setConsoleMode.Call(stdout, uintptr(mode))
}

func main() {
	enableColours()
	if len(os.Args) < 3 {
		fmt.Println("Usage: <host> <port> [-h timeout in ms]")
		os.Exit(1)
	}
	host := os.Args[1]
	port := os.Args[2]
	timeoutMs := 100
	for i := 3; i < len(os.Args); i++ {
		if os.Args[i] == "-h" && i+1 < len(os.Args) {
			if t, err := strconv.Atoi(os.Args[i+1]); err == nil {
				timeoutMs = t
			}
			i++
		}
	}

	ip := fmt.Sprintf("%s:%s", host, port)
	timeout := time.Duration(timeoutMs) * time.Millisecond
	for {
		start := time.Now()
		conn, err := net.DialTimeout("tcp", ip, timeout)
		elapsed := time.Since(start)

		if err != nil {
			fmt.Printf("\033[31mConnection timed out\033[0m\n")
		} else {
			ms := float64(elapsed.Nanoseconds()) / 1e6
			fmt.Printf(
				"\033[37mConnected to \033[32m%s\033[37m: time=\033[32m%.2f\033[37mms protocol=\033[32mTCP \033[37mport=\033[32m%s\033[0m\n",
				host, ms, port,
			)
			conn.Close()
		}
		if remaining := timeout - elapsed; remaining > 0 {
			time.Sleep(remaining)
		}
	}
}
