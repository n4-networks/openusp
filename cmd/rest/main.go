package main

import (
	"log"

	"github.com/n4-networks/usp/pkg/rest"
)

func main() {
	log.SetFlags(log.Lshortfile)

	re := &rest.Rest{}
	log.Println("Initializing API Server...")
	if err := re.Init(); err != nil {
		log.Println("Error:", err)
	}
	log.Println("Starting RESt Server...")
	if err := re.Server(); err != nil {
		log.Println("Error: RESt Server is exiting...Err:", err)
	}
}
