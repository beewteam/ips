package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type UserConfig struct {
	UserData struct {
		Username string
		FullName string
		Nickname string
		Password string
		Chat     string
	}
	ServerData struct {
		Hostname string
		Port     string
	}
}

func ParseConfigFile(configPath string) (config UserConfig, err error) {
	config = UserConfig{}
	if _, err = os.Stat(configPath); err == nil {
		var file *os.File

		file, err = os.Open(configPath)
		defer file.Close()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		err = json.NewDecoder(file).Decode(&config)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}
	return
}
