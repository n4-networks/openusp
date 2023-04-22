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

func (cli *Cli) registerNounsDevice() {
	cmds := []noun{
		{"show", "device", showDeviceHelp, cli.showDevice},
		{"operate", "device", operateDeviceHelp, cli.operateDevice},
	}
	cli.registerNouns(cmds)
}

const showDeviceHelp = "show device"

func (cli *Cli) showDevice(c *ishell.Context) {
	cli.showParams(c, "device")
}

const operateDeviceHelp = "operate device <reboot|factory-reset|self-test-diag>"

func (cli *Cli) operateDevice(c *ishell.Context) {
	cli.operateCmd(c, "device")
}
