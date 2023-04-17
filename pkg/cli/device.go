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
