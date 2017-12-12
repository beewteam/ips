package main

import (
	"encoding/json"
	"fmt"
	"os"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

type StartupSettings struct {
	Nickname string
	Password string
	Chat     string
}

func ParseConfig(configPath string) StartupSettings {
	settings := StartupSettings{}
	if _, err := os.Stat(configPath); err == nil {
		file, err := os.Open(configPath)
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
	return settings
}
