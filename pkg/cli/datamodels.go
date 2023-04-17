package cli

import (
	"log"

	"github.com/abiosoft/ishell"
)

func (cli *Cli) registerNounsDatamodel() {
	cmds := []noun{
		{"show", "datamodel", showDatamodelHelp, cli.showDatamodel},
		{"update", "datamodel", updateDatamodelHelp, cli.updateDatamodel},
	}
	cli.registerNouns(cmds)
}

const updateDatamodelHelp = "update datamodel path (must have . at the end)"

func (cli *Cli) updateDatamodel(c *ishell.Context) {
	var err error
	if err = cli.checkDefault(); err != nil {
		c.Println(err)
		return
	}
	path := getPath(c.Args)
	if err = cli.restUpdateDm(path); err != nil {
		c.Println(err)
	}
}

const showDatamodelHelp = "show datamodel <path>"

func (cli *Cli) showDatamodel(c *ishell.Context) {
	if err := cli.checkDefault(); err != nil {
		c.Println(err)
		return
	}

	path := getPath(c.Args)

	dmObjs, err := cli.restReadDm(path)
	if err != nil {
		log.Println("restErr:", err)
		return
	}
	log.Println("Fetched datamodel of:", path)
	log.Println("Len of datamodel objects:", len(dmObjs))
	for _, d := range dmObjs {
		c.Printf("path: %-24s, MultiInstance: %v Access: %v\n", d.Path, d.MultiInstance, d.Access)
		c.Printf("Commands:\n")
		for _, cmd := range d.Cmds {
			c.Printf("  %-24s, Input: %12s Output: %12s\n", cmd.Name, cmd.Inputs, cmd.Outputs)
		}
		c.Printf("Events:\n")
		for _, evt := range d.Events {
			c.Printf("  %-24s Args: %24s\n", evt.Name, evt.Args)
		}
		c.Printf("Params:\n")
		for _, param := range d.Params {
			c.Printf("  %-24s AccessType : %24s\n", param.Name, param.Access)
		}
	}
}
