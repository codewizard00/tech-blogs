package main

import (
	"bufio"
	"fmt"
	"net"
	"os/exec"
	"strings"
	"io"
	"github.com/creack/pty"
)

const password = "secret123"

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Println("Telnet-like server on port 8080...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	conn.Write([]byte("Password: "))
	pass, _ := reader.ReadString('\n')
	if strings.TrimSpace(pass) != password {
		conn.Write([]byte("Authentication failed.\n"))
		return
	}

	conn.Write([]byte("Welcome! You now have a real shell.\n"))

	// Start bash (or sh if bash not available)
	cmd := exec.Command("bash")

	// Create a pty
	ptmx, err := pty.Start(cmd)
	if err != nil {
		conn.Write([]byte("Failed to start shell.\n"))
		return
	}
	defer func() { _ = ptmx.Close() }() // close the PTY

	// Copy data between PTY and network connection
	go func() { _, _ = io.Copy(ptmx, conn) }()
	_, _ = io.Copy(conn, ptmx)
}

