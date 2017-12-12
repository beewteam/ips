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
		Hostname: settings.ServerData.Hostname,
		Port:     settings.ServerData.Port}
	var client = irc.Client{
		Username: settings.UserData.Username,
		FullName: settings.UserData.FullName,
		Server:   server,
	}

	if !client.Server.Connect() {
		os.Exit(1)
	}

	if !client.Login(settings.UserData.Nickname) {
		os.Exit(1)
	}

	client.Auth(string(settings.UserData.Password))

	// Should wait NOTIFY message
	time.Sleep(10 * time.Second)

	client.JoinChannel(settings.UserData.Chat)

	for client.HandleData() {
	}

	client.Close()
}
