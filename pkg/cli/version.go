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
