package main

import (
	"fmt"

	"../pkg/irc"
)

const (
	configFile = "./UserConfigs.json"
)

func main() {
	fmt.Printf("Program version: %s\n", VERSION)

	var settings = ParseConfig(configFile)

	var client = irc.Client{
		Account: irc.Account{
			Username: settings.UserData.Username,
			FullName: settings.UserData.FullName,
		},
		Server: irc.Server{
			Hostname: settings.ServerData.Hostname,
			Port:     settings.ServerData.Port},
	}

	client.Run()
}
