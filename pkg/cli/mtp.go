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

package cli

import (
	"github.com/abiosoft/ishell"
)

var resultStr = map[bool]string{
	true:  "Passed",
	false: "Failed",
}

func (cli *Cli) registerNounsMtp() {
	cmds := []noun{
		{"reconnect", "mtp", connectMtpHelp, cli.reconnectMtp},
		{"show", "mtp", showMtpHelp, cli.showMtp},
	}
	cli.registerNouns(cmds)
}

const showMtpHelp = "show mtp"

func (cli *Cli) showMtp(c *ishell.Context) {
	c.Printf("%-25s\n", "MTP Status")
	if info, err := cli.restMtpGetInfo(); err == nil {
		c.Printf(" %-24s : %-12s\n", "Version", info.Version)
	}
	c.Println("-------------------------------------------------")
}

const connectMtpHelp = "reconnect mtp"

func (cli *Cli) reconnectMtp(c *ishell.Context) {

	c.Println("RESt server is reconnecting to MTP")
	if err := cli.restReconnectMtp(); err != nil {
		c.Println("Error:", err)
		return
	}
	c.Println("Success")
}

type MtpInfo struct {
	Version string `json:"version"`
}
