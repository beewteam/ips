package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"syscall"
	"time"

	"github.com/beewteam/ips/pkg/irc"
	"golang.org/x/crypto/ssh/terminal"
)

type StartupSettings struct {
	Nickname string
	Password string
	Chat     string
}

func main() {
	configFile := "./UserConfigs.json"

	dat, err := ioutil.ReadFile("../program-version")
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
	fmt.Print(string(dat))

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

	client := irc.Client{}
	if !client.Connect("irc.freenode.net", "8000") {
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
