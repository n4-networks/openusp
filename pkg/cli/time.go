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
