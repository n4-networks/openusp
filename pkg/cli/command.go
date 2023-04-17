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
