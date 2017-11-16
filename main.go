package main

import (
	"io/ioutil"
	"fmt"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type server struct {
	greetingMsg	string
	allMsg		[]string
	name		string
	userNumber	int
}

func main() {
	dat, err := ioutil.ReadFile("./program-version")
	check(err)
	fmt.Print(string(dat))

	server := server{"Hello", []string{}, "Test-Server", 0}
	fmt.Print(server.greetingMsg)
}
