package main

import (
	"log"

	"ivpn.net/auth/services/token/client"
	"ivpn.net/auth/services/token/service"
)

func main() {
	hsm := client.NewMockHSMClient()
	server := service.NewServer(hsm)
	err := server.Start()
	if err != nil {
		log.Println(err)
	}
}
