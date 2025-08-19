// server.go
package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"

	"github.com/creack/pty"
	"golang.org/x/crypto/ssh"
)

func must(err error) { if err != nil { panic(err) } }

func main() {
	// Host key (server identity, persistent)
	hostKey, err := os.ReadFile("host_ed25519")
	must(err)
	signer, err := ssh.ParsePrivateKey(hostKey)
	must(err)

	cfg := &ssh.ServerConfig{
		Config: ssh.Config{
			KeyExchanges: []string{
				"curve25519-sha256", "curve25519-sha256@libssh.org",
			},
		},
		PasswordCallback: func(meta ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			// Demo: hardcoded username + password
			if meta.User() == "admin" && string(pass) == "pass123" {
				return nil, nil
			}
			return nil, fmt.Errorf("invalid login for %s", meta.User())
		},
	}
	cfg.AddHostKey(signer)

	ln, err := net.Listen("tcp", ":2222")
	must(err)
	fmt.Println("SSH-like server listening on :2222")

	for {
		tcpConn, err := ln.Accept()
		if err != nil {
			continue
		}
		go handleConn(tcpConn, cfg)
	}
}

func handleConn(tcpConn net.Conn, cfg *ssh.ServerConfig) {
	sshConn, chans, reqs, err := ssh.NewServerConn(tcpConn, cfg)
	if err != nil {
		_ = tcpConn.Close()
		return
	}
	defer sshConn.Close()
	go ssh.DiscardRequests(reqs)

	for ch := range chans {
		if ch.ChannelType() != "session" {
			_ = ch.Reject(ssh.UnknownChannelType, "only session channels supported")
			continue
		}
		channel, requests, err := ch.Accept()
		if err != nil {
			continue
		}
		go handleSession(channel, requests)
	}
}

func handleSession(channel ssh.Channel, requests <-chan *ssh.Request) {
	defer channel.Close()

	cmd := exec.Command("sh") // change to "bash" if you prefer
	ptyFile, err := pty.Start(cmd)
	if err != nil {
		_, _ = channel.Stderr().Write([]byte("failed to start shell: " + err.Error()))
		return
	}
	defer ptyFile.Close()

	// connect input/output
	go io.Copy(ptyFile, channel)
	go io.Copy(channel, ptyFile)

	for req := range requests {
		switch req.Type {
		case "pty-req", "shell":
			req.Reply(true, nil)
		case "exec":
			req.Reply(false, nil) // optional: implement one-shot exec
		default:
			req.Reply(false, nil)
		}
	}
	_ = cmd.Wait()
}
