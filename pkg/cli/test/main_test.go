package clitest

import (
	"errors"
	"log"
	"os"
	"testing"

	"github.com/n4-networks/openusp/pkg/cli"
)

var (
	cliH *cli.Cli
)

func TestMain(m *testing.M) {
	log.Println("Inside of TestMain")
	if err := setup(); err != nil {
		log.Println(err)
		os.Exit(-1)
	}
	log.Println("Initializtion of cli is completed")
	os.Exit(m.Run())
}

func setup() error {

	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	cliH = &cli.Cli{}
	if err := cliH.Init(); err != nil {
		log.Println("Error in initializing shell, exiting...:", err)
		return err
	}
	if !cliH.IsConnectedToDb() {
		return errors.New("Need valid DB connection to run tests")
	}
	if !cliH.IsConnectedToMtp() {
		return errors.New("Need valid MTP connection to run tests")
	}
	return nil
}
