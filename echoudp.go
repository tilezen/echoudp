package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

func dieUsage(args []string) {
	fatalf("Usage: %s <port>\n", args[0])
}

func fatalf(msg string, args ...interface{}) {
	fmt.Printf(msg, args...)
	os.Exit(1)
}

func main() {
	if len(os.Args) != 2 {
		dieUsage(os.Args)
	}
	portStr := os.Args[1]
	port, err := strconv.ParseUint(portStr, 10, 32)
	if err != nil {
		dieUsage(os.Args)
	}
	listenSpec := fmt.Sprintf(":%d", port)

	addr, err := net.ResolveUDPAddr("udp", listenSpec)
	if err != nil {
		fatalf("Invalid addr %s: %v\n", listenSpec, err)
	}

	udpConn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fatalf("Failed to listen on %s: %v\n", listenSpec, err)
	}

	c := make(chan struct{}, 1)

	go func() {
		var input string
		fmt.Scanln(&input)
		c <- struct{}{}
	}()

	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := udpConn.Read(buf)
			if err != nil {
				fatalf("Failed to read from connection: %v\n", err)
			}
			if n > 0 {
				readBytes := buf[:n]
				s := string(readBytes)
				fmt.Printf(s)
			}
		}
	}()

	<-c
	os.Exit(0)

}
