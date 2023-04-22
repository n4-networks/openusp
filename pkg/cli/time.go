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

import "github.com/abiosoft/ishell"

func (cli *Cli) registerNounsTime() {
	cmds := []noun{
		{"show", "time", showTimeHelp, cli.showTime},
		{"set", "ip", setIpHelp, cli.setIp},
	}
	cli.registerNouns(cmds)
}

const showTimeHelp = "show time"

func (cli *Cli) showTime(c *ishell.Context) {
	if err := cli.restUpdateParams("Device.Time."); err != nil {
		c.Println(err)
		return
	}
	cli.showParams(c, "time")
}

const setTimeHelp = "set time <param> <value>"

func (cli *Cli) setTime(c *ishell.Context) {
	if err := cli.setParam(c.Args, "time"); err != nil {
		c.Println(err)
		c.Println(setTimeHelp)
	}
}
