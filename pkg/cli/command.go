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

func (cli *Cli) registerNounsCommand() {
	cmds := []noun{
		{"operate", "command", operateCommandHelp, cli.operateCommand},
	}
	cli.registerNouns(cmds)
}

const operateCommandHelp = "operate command <path> <input-params>"

func (cli *Cli) operateCommand(c *ishell.Context) {
	var err error
	if err = cli.checkDefault(); err != nil {
		c.Println(err)
		return
	}
	if len(c.Args) <= 0 {
		c.Println("Please provide command path")
		return
	}
	path := c.Args[0]
	inputArgs, err := getMapFromArgs(c.Args[1:])
	if err != nil {
		c.Println(err)
		return
	}

	if err := cli.restOperateCmd(path, inputArgs); err != nil {
		c.Println(err)
		return
	}
	c.Println("Cmd executed successfully on the device")
}
