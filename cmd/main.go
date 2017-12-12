package main

import (
	"encoding/json"
	"fmt"
	"os"
	"syscall"
	"time"

	"../pkg/irc"

	"golang.org/x/crypto/ssh/terminal"
)

type StartupSettings struct {
	Nickname string
	Password string
	Chat     string
}

func main() {
	fmt.Printf("Program version: %s\n", VERSION)

	configFile := "./UserConfigs.json"

	settings := StartupSettings{}
	if _, err := os.Stat(configFile); err == nil {
		file, err := os.Open(configFile)
		defer file.Close()
		if err != nil {
			fmt.Println(err.Error())
		}

		err = json.NewDecoder(file).Decode(&settings)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	} else {
		fmt.Printf("Nick: ")
		fmt.Scanln(&settings.Nickname)

		fmt.Printf("Pass: ")
		passBytes, _ := terminal.ReadPassword(int(syscall.Stdin))
		settings.Password = string(passBytes)
		fmt.Println()

		fmt.Printf("Chat: ")
		fmt.Scanln(&settings.Chat)
	}

	var server = irc.Server{
		Hostname: "irc.freenode.net",
		Port:     "8000"}
	client := irc.Client{}
	client.SetNames("opq", "eqq")
	if !server.Connect() {
		os.Exit(1)
	}
	client.Server = server

	if !client.Login(settings.Nickname) {
		os.Exit(1)
	}

	client.Auth(string(settings.Password))

	// Should wait NOTIFY message
	time.Sleep(10 * time.Second)

	client.JoinChannel(settings.Chat)

	for client.HandleData() {
	}

	client.Close()
}
