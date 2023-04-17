package cli

import (
	"os"

	"github.com/abiosoft/ishell"
)

func (cli *Cli) registerNounsHistory() {
	cmds := []noun{
		{"show", "history", showHistoryHelp, cli.showHistory},

		{"remove", "history", removeHistoryHelp, cli.removeHistory},

		{"set", "history", setHistoryHelp, cli.setHistory},
	}
	cli.registerNouns(cmds)
}

const showHistoryHelp = "show history settings"

func (cli *Cli) showHistory(c *ishell.Context) {
	c.Printf("%-25s\n", "History Settings")
	if cli.sh.histFile != "" {
		c.Printf(" %-24s : %-12s\n", "Status", "ON")
		c.Printf(" %-24s : %-12s\n", "File", cli.sh.histFile)
	} else {
		c.Printf(" %-24s : %-12s\n", "Status", "OFF")
	}
	c.Println("-------------------------------------------------")
}

const removeHistoryHelp = "remove history"

func (cli *Cli) removeHistory(c *ishell.Context) {
	if cli.sh.histFile != "" {
		cli.sh.shell.SetHistoryPath("")
		os.Remove(cli.sh.histFile)
		cli.sh.shell.SetHistoryPath(cli.sh.histFile)
	}
}

const setHistoryHelp = "set history on|off|file <filename>"

func (cli *Cli) setHistory(c *ishell.Context) {
	argLen := len(c.Args)
	if argLen < 1 {
		c.Println("Wrong input.", setHistoryHelp)
		return
	}
	switch c.Args[0] {
	case "on":
		cli.sh.shell.SetHistoryPath(cli.sh.histFile)
	case "off":
		cli.sh.shell.SetHistoryPath("")
	case "file":
		if argLen < 2 {
			c.Println("Wrong input.", setHistoryHelp)
			return
		}
		cli.sh.shell.SetHistoryPath(c.Args[1])
		cli.sh.histFile = c.Args[1]
	}
}
