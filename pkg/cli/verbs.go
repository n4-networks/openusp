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
	"errors"

	"github.com/abiosoft/ishell"
)

type verb struct {
	cmd  string
	objs []string
}

func (cli *Cli) registerVerbs() {
	verbs := []verb{
		{"add", []string{"bridging", "devinfo", "dhcpv4", "ip", "nat", "wifi", "instance"}},
		{"addcfg", []string{"bridging", "dhcpv4", "ip", "nat", "wifi"}},
		{"reconnect", []string{"db", "mtp", "stomp"}},
		{"operate", []string{"bridging", "command", "device", "devinfo", "ip", "wifi", "param", "instance"}},
		{"set", []string{"agent", "devinfo", "bridging", "history", "ip", "logging", "nat", "wifi", "param"}},
		{"setcfg", []string{"bridging", "devinfo", "ip", "nat", "wifi"}},
		{"show", []string{"agent", "bridging", "devinfo", "eth", "dhcpv4", "history", "ip", "logging", "nat", "nw", "wifi", "datamodel", "param", "instance", "version"}},
		{"showcfg", []string{"bridging", "devinfo", "eth", "dhcpv4", "ip", "nat", "wifi"}},
		{"remove", []string{"bridging", "db", "devinfo", "dhcpv4", "history", "ip", "nat", "stomp", "wifi", "param", "instance"}},
		{"removecfg", []string{"bridging", "dhcpv4", "ip", "nat", "wifi"}},
		{"update", []string{"bridging", "dhcpv4", "ip", "nat", "wifi", "datamodel", "param", "instance"}},
		{"unset", []string{"agent"}},
	}
	cli.addVerbCmds(verbs)
}

func (cli *Cli) addVerbCmds(cmds []verb) {
	for _, v := range cmds {
		help := v.cmd + "\t"
		for _, o := range v.objs {
			help = help + o + "|"
		}
		help = help + "..."
		cmd := &ishell.Cmd{
			Name: v.cmd,
			Help: help,
			Func: func(c *ishell.Context) {
				c.Println(help)
				cli.lastCmdErr = errors.New("Wrong command")
			},
		}
		cli.sh.shell.AddCmd(cmd)
		cli.sh.cmds[v.cmd] = cmd
	}
}
