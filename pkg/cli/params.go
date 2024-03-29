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
	"log"

	"github.com/abiosoft/ishell"
)

func (cli *Cli) registerNounsParam() {
	cmds := []noun{
		{"show", "param", showParamHelp, cli.showParamCmd},
		{"update", "param", updateParamHelp, cli.updateParamCmd},
		{"set", "param", setParamHelp, cli.setParamCmd},
	}
	cli.registerNouns(cmds)
}

const setParamHelp = "set param <path> <name> <value>"

func (cli *Cli) setParamCmd(c *ishell.Context) {
	var err error
	if err = cli.checkDefault(); err != nil {
		c.Println(err)
		return
	}
	if len(c.Args) < 3 {
		c.Println("Wrong input. Minimum 3 parameters are required")
		return
	}
	path := getPath(c.Args)

	name := c.Args[1]
	value := c.Args[2]

	params := map[string]string{name: value}
	if err = cli.restSetParams(path, params); err != nil {
		c.Println(err)
		return
	}
	c.Println("Set param for:", path)
}

const updateParamHelp = "update param <path>"

func (cli *Cli) updateParamCmd(c *ishell.Context) {
	var err error
	if err = cli.checkDefault(); err != nil {
		c.Println(err)
		return
	}
	path := getPath(c.Args)

	if err = cli.restUpdateParams(path); err != nil {
		c.Println(err)
		return
	}
	c.Println("Updated param for:", path)
}

const showParamHelp = "show param <path>"

func (cli *Cli) showParamCmd(c *ishell.Context) {
	var err error
	if err = cli.checkDefault(); err != nil {
		c.Println(err)
		return
	}

	path := getPath(c.Args)

	objParams, err := cli.restReadParams(path)
	if err != nil {
		log.Println("Err:", err)
		return
	}
	for _, obj := range objParams {
		c.Printf("%-25s : %-12s\n", "Object Path", obj.Path)
		for _, p := range obj.Params {
			c.Printf(" %-24s : %-12s\n", p.Name, p.Value)
		}
		c.Println("-------------------------------------------------")
	}
}
