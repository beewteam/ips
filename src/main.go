package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	dat, err := ioutil.ReadFile("../program-version")
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Print(string(dat))

	var client Client
	client.username = "Vperus"
	client.fullName = "Vperus"
	if !client.Connect("chat.freenode.net", "8000") {
		return
	}

	if !client.Login("Chicken") {
		return
	}

	client.JoinChannel("#go-nuts")

	for client.HandleData() {
	}

	client.Close()
}
