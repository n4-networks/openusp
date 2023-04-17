package cli

import (
	"errors"

	"github.com/abiosoft/ishell"
)

func (cli *Cli) registerNounsBridging() {
	cmds := []noun{
		{"add", "bridging", addBridgingHelp, cli.addBridging},
		{"add.bridging", "port", addBridgingPortHelp, cli.addBridgingPort},
		{"add.bridging", "bridge", addBridgingBridgeHelp, cli.addBridgingBridge},

		{"addcfg", "bridging", addCfgBridgingHelp, cli.addCfgBridging},
		{"addcfg.bridging", "port", addCfgBridgingPortHelp, cli.addCfgBridgingPort},
		{"addcfg.bridging", "bridge", addCfgBridgingBridgeHelp, cli.addCfgBridgingBridge},

		{"show", "bridging", showBridgingHelp, cli.showBridging},
		{"showcfg", "bridging", showCfgBridgingHelp, cli.showCfgBridging},

		{"remove", "bridging", removeBridgingHelp, cli.removeBridging},
		{"removecfg", "bridging", removeCfgBridgingHelp, cli.removeCfgBridging},

		{"set", "bridging", setBridgingHelp, cli.setBridging},
		{"setcfg", "bridging", setCfgBridgingHelp, cli.setCfgBridging},
	}
	cli.registerNouns(cmds)
}

const setBridgingHelp = "set bridging <bridge|filter|provbridge|> [id] <port|vlan|vlanport> [id] <param> <value>"

func (cli *Cli) setBridging(c *ishell.Context) {
	if err := cli.setParam(c.Args, "bridging"); err != nil {
		c.Println(err)
		c.Println(setBridgingHelp)
	}
}

const setCfgBridgingHelp = "setcfg bridging <bridge|filter|provbridge|> [id] <port|vlan|vlanport> [id] <param> <value>"

func (cli *Cli) setCfgBridging(c *ishell.Context) {
	if err := cli.setCfgParam(c.Args, "ip"); err != nil {
		c.Println(err)
		c.Println(setCfgBridgingHelp)
	}
}

const removeBridgingHelp = "remove bridging <bridge|filter> <id|name> [port] [id|name]"

func (cli *Cli) removeBridging(c *ishell.Context) {
	cli.removeInst(c, "bridging")
}

const removeCfgBridgingHelp = "removecfg bridging <bridge|filter> <id|name> [port] [id|name]"

func (cli *Cli) removeCfgBridging(c *ishell.Context) {
	cli.removeCfgInst(c, "bridging")
}

const showBridgingHelp = "show bridging <bridge|filter|provbridge|> [id] <port|vlan|vlanport> [id] <stats|pricodepoint>"

func (cli *Cli) showBridging(c *ishell.Context) {
	cli.showParams(c, "bridging")
}

const showCfgBridgingHelp = "showcfg bridging <bridge|filter|provbridge|> [id] <port|vlan|vlanport> [id] <stats|pricodepoint>"

func (cli *Cli) showCfgBridging(c *ishell.Context) {
	cli.showCfg(c, "bridging")
}

const addBridgingBridgeHelp = "add bridging bridge alias <string> name <name>"

func (cli *Cli) addBridgingBridge(c *ishell.Context) {
	if err := cli.checkDefault(); err != nil {
		c.Println(err)
		return
	}
	instInfo, err := parseAddBridgingBridgeArgs(c.Args)
	if err != nil {
		c.Println(err)
		return
	}
	if instId, err := cli.addInst(instInfo, destMtp); err != nil {
		c.Println(err)
		return
	} else {
		c.Println("Bridge created with path:", instId)
	}
}

const addCfgBridgingBridgeHelp = "addcfg bridging bridge alias <string> name <name>"

func (cli *Cli) addCfgBridgingBridge(c *ishell.Context) {
	if err := cli.checkCfgDevSet(); err != nil {
		c.Println(err)
		//c.Println("Use ", SetCfgDevInfoHelp)
		return
	}
	instInfo, err := parseAddBridgingBridgeArgs(c.Args)
	if err != nil {
		c.Println(err)
		return
	}
	if _, err := cli.addInst(instInfo, destDb); err != nil {
		c.Println(err)
		return
	}
	c.Println("Bridging Instance Bridge added to cfg datastore")
}

func parseAddBridgingBridgeArgs(args []string) (*addInstInfo, error) {
	argLen := len(args)
	if argLen < 2 {
		return nil, errors.New("Wrong input.")
	}
	am, _ := getMapFromArgs(args) // argMap

	// Add Instance
	params := make(map[string]string)
	parent := "Device.Bridging."
	path := parent + "Bridge."

	if name, ok := am["name"]; ok {
		params["Name"] = name
	} else {
		return nil, errors.New("No name found")
	}

	if alias, ok := am["alias"]; ok {
		params["Alias"] = alias
	} else {
		return nil, errors.New("No alias found")
	}

	info := &addInstInfo{
		path:   path,
		parent: parent,
		params: params,
	}
	return info, nil
}

const addCfgBridgingHelp = "addcfg bridging port|bridge..."

func (cli *Cli) addCfgBridging(c *ishell.Context) {
	c.Println(addCfgBridgingHelp)
}

const addBridgingPortHelp = "add bridging port alias <string> bridge <id|name> intftype <eth|ppp> id <id|name>"

func (cli *Cli) addBridgingPort(c *ishell.Context) {
	if err := cli.checkDefault(); err != nil {
		c.Println(err)
		return
	}
	instInfo, err := parseAddBridgingPortArgs(c.Args)
	if err != nil {
		c.Println(err)
		return
	}
	if instId, err := cli.addInst(instInfo, destMtp); err != nil {
		c.Println(err)
		return
	} else {
		c.Println("Bridge port created with path:", instId)
	}
}

const addCfgBridgingPortHelp = "addcfg bridging port alias <string> bridge <id|name> intftype <eth|ppp> id <id|name>"

func (cli *Cli) addCfgBridgingPort(c *ishell.Context) {
	if err := cli.checkCfgDevSet(); err != nil {
		c.Println(err)
		//c.Println("Use ", SetCfgDevInfoHelp)
		return
	}
	instInfo, err := parseAddBridgingPortArgs(c.Args)
	if err != nil {
		c.Println(err)
		return
	}
	if _, err := cli.addInst(instInfo, destDb); err != nil {
		c.Println(err)
		return
	}
	c.Println("Bridging Instance Port added to cfg datastore")
}

func parseAddBridgingPortArgs(args []string) (*addInstInfo, error) {
	argLen := len(args)
	if argLen < 8 {
		return nil, errors.New("Wrong input.")
	}

	am, _ := getMapFromArgs(args) // argMap

	// Validate Inputs

	// Add Instance
	params := make(map[string]string)
	parent := "Device.Bridging.Bridge."
	var path string
	if isDigit(am["bridge"]) {
		path = parent + am["bridge"] + ".Port."
	} else {
		path = parent + "[Alias==\"" + am["bridge"] + "\"].Port"
	}

	if alias, ok := am["alias"]; ok {
		params["Alias"] = alias
	} else {
		return nil, errors.New("Wrong interface type")
	}

	var lowerLayerPath string
	switch am["intftype"] {
	case "eth":
		lowerLayerPath = "Device.Ethernet.Interface."
	case "ppp":
		lowerLayerPath = "Device.Bridging.PPP."
	default:
		return nil, errors.New("Wrong interface type")
	}

	if isDigit(am["intftype_id"]) {
		lowerLayerPath = lowerLayerPath + am["intftype_id"] + "."
	} else {
		lowerLayerPath = lowerLayerPath + "[Alias==\"" + am["intftype_id"] + "\"]."
	}
	params["LowerLayers"] = lowerLayerPath

	info := &addInstInfo{
		path:   path,
		parent: parent,
		params: params,
	}
	return info, nil
}

const addBridgingHelp = "add bridging port|bridge..."

func (cli *Cli) addBridging(c *ishell.Context) {
	c.Println(addBridgingHelp)
}
