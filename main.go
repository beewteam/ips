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

func (s *server) addMessage(msg string) {
	s.allMsg = append(s.allMsg, msg)
}

func (s server) getMessage(n int) string {
	return s.allMsg[n]
}

func main() {
	dat, err := ioutil.ReadFile("./program-version")
	check(err)
	fmt.Print(string(dat))

	server := server{"Hello", []string{}, "Test-Server", 0}
	fmt.Printf("%s\n", server.greetingMsg)
	
	server.addMessage("new msg")
	fmt.Printf("New msg: %s\n", server.getMessage(0))
}
