package main

import (
	"fmt"
	"io/ioutil"
	"syscall"
	"time"

	"github.com/beewteam/ips/pkg/irc"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	dat, err := ioutil.ReadFile("../program-version")
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Print(string(dat))

	var nick string
	fmt.Printf("Nick: ")
	fmt.Scanln(&nick)

	fmt.Printf("Pass: ")
	pass, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println()

	var chat string
	fmt.Printf("Chat: ")
	fmt.Scanln(&chat)

	var client irc.Client
	if !client.Connect("irc.freenode.net", "8000") {
		return
	}

	if !client.Login(nick) {
		return
	}

	client.Auth(string(pass))

	time.Sleep(10 * time.Second)

	client.JoinChannel(chat)

	for client.HandleData() {
	}

	client.Close()
}
