// client.go
package main

import (
	"fmt"
	"os"
	"strings"
	"bufio"
	"golang.org/x/crypto/ssh"
)

func must(err error) { if err != nil { panic(err) } }

func main() {
	fmt.Println("Connecting to SSH-like server at 127.0.0.1:2222\n")
	fmt.Print("Enter password for user 'admin': ")
	password, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	new_password := strings.TrimSpace(password)
	cfg := &ssh.ClientConfig{
		User: "admin",
		Auth: []ssh.AuthMethod{
			ssh.Password(new_password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // accept any host key (⚠️ unsafe, but simple)
		Config: ssh.Config{
			KeyExchanges: []string{
				"curve25519-sha256", "curve25519-sha256@libssh.org",
			},
		},
	}

	client, err := ssh.Dial("tcp", "127.0.0.1:2222", cfg)
	must(err)
	defer client.Close()

	sess, err := client.NewSession()
	must(err)
	defer sess.Close()

	// allocate PTY
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	must(sess.RequestPty("xterm-256color", 24, 80, modes))

	sess.Stdin = os.Stdin
	sess.Stdout = os.Stdout
	sess.Stderr = os.Stderr

	must(sess.Shell())
	must(sess.Wait())
	fmt.Println("session closed")
}
