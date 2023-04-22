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

import (
	"errors"

	"github.com/abiosoft/ishell"
)

func (cli *Cli) registerNounsDhcpv4() {
	cmds := []noun{
		{"add", "dhcpv4", addDhcpv4Help, cli.addDhcpv4},
		{"add.dhcpv4", "server", addDhcpv4ServerHelp, cli.addDhcpv4Server},
		{"add.dhcpv4.server", "pool", addDhcpv4ServerPoolHelp, cli.addDhcpv4ServerPool},
		/*
			{"add.dhcpv4.server.pool", "static", addDhcpv4ServerPoolStaticHelp, cli.addDhcpv4ServerPoolStatic},
			{"add.dhcpv4.server.pool", "option", addDhcpv4ServerPoolOptionHelp, cli.addDhcpv4ServerPoolOption},
			{"add.dhcpv4.server.pool", "client", addDhcpv4ServerPoolClientHelp, cli.addDhcpv4ServerPoolClient},
			{"add.dhcpv4.server.pool.client", "option", addDhcpv4ServerPoolClientOptionHelp, cli.addDhcpv4ServerPoolClientOption},
		*/

		{"add.dhcpv4", "relay", addDhcpv4RelayHelp, cli.addDhcpv4Relay},
		{"add.dhcpv4.relay", "forwarding", addDhcpv4RelayForwardingHelp, cli.addDhcpv4RelayForwarding},

		{"addcfg", "dhcpv4", addCfgDhcpv4Help, cli.addCfgDhcpv4},
		{"addcfg.dhcpv4", "server", addCfgDhcpv4ServerHelp, cli.addCfgDhcpv4Server},
		{"addcfg.dhcpv4.server", "pool", addCfgDhcpv4ServerPoolHelp, cli.addCfgDhcpv4ServerPool},

		{"addcfg.dhcpv4", "relay", addCfgDhcpv4RelayHelp, cli.addCfgDhcpv4Relay},
		{"addcfg.dhcpv4.relay", "forwarding", addCfgDhcpv4RelayForwardingHelp, cli.addCfgDhcpv4RelayForwarding},

		{"show", "dhcpv4", showDhcpv4Help, cli.showDhcpv4},
		{"showcfg", "dhcpv4", showCfgDhcpv4Help, cli.showCfgDhcpv4},

		{"remove", "dhcpv4", removeDhcpv4Help, cli.removeDhcpv4},
		{"removecfg", "dhcpv4", removeCfgDhcpv4Help, cli.removeCfgDhcpv4},

		{"set", "dhcpv4", setDhcpv4Help, cli.setDhcpv4},
		{"setcfg", "dhcpv4", setCfgDhcpv4Help, cli.setCfgDhcpv4},
	}
	cli.registerNouns(cmds)
}

const setCfgDhcpv4Help = "setcfg dhcpv4 server pool <id|name> <static|option|client> [id] <param> <value>"

func (cli *Cli) setCfgDhcpv4(c *ishell.Context) {
	if err := cli.setCfgParam(c.Args, "dhcpv4"); err != nil {
		c.Println(err)
		c.Println(setCfgDhcpv4Help)
	}
}

const setDhcpv4Help = "set dhcpv4 server pool <id|name> <static|option|client> [id] <param> <value>"

func (cli *Cli) setDhcpv4(c *ishell.Context) {
	if err := cli.setParam(c.Args, "dhcpv4"); err != nil {
		c.Println(err)
		c.Println(setDhcpv4Help)
	}
}

const removeDhcpv4Help = "remove dhcpv4 server pool <id|name> <static|option|client> <id|name>"

func (cli *Cli) removeDhcpv4(c *ishell.Context) {
	cli.removeInst(c, "dhcpv4")
}

const removeCfgDhcpv4Help = "removecfg dhcpv4 server pool <id|name> <static|option|client> <id|name>"

func (cli *Cli) removeCfgDhcpv4(c *ishell.Context) {
	cli.removeCfgInst(c, "dhcpv4")
}

const showDhcpv4Help = "show dhcpv4 server <pool> [id|name] <option|client|staticaddr> [id|name]"

func (cli *Cli) showDhcpv4(c *ishell.Context) {
	cli.showParams(c, "dhcpv4")
}

const showCfgDhcpv4Help = "showcfg dhcpv4 server <pool> [id|name] <option|client|staticaddr> [id|name]"

func (cli *Cli) showCfgDhcpv4(c *ishell.Context) {
	cli.showCfg(c, "dhcpv4")
}

const addDhcpv4ServerPoolHelp = "add dhcpv4 server pool alias <name|id> intf <id|name> minaddr <address> maxaddr <maxaddr> subnet <subnet> gw <address(es)> dns <address(es)>"

func (cli *Cli) addDhcpv4ServerPool(c *ishell.Context) {
	if err := cli.checkDefault(); err != nil {
		c.Println(err)
		return
	}
	dest := destMtp

	instInfo, err := parseAddDhcpv4ServerPoolArgs(c.Args, dest)
	if err != nil {
		c.Println(err)
		c.Printf(addDhcpv4ServerPoolHelp)
		return
	}
	printAddInstInfo(instInfo)

	if instId, err := cli.addInst(instInfo, dest); err != nil {
		c.Println(err)
		return
	} else {
		c.Println("Instance created with the path:", instId)
		// Enable the server
		params := map[string]string{"Enable": "true"}
		if err := cli.restSetParams(instId, params); err != nil {
			c.Println("Error in activating DHCPv4 server pool:", err)
			return
		} else {
			c.Println("Enabled DHCP server pool")
		}
	}
}

const addCfgDhcpv4ServerPoolHelp = "addcfg dhcpv4 server pool alias <name|id> intf <id|name> minaddr <address> maxaddr <maxaddr> subnet <subnet> gw <address(es)> dns <address(es)>"

func (cli *Cli) addCfgDhcpv4ServerPool(c *ishell.Context) {
	if err := cli.checkDefault(); err != nil {
		c.Println(err)
		return
	}
	dest := destDb

	instInfo, err := parseAddDhcpv4ServerPoolArgs(c.Args, dest)
	if err != nil {
		c.Println(err)
		c.Printf(addDhcpv4ServerPoolHelp)
		return
	}

	if _, err := cli.addInst(instInfo, dest); err != nil {
		c.Println(err)
		return
	}
	c.Println("add Instance has been  added to cfg store")

	// Set Enable to true
	/*
		cfgInstPath := instInfo.path + "[Alias==\"" + instInfo.params["Alias"] + "\"]."
		setParams := make(map[string]string)
		setParams["Enable"] = "true"
			if err := cli.dbWriteCfgParamNode(cfgInstPath, setParams); err != nil {
				c.Println("Error:", err)
			}
	*/
}

func parseAddDhcpv4ServerPoolArgs(args []string, dest destType) (*addInstInfo, error) {
	// Receive arguments
	argLen := len(args)
	if argLen < 12 {
		return nil, errors.New("Wrong input")
	}
	am, _ := getMapFromArgs(args) // argMap

	// Validate Inputs and form param map

	params := make(map[string]string)

	if alias, ok := am["alias"]; ok {
		params["Alias"] = alias
	} else if dest == destDb {
		return nil, errors.New("Alias is must for cfg")
	}

	/*
		if isDigit(am["intf"]) {
			params["Interface"] = "Device.IP.Interface." + am["intf"] + "."
		} else {
			params["Interface"] = "Device.IP.Interface." + "[Alias==\"" + am["intf"] + "\"]."
		}
	*/

	if minaddr, ok := am["minaddr"]; ok {
		params["MinAddress"] = minaddr
	} else {
		return nil, errors.New("minaddr is not provided")
	}

	if maxaddr, ok := am["maxaddr"]; ok {
		params["MaxAddress"] = maxaddr
	} else {
		return nil, errors.New("maxaddr is not provided")
	}

	if subnet, ok := am["subnet"]; ok {
		params["SubnetMask"] = subnet
	} else {
		return nil, errors.New("subnet is not provided")
	}

	params["IPRouters"] = am["gw"]
	if gw, ok := am["gw"]; ok {
		params["IPRouters"] = gw
	} else {
		return nil, errors.New("gw is not provided")
	}
	params["DNSServers"] = am["dns"]
	params["Enable"] = "true"

	parent := "Device.DHCPv4.Server."
	path := parent + "Pool."
	info := &addInstInfo{
		path:   path,
		parent: parent,
		params: params,
	}
	return info, nil
}

const addDhcpv4RelayForwardingHelp = "add dhcpv4 relay forwarding alias <name|id> intf <id|name> server-ip <address>"

func (cli *Cli) addDhcpv4RelayForwarding(c *ishell.Context) {
	if err := cli.checkDefault(); err != nil {
		c.Println(err)
		return
	}
	dest := destMtp

	instInfo, err := parseAddDhcpv4RelayForwardingArgs(c.Args, dest)
	if err != nil {
		c.Println(err)
		c.Printf(addDhcpv4RelayForwardingHelp)
		return
	}
	printAddInstInfo(instInfo)

	if instId, err := cli.addInst(instInfo, dest); err != nil {
		c.Println(err)
		return
	} else {
		c.Println("Instance created with the path:", instId)
		// Enable the server
		params := map[string]string{"Enable": "true"}
		if err := cli.restSetParams(instId, params); err != nil {
			c.Println("Error in activating AccessPoint:", err)
			return
		} else {
			c.Println("Enabled dhcp server pool")
		}
	}
}

const addCfgDhcpv4RelayForwardingHelp = "addcfg dhcpv4 relay forwarding alias <name|id> intf <id|name> server-ip <address>"

func (cli *Cli) addCfgDhcpv4RelayForwarding(c *ishell.Context) {
	if err := cli.checkDefault(); err != nil {
		c.Println(err)
		return
	}
	dest := destDb

	instInfo, err := parseAddDhcpv4RelayForwardingArgs(c.Args, dest)
	if err != nil {
		c.Println(err)
		c.Printf(addDhcpv4RelayForwardingHelp)
		return
	}

	if _, err := cli.addInst(instInfo, dest); err != nil {
		c.Println(err)
		return
	}
	c.Println("add Instance has been  added to cfg store")

	// Set Enable to true
	/*
		cfgInstPath := "Device.DHCPv4.Relay.Forwarding." + "[Alias==\"" + instInfo.params["Alias"] + "\"]."
		setParams := make(map[string]string)
		setParams["Enable"] = "true"
			if err := cli.dbWriteCfgParamNode(cfgInstPath, setParams); err != nil {
				c.Println("Error:", err)
			}
	*/
}

func parseAddDhcpv4RelayForwardingArgs(args []string, dest destType) (*addInstInfo, error) {
	// Receive arguments
	argLen := len(args)
	if argLen < 4 {
		return nil, errors.New("Wrong input")
	}
	am, _ := getMapFromArgs(args) // argMap

	// Validate Inputs and form param map

	params := make(map[string]string)

	if alias, ok := am["alias"]; ok {
		params["Alias"] = alias
	} else if dest == destDb {
		return nil, errors.New("Alias is must for cfg")
	}

	if isDigit(am["intf"]) {
		params["Interface"] = "Device.IP.Interface." + am["intf"] + "."
	} else {
		params["Interface"] = "Device.IP.Interface." + "[Alias==\"" + am["intf"] + "\"]."
	}

	if serverAddr, ok := am["server-ip"]; ok {
		params["DHCPServerIPAddress"] = serverAddr
	} else {
		return nil, errors.New("DHCP Server IP is not provided")
	}

	//params["Enable"] = "true"

	parent := "Device.DHCPv4.Relay."
	path := parent + "Forwarding."
	info := &addInstInfo{
		path:   path,
		parent: parent,
		params: params,
	}
	return info, nil
}

const addDhcpv4Help = "add dhcpv4 server..."

func (cli *Cli) addDhcpv4(c *ishell.Context) {
	c.Printf(addDhcpv4Help)
}

const addDhcpv4ServerHelp = "add dhcpv4 server pool..."

func (cli *Cli) addDhcpv4Server(c *ishell.Context) {
	c.Printf(addDhcpv4ServerHelp)
}

const addDhcpv4RelayHelp = "add dhcpv4 relay forwarding..."

func (cli *Cli) addDhcpv4Relay(c *ishell.Context) {
	c.Printf(addDhcpv4RelayHelp)
}

const addCfgDhcpv4Help = "addcfg dhcpv4 server..."

func (cli *Cli) addCfgDhcpv4(c *ishell.Context) {
	c.Printf(addCfgDhcpv4Help)
}

const addCfgDhcpv4RelayHelp = "addcfg dhcpv4 relay..."

func (cli *Cli) addCfgDhcpv4Relay(c *ishell.Context) {
	c.Printf(addCfgDhcpv4RelayHelp)
}

const addCfgDhcpv4ServerHelp = "addcfg dhcpv4 server pool..."

func (cli *Cli) addCfgDhcpv4Server(c *ishell.Context) {
	c.Printf(addCfgDhcpv4ServerHelp)
}
