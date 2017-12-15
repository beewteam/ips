package main

import (
	"bufio"
	"fmt"
	"os"

	"../pkg/irc"
)

const (
	configFile = "./UserConfigs.json"
)

func print(reply string, err string) {
	fmt.Println(reply)
}

func main() {
	fmt.Printf("Program version: %s\n", VERSION)

	//var settings = ParseConfig(configFile)

	/*var client = irc.Client{
		Account: irc.Account{
			Username: settings.UserData.Username,
			FullName: settings.UserData.FullName,
		},
		Server: irc.Server{
			Hostname: settings.ServerData.Hostname,
			Port:     settings.ServerData.Port},
	}*/

	//client.Init()
	//client.Run()
	var com irc.Communicator
	com.Init()
	defer com.Close()

	err := com.Run("irc.freenode.com", "8000")
	if err != nil {
		fmt.Printf("Cannot run communicator")
		return
	}

	for {
		text, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		if text[0] == 'q' {
			return
		}
	}
}
