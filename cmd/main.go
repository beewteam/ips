package main

import (
	"fmt"
	"os"
	"time"

	"../pkg/irc"
)

func main() {
	fmt.Printf("Program version: %s\n", VERSION)

	configFile := "./UserConfigs.json"
	var settings = ParseConfig(configFile)

	var server = irc.Server{
		Hostname: "irc.freenode.net",
		Port:     "8000"}
	var client = irc.Client{
		Username: "In work",
		FullName: "In work",
		Server:   server,
	}

	if !server.Connect() {
		os.Exit(1)
	}

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
