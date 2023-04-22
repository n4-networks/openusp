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

func (cli *Cli) registerNounsInstance() {
	cmds := []noun{
		{"add", "instance", addInstanceHelp, cli.addInstanceCmd},
		{"update", "instance", updateInstanceHelp, cli.updateInstance},
		{"show", "instance", showInstanceHelp, cli.showInstance},
	}
	cli.registerNouns(cmds)
}

const showInstanceHelp = "show instance <path|objname>"

func (cli *Cli) showInstance(c *ishell.Context) {
	var err error
	if err = cli.checkDefault(); err != nil {
		c.Println(err)
		return
	}
	path := getPath(c.Args)

	instances, err := cli.restReadInstances(path)
	if err != nil {
		log.Println("restErr:", err)
		return
	}
	for _, instance := range instances {
		c.Printf("%-25s : %-12s\n", "Instance Path", instance.Path)
		c.Printf("%-25s :\n", "Unique Keys")
		for key, value := range instance.UniqueKeys {
			c.Printf(" %-24s : %-12s\n", key, value)
		}
		c.Println("-------------------------------------------------")
	}
}

const updateInstanceHelp = "update instance [path]"

func (cli *Cli) updateInstance(c *ishell.Context) {
	var err error
	if err = cli.checkDefault(); err != nil {
		c.Println(err)
		return
	}
	// update instances Device. 1
	path := getPath(c.Args)

	if err = cli.restUpdateInstances(path); err != nil {
		c.Println(err)
	}
}

const addInstanceHelp = "add instance <path> <param name> <value>"

func (cli *Cli) addInstanceCmd(c *ishell.Context) {
	var err error
	if err = cli.checkDefault(); err != nil {
		c.Println(err)
		return
	}
	if len(c.Args) < 3 {
		c.Println("Wrong input. add instance <path> <param name> <value>)")
		return
	}
	path := c.Args[0]
	name := c.Args[1]
	value := c.Args[2]

	// Add Instance
	params := make(map[string]string)
	params[name] = value
	instInfo := &addInstInfo{
		parent: "Device.",
		path:   path,
		params: params,
	}
	instPath, err := cli.addInst(instInfo, destMtp)
	if err != nil {
		c.Println(err)
		return
	}
	c.Println("Instance created with path:", instPath)
}

/*
func (cli *Cli) GetInstancePathByAlias(name string) (string, error) {
	inst, err := cli.dbGetInstanceByAlias(name)
	if err != nil {
		return "", err
	}
	return inst.Path, nil
}
*/
