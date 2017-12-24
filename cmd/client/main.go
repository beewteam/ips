package main

import (
	"flag"
	"fmt"

	"github.com/fatih/color"
)

var configFile string

func main() {
	fmt.Println("Client version: " + color.GreenString(VERSION))

	flag.StringVar(&configFile, "c", "UserConfigs.json", "path to config file in json, default==UserConfigs.json in pwd")
	flag.Parse()

	client := NewClient(configFile)
	if client == nil {
		panic(color.RedString("client") + ": cannot create client")
	}

	if err := client.Run(); err != nil {
		panic(color.RedString("client") + ": " + err.Error())
	}
}
