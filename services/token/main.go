package main

import (
	"log"

	"ivpn.net/auth/services/token/client"
)

func main() {
	hsm := client.NewMockHSMClient()
	server := New(hsm)
	err := server.Start()
	if err != nil {
		log.Println(err)
	}
}
