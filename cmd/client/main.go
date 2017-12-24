package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/beewteam/ips/pkg/irc"
)

const (
	configFile = "./UserConfigs.json"
)

func printM(reply string) {
	fmt.Println(reply)
}

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

	var com irc.Communicator
	com.Init()
	defer com.Close()
	errlog := com.SetLog("irc.log")
	if errlog != nil {
		fmt.Println("Cannot setup log file")
		return
	}
	com.Subscribe("*", printM)

	err := com.Run("irc.freenode.com", "8000")
	if err != nil {
		fmt.Printf("Cannot run communicator\n")
		return
	}

	for {
		text, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		if text[0] == 'q' {
			return
		} else if text[0] == 'p' {
			com.SendMessage("PING", print, "irc.freenode.com")
		}
	}
}
