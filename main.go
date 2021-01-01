package main

import (
	"github.com/darmiel/whgoxy/whgoxy"
	"log"
)

func main() {
	opt, err := whgoxy.Parse()
	if err != nil {
		log.Fatalln("Error while parsing:", err.Error())
		return
	}
	whgoxy.New(opt)
}
