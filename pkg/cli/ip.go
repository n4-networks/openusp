package cli

import (
	"errors"

	"github.com/abiosoft/ishell"
)

func (cli *Cli) registerNounsIp() {
	cmds := []noun{
		{"add", "ip", addIpHelp, cli.addIp},
		{"add.ip", "intf", addIpIntfHelp, cli.addIpIntf},
		{"add.ip", "addr", addIpAddrHelp, cli.addIpAddr},

		{"addcfg", "ip", addCfgIpHelp, cli.addCfgIp},
		{"addcfg.ip", "intf", addCfgIpIntfHelp, cli.addCfgIpIntf},
		{"addcfg.ip", "addr", addCfgIpAddrHelp, cli.addCfgIpAddr},

		{"show", "ip", showIpHelp, cli.showIp},
		{"showcfg", "ip", showCfgIpHelp, cli.showCfgIp},

		{"remove", "ip", removeIpHelp, cli.removeIp},
		{"removecfg", "ip", removeCfgIpHelp, cli.removeCfgIp},

		{"set", "ip", setIpHelp, cli.setIp},
		{"setcfg", "ip", setCfgIpHelp, cli.setCfgIp},
	}
	cli.registerNouns(cmds)
}

const setIpHelp = "set ip <intf|acport|twamp> <id> <ipv4addr|ipv6addr|ipv6prefix> [id] <param> <value>"

func (cli *Cli) setIp(c *ishell.Context) {
	if err := cli.setParam(c.Args, "ip"); err != nil {
		c.Println("Error:", err)
		cli.lastCmdErr = err
		return
	}
	cli.lastCmdErr = nil
	return
}

const setCfgIpHelp = "setcfg ip <intf|acport|twamp> <id> <ipv4addr|ipv6addr|ipv6prefix> [id] <param> <value>"

func (cli *Cli) setCfgIp(c *ishell.Context) {
	if err := cli.setCfgParam(c.Args, "ip"); err != nil {
		c.Println(err)
		c.Println(setCfgIpHelp)
		cli.lastCmdErr = err
		return
	}
	cli.lastCmdErr = nil
}

//const removeCfgIpHelp = "removecfg ip <intf|port> <id|name> <ipv4addr|ipv6addr> <id|name>"
const removeCfgIpHelp = "removecfg ip <intf|port> <id|name> <ipv4addr|ipv6addr> <id|name>"

func (cli Cli) removeCfgIp(c *ishell.Context) {
	if err := cli.removeCfgInst(c, "ip"); err != nil {
		c.Println(err)
		cli.lastCmdErr = err
		return
	}
	cli.lastCmdErr = nil
}

const removeIpHelp = "remove ip <intf|port> <id|name> <ipv4addr|ipv6addr> <id|name>"

func (cli *Cli) removeIp(c *ishell.Context) {
	if err := cli.removeInst(c, "ip"); err != nil {
		c.Println(err)
		cli.lastCmdErr = err
		return
	}
	cli.lastCmdErr = nil
	return
}

const addCfgIpAddrHelp = "addcfg ip addr alias <string> intf <id|name> type <v4|v6> mode <static|dhcp> address <address> subnet <subnet>"

func (cli *Cli) addCfgIpAddr(c *ishell.Context) {
	if err := cli.checkCfgDevSet(); err != nil {
		c.Println(err)
		//c.Println("Use ", SetCfgDevInfoHelp)
		cli.lastCmdErr = err
		return
	}
	dest := destDb
	instInfo, err := parseAddIpAddrArgs(c.Args, dest)
	if err != nil {
		c.Println(err)
		cli.lastCmdErr = err
		return
	}
	if _, err := cli.addInst(instInfo, dest); err != nil {
		c.Println(err)
		cli.lastCmdErr = err
		return
	}
	c.Println("IP Instance Addr added to cfg datastore")
	cli.lastCmdErr = nil
}

const addCfgIpIntfHelp = "addcfg ip intf alias <string> type <normal|lo> version <v4|v6|both> lowerlayer <eth|bridge> id <id|name>"

func (cli *Cli) addCfgIpIntf(c *ishell.Context) {
	if err := cli.checkCfgDevSet(); err != nil {
		c.Println(err)
		//c.Println("Use ", setCfgDevInfoHelp)
		cli.lastCmdErr = err
		return
	}
	dest := destDb
	instInfo, err := parseAddIpIntfArgs(c.Args, dest)
	if err != nil {
		c.Println(err)
		cli.lastCmdErr = err
		return
	}
	if _, err := cli.addInst(instInfo, dest); err != nil {
		c.Println(err)
		cli.lastCmdErr = err
		return
	}
	c.Println("IP interface instance added to cfg datastore")
	cli.lastCmdErr = nil
}

const addCfgIpHelp = "addcfg ip intf|addr|acport..."

func (cli *Cli) addCfgIp(c *ishell.Context) {
	c.Println(addCfgIpHelp)
}

const addIpAddrHelp = "add ip addr alias <string> intf <id|name> type <v4|v6> mode <static|dhcp> address <address> subnet <subnet>"

func (cli *Cli) addIpAddr(c *ishell.Context) {
	cli.lastCmdErr = errors.New("addIpAddr Error")
	if err := cli.checkDefault(); err != nil {
		cli.lastCmdErr = err
		c.Println(err)
		return
	}
	dest := destMtp
	instInfo, err := parseAddIpAddrArgs(c.Args, dest)
	if err != nil {
		cli.lastCmdErr = err
		c.Println(err)
		return
	}
	if instId, err := cli.addInst(instInfo, dest); err != nil {
		c.Println(err)
		cli.lastCmdErr = err
		return
	} else {
		c.Println("Instance created with the path:", instId)
	}
	cli.lastCmdErr = nil
}

func parseAddIpAddrArgs(args []string, dest destType) (*addInstInfo, error) {
	// Receive arguments
	argLen := len(args)
	if argLen < 10 {
		return nil, errors.New("Wrong input")
	}
	am, _ := getMapFromArgs(args) // argMap

	/* log.Printf("intf: %v, type: %v, mode: %v address: %v, subnet: %v\n",
	am["intf"], am["type"], am["mode"], am["address"], am["subnet"]) */

	params := make(map[string]string)
	switch am["mode"] {
	case "static":
		params["IPAddress"] = am["address"]
		params["SubnetMask"] = am["subnet"]
		params["AddressingType"] = "Static"
	case "dhcp":
		params["AddressingType"] = "Dhcp"
	default:
		return nil, errors.New("Mode not supported. Valid values: static|dhcp")
	}

	if alias, ok := am["alias"]; ok {
		params["Alias"] = alias
	} else if dest == destDb {
		return nil, errors.New("Alias is must for cfg")
	}
	// Validate Inputs

	// Create IPv4 address object
	var path string
	parent := "Device.IP.Interface."
	if isDigit(am["intf"]) {
		path = parent + am["intf"] + ".IPv4Address."
	} else {
		path = parent + "[Alias==\"" + am["intf"] + "\"].IPv4Address."
	}
	info := &addInstInfo{
		path:   path,
		parent: parent,
		params: params,
	}
	return info, nil
}

const showIpHelp = "show ip <intf|port|twmp> [id] <stats|ipv4addr|ipv6addr|ipv6prefix> [id]"

func (cli *Cli) showIp(c *ishell.Context) {

	if err := cli.showParams(c, "ip"); err != nil {
		c.Println(err)
		cli.lastCmdErr = err
		return
	}
}

const showCfgIpHelp = "showcfg ip <intf|port|twmp> [id] <stats|ipv4addr|ipv6addr|ipv6prefix> [id]"

func (cli *Cli) showCfgIp(c *ishell.Context) {
	if err := cli.showCfg(c, "ip"); err != nil {
		c.Println(err)
		cli.lastCmdErr = err
		return
	}
}

const addIpHelp = "add ip intf|addr|acport..."

func (cli *Cli) addIp(c *ishell.Context) {
	c.Printf(addIpHelp)
}

const addIpIntfHelp = "add ip intf alias <string> type <normal|lo> version <v4|v6|both> lowerlayer <eth|bridge> id <id|name>"

func (cli *Cli) addIpIntf(c *ishell.Context) {
	cli.lastCmdErr = errors.New("addIpAddr Error")
	if err := cli.checkDefault(); err != nil {
		c.Println(err)
		cli.lastCmdErr = err
		return
	}
	dest := destMtp
	instInfo, err := parseAddIpIntfArgs(c.Args, dest)
	if err != nil {
		c.Println(err)
		cli.lastCmdErr = err
		return
	}
	if instId, err := cli.addInst(instInfo, dest); err != nil {
		c.Println(err)
		cli.lastCmdErr = err
		return
	} else {
		c.Println("Instance created with the path:", instId)
	}
	cli.lastCmdErr = nil
}

func parseAddIpIntfArgs(args []string, dest destType) (*addInstInfo, error) {
	// Receive arguments
	argLen := len(args)
	var requiredArgs int

	if dest == destMtp {
		requiredArgs = 8
	} else {
		requiredArgs = 10
	}
	if argLen < requiredArgs {
		return nil, errors.New("Wrong input")
	}
	am, _ := getMapFromArgs(args) // argMap

	// "add ip intf <normal|lo> <v4|v6> <eth|bridge> <id|name>"
	/*log.Printf("type: %s, version: %s, lowerlayer: %s, lowerlayer_id: %s\n",
	am["type"], am["version"], am["lowerlayer"], am["lowerlayer_id"])*/

	// Validate Inputs and form param map

	params := make(map[string]string)

	var lowerlayer string
	switch am["lowerlayer"] {
	case "eth":
		lowerlayer = "Device.Ethernet.Interface."
	case "bridge":
		lowerlayer = "Device.Bridging.Bridge."
	default:
		return nil, errors.New("Unsupported lowerlayer. Valid values: eth|bridge")
	}

	if alias, ok := am["alias"]; ok {
		params["Alias"] = alias
	} else if dest == destDb {
		return nil, errors.New("Alias is must for cfg")
	}

	if isDigit(am["lowerlayer_id"]) {
		params["LowerLayers"] = lowerlayer + am["lowerlayer_id"] + "."
	} else {
		params["LowerLayers"] = lowerlayer + "[Alias==\"" + am["lowerlayer_id"] + "\"]."
	}

	/*
			switch am["version"] {
			case "v4":
				params["IPv4Enable"] = "true"
			case "v6":
				params["IPv6Enable"] = "true"
			case "both":
				params["IPv4Enable"] = "true"
				params["IPv6Enable"] = "true"
			}

		if am["type"] == "normal" {
			params["Type"] = "Normal"
		}
	*/

	parent := "Device.IP."
	path := parent + "Interface."
	info := &addInstInfo{
		path:   path,
		parent: parent,
		params: params,
	}
	return info, nil
}
