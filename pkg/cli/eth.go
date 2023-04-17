package cli

import "github.com/abiosoft/ishell"

func (cli *Cli) registerNounsEth() {
	cmds := []noun{
		{"show", "eth", showEthHelp, cli.showEth},

		{"set", "eth", setEthHelp, cli.setEth},
		{"setcfg", "eth", setCfgEthHelp, cli.setCfgEth},

		//{"unset", "eth", unsetCfgEthHelp, cli.unsetCfgEth},
	}
	cli.registerNouns(cmds)
}

const showEthHelp = "show eth <intf|link|vlanterm|rmonstats|wol|lag> [id] <stats>"

func (cli *Cli) showEth(c *ishell.Context) {
	if err := cli.showParams(c, "eth"); err != nil {
		c.Println(err)
		cli.lastCmdErr = err
	}
}

const setEthHelp = "set eth <intf|link|vlanterm|wol|lag> <id> <param> <value>"

func (cli *Cli) setEth(c *ishell.Context) {
	if err := cli.setParam(c.Args, "eth"); err != nil {
		c.Println(err)
		c.Println(setEthHelp)
	}
}

const setCfgEthHelp = "setcfg eth <intf|link|vlanterm|wol|lag> <id> <param> <value>"

func (cli *Cli) setCfgEth(c *ishell.Context) {
	if err := cli.setCfgParam(c.Args, "eth"); err != nil {
		c.Println(err)
		c.Println(setCfgEthHelp)
	}
}
