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

package main

import (
	"flag"
	"log"

	"github.com/n4-networks/openusp/pkg/cntlr"
)

func main() {

	var cliMode bool
	var err error

	flag.BoolVar(&cliMode, "c", false, "run with cli")
	flag.Parse()

	var c cntlr.Cntlr
	err = c.Init()
	if err != nil {
		log.Println("Could not initialize Mtp, err:", err)
		return
	}

	var status chan int32
	status, err = c.Run()
	if err != nil {
		log.Println("Error in running server, exiting...")
		return
	}
	if cliMode {
		log.Println("Going to wait here")
		go c.WaitForExit(status)
		c.Cli()
	} else {
		log.Printf("Started MTP server successfully")
		c.WaitForExit(status)
	}
}
