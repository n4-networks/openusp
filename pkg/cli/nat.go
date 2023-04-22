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

func (cli *Cli) registerNounsNat() {
	cmds := []noun{
		{"add", "nat", addNatHelp, cli.addNat},
		{"add.nat", "intf-setting", addNatIntfSettingHelp, cli.addNatIntfSetting},
		{"add.nat", "port-mapping", addNatPortMappingHelp, cli.addNatPortMapping},

		{"addcfg", "nat", addCfgNatHelp, cli.addCfgNat},
		{"addcfg.nat", "intf-setting", addCfgNatIntfSettingHelp, cli.addCfgNatIntfSetting},
		{"addcfg.nat", "port-mapping", addCfgNatPortMappingHelp, cli.addCfgNatPortMapping},

		{"show", "nat", showNatHelp, cli.showNat},
		{"showcfg", "nat", showCfgNatHelp, cli.showCfgNat},

		{"remove", "nat", removeNatHelp, cli.removeNat},
		{"removecfg", "nat", removeCfgNatHelp, cli.removeCfgNat},

		{"set", "nat", setNatHelp, cli.setNat},
		{"setcfg", "nat", setCfgNatHelp, cli.setCfgNat},
	}
	cli.registerNouns(cmds)
}

const showNatHelp = "show nat <intf-setting|port-mapping> [id|name]"

func (cli *Cli) showNat(c *ishell.Context) {
	cli.showParams(c, "nat")
}

const showCfgNatHelp = "showcfg nat <intf-setting|port-mapping> [id|name]"

func (cli *Cli) showCfgNat(c *ishell.Context) {
	cli.showCfg(c, "nat")
}

const addNatHelp = "add nat intf-setting|port-mapping..."

func (cli *Cli) addNat(c *ishell.Context) {
	c.Printf(addNatHelp)
}

const addCfgNatHelp = "addcfg nat intf-setting|port-mapping..."

func (cli *Cli) addCfgNat(c *ishell.Context) {
	c.Println(addCfgNatHelp)
}

const addNatIntfSettingHelp = "add nat intf-setting alias <string> intf <id|name> tcp-timeout|udp-timeout <seconds> "

func (cli *Cli) addNatIntfSetting(c *ishell.Context) {
	if err := cli.checkDefault(); err != nil {
		c.Println(err)
		return
	}
	dest := destMtp
	instInfo, err := parseAddNatIntfSettingArgs(c.Args, dest)
	if err != nil {
		c.Println(err)
		return
	}
	if instId, err := cli.addInst(instInfo, dest); err != nil {
		c.Println(err)
		return
	} else {
		c.Println("Instance created with the path:", instId)
		// Enable
		params := map[string]string{"Enable": "true"}
		if err := cli.restSetParams(instId, params); err != nil {
			c.Println("Error in activating NAT Interface Setting:", err)
			return
		} else {
			c.Println("Enabled NAT Interface Setting")
		}
	}
}

const addCfgNatIntfSettingHelp = "addcfg nat intf-setting alias <string> intf <id|name> tcp-timeout|udp-timeout <seconds> "

func (cli *Cli) addCfgNatIntfSetting(c *ishell.Context) {
	if err := cli.checkDefault(); err != nil {
		c.Println(err)
		return
	}
	dest := destDb
	instInfo, err := parseAddNatIntfSettingArgs(c.Args, dest)
	if err != nil {
		c.Println(err)
		return
	}
	if instId, err := cli.addInst(instInfo, dest); err != nil {
		c.Println(err)
		return
	} else {
		c.Println("Instance created with the path:", instId)
	}

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

func parseAddNatIntfSettingArgs(args []string, dest destType) (*addInstInfo, error) {
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

	if tcpTimeout, ok := am["tcp-timeout"]; ok {
		params["TCPTransalationTimeout"] = tcpTimeout
	} else if udpTimeout, ok := am["udp-timeout"]; ok {
		params["UDPTransalationTimeout"] = udpTimeout
	} else {
		return nil, errors.New("TCP/UCP Translation timeout is not provided")
	}

	//params["Enable"] = "true"
	//params["Status"] = "Enable"

	parent := "Device.NAT."
	path := parent + "InterfaceSetting."
	info := &addInstInfo{
		path:   path,
		parent: parent,
		params: params,
	}
	return info, nil
}

const addNatPortMappingHelp = "add nat port-mapping alias <string> intf <id|name> all-intf <true|false> lease-duration <seconds> remote <address> extport <port> ext-port-range <range> int-port-range <range> <proto> <tcp|udp>"

func (cli *Cli) addNatPortMapping(c *ishell.Context) {
	if err := cli.checkDefault(); err != nil {
		c.Println(err)
		return
	}
	dest := destMtp
	instInfo, err := parseAddNatPortMappingArgs(c.Args, dest)
	if err != nil {
		c.Println(err)
		return
	}
	if instId, err := cli.addInst(instInfo, destMtp); err != nil {
		c.Println(err)
		return
	} else {
		c.Println("Instance created with the path:", instId)
		// Enable
		//if err := cli.MtpSetParamReq(instId, "Enable", "true"); err != nil {
		params := map[string]string{"Enable": "true"}
		if err := cli.restSetParams(instId, params); err != nil {
			c.Println("Error in activating Port Mapping:", err)
			return
		} else {
			c.Println("Enabled NAT PortMapping")
		}
	}

}

const addCfgNatPortMappingHelp = "addcfg nat port-mapping alias <string> intf <id|name> all-intf <true|false> lease-duration <seconds> remote <address> ext-port <port> ext-port-range <range> int-port-range <range> <proto> <tcp|udp>"

func (cli *Cli) addCfgNatPortMapping(c *ishell.Context) {
	if err := cli.checkCfgDevSet(); err != nil {
		c.Println(err)
		//c.Println("Use ", SetCfgDevInfoHelp)
		return
	}

	dest := destDb
	instInfo, err := parseAddNatPortMappingArgs(c.Args, dest)
	if err != nil {
		c.Println(err)
		return
	}
	if _, err := cli.addInst(instInfo, dest); err != nil {
		c.Println(err)
		return
	}
	c.Println("NAT PortMapping added to cfg datastore")

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

func parseAddNatPortMappingArgs(args []string, dest destType) (*addInstInfo, error) {
	// Receive arguments
	argLen := len(args)
	if argLen < 16 {
		return nil, errors.New("Wrong input")
	}
	am, _ := getMapFromArgs(args) // argMap

	params := make(map[string]string)

	if alias, ok := am["alias"]; ok {
		params["Alias"] = alias
	} else if dest == destDb {
		return nil, errors.New("alias is must for cfg")
	}
	//const addCfgNatPortMappingHelp = "addcfg nat port-mapping alias <string> intf <id|name> all-intf <true|false> lease-duration <seconds> remote <address> ext-port <port> ext-port-range <range> int-port-range <range> <proto> <tcp|udp>"

	intf := "Device.IP.Interface."
	if isDigit(am["intf"]) {
		intf = intf + am["intf"] + "."
	} else {
		intf = intf + "[Alias==\"" + am["intf"] + "\"]."
	}
	params["Interface"] = intf

	if am["all-intf"] == "true" {
		params["AllInterface"] = "true"
	} else if am["all-intf"] == "false" {
		params["AllInterface"] = "false"
	}
	params["LeaseDuration"] = am["lease-duration"]
	params["RemoteHost"] = am["remote"]
	params["ExternalPort"] = am["ext-port"]
	params["ExternalPortEndRange"] = am["ext-port-range"]
	params["InternalPort"] = am["int-port-range"]
	params["Protocol"] = am["proto"]

	// Validate Inputs

	// Create IPv4 address object
	var path string
	parent := "Device.NAT.PortMapping."
	info := &addInstInfo{
		path:   path,
		parent: parent,
		params: params,
	}
	return info, nil
}

const removeCfgNatHelp = "removecfg nat <intf-setting|port-mapping> <id|name>"

func (cli Cli) removeCfgNat(c *ishell.Context) {
	cli.removeCfgInst(c, "nat")
}

const removeNatHelp = "remove nat <intf-settins|port-mapping> <id|name>"

func (cli *Cli) removeNat(c *ishell.Context) {
	cli.removeInst(c, "nat")
}

const setCfgNatHelp = "setcfg nat <intf-setting|port-mapping> <id|name> <param> <value>"

func (cli *Cli) setCfgNat(c *ishell.Context) {
	if err := cli.setCfgParam(c.Args, "nat"); err != nil {
		c.Println(err)
		c.Println(setCfgNatHelp)
	}
}

const setNatHelp = "set nat <intf-setting|port-mapping> <id|name> <param> <value>"

func (cli *Cli) setNat(c *ishell.Context) {
	if err := cli.setParam(c.Args, "nat"); err != nil {
		c.Println(err)
		c.Println(setNatHelp)
	}
}
