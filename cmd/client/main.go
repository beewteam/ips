package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/beewteam/ips/pkg/irc"
	"github.com/chzyer/readline"
	"github.com/fatih/color"
)

var configFile string

// func printM(reply string) {
// 	fmt.Println(reply)
// }

// func print(reply string, err string) {
// 	fmt.Println(reply)
// }

func main() {
	var settings UserConfig

	fmt.Println("Client version: " + color.GreenString(VERSION))

	flag.StringVar(&configFile, "c", "UserConfigs.json", "path to config file in json, default==UserConfigs.json in pwd")
	flag.Parse()

	rl, err := readline.New("irc> ")
	if err != nil {
		fmt.Println(color.RedString("client") + ": cannot init input interface")
		panic(err)
	}
	defer rl.Close()

	fmt.Printf("Trying to load %s\n", color.GreenString(configFile))
	if _, err = os.Stat(configFile); err != nil {
		settings, err = ParseConfigFile(configFile)
		if err != nil {
			fmt.Println(color.RedString("client") + ": cannot load config")
			panic(err)
		}
	}

	var client = irc.Client{
		Account: irc.Account{
			Username: settings.UserData.Username,
			FullName: settings.UserData.FullName,
		},
		Server: irc.Server{
			Hostname: settings.ServerData.Hostname,
			Port:     settings.ServerData.Port},
	}

	fmt.Println(client.Account.Username)

	// var com irc.Communicator
	// com.Init()
	// defer com.Close()
	// errlog := com.SetLog("irc.log")
	// if errlog != nil {
	// 	fmt.Println("Cannot setup log file")
	// 	return
	// }
	// com.Subscribe("*", printM)

	// err := com.Run("irc.freenode.com", "8000")
	// if err != nil {
	// 	fmt.Printf("Cannot run communicator\n")
	// 	return
	// }

	// for {
	// 	text, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	// 	if text[0] == 'q' {
	// 		return
	// 	} else if text[0] == 'p' {
	// 		com.SendMessage("PING", print, "irc.freenode.com")
	// 	} else if text[0] == 'u' {
	// 		com.SendMessage("USER", "testion", "0", "*", "Oliak")
	// 	}
	// }
}
