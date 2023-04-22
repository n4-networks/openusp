// Copyright 2023 N4-Networks.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
