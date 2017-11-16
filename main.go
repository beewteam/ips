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

func main() {
    dat, err := ioutil.ReadFile("./program-version")
    check(err)
    fmt.Print(string(dat))
}
