package main

import (
	"log"

	"github.com/n4-networks/openusp/pkg/apiserver"
)

func main() {
	log.SetFlags(log.Lshortfile)

	as := &apiserver.ApiServer{}
	log.Println("Initializing API Server...")
	if err := as.Init(); err != nil {
		log.Println("Error:", err)
	}
	log.Println("Starting API Server...")
	if err := as.Server(); err != nil {
		log.Println("Error: API Server is exiting...Err:", err)
	}
}
