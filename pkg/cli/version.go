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
	"strconv"

	"github.com/abiosoft/ishell"
)

type verType struct {
	major uint
	minor uint
}

var buildtime string

var (
	ver verType = verType{
		major: 1,
		minor: 0,
	}
)

func (cli *Cli) registerNounsVersion() {
	cmds := []noun{
		{"show", "version", showVersionHelp, cli.showVersion},
	}
	cli.registerNouns(cmds)
}

const showVersionHelp = "show version"

func (cli *Cli) showVersion(c *ishell.Context) {
	c.Printf("%-25s : %-12s\n", "CLI Version", getVer())
}

func getVer() string {

	v := strconv.FormatUint(uint64(ver.major), 10) + "." + strconv.FormatUint(uint64(ver.minor), 10)
	if buildtime != "" {
		v = v + "." + buildtime
	}
	return v
}
