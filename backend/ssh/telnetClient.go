package main

import (
	"fmt"
	"net"
	"bufio"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")

	if(err!=nil){
		panic(err);
	}

	defer conn.Close();

	fmt.Println("Connected to server at:", conn.RemoteAddr())

	go func() {
		reader := bufio.NewReader(conn)
		for {
			msg, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("❌ Server disconnected")
				os.Exit(0)
			}
			fmt.Print(msg)
		}
	}()

	for {
		msg, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		conn.Write([]byte(msg))
		fmt.Println("📤 Sent to server:", msg);
	}
}