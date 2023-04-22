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
	"log"

	"github.com/abiosoft/ishell"
)

func (cli *Cli) registerNounsNw() {
	cmds := []noun{
		/*
			{"add", "network", addIpHelp, cli.addIp},
			{"add.ip", "intf", addIpIntfHelp, cli.addIpIntf},
			{"add.ip", "addr", addIpAddrHelp, cli.addIpAddr},

			{"addcfg", "ip", addCfgIpHelp, cli.addCfgIp},
			{"addcfg.ip", "intf", addCfgIpIntfHelp, cli.addCfgIpIntf},
			{"addcfg.ip", "addr", addCfgIpAddrHelp, cli.addCfgIpAddr},
		*/

		{"show", "nw", showNwHelp, cli.showNw},
		{"show.nw", "wireless", showNwWirelessHelp, cli.showNwWireless},
		{"show.nw", "wired", showNwWiredHelp, cli.showNwWired},

		/*
			{"showcfg", "ip", showCfgIpHelp, cli.showCfgIp},

			{"remove", "ip", removeIpHelp, cli.removeIp},
			{"removecfg", "ip", removeCfgIpHelp, cli.removeCfgIp},

			{"set", "ip", setIpHelp, cli.setIp},
			{"setcfg", "ip", setCfgIpHelp, cli.setCfgIp},
		*/
	}
	cli.registerNouns(cmds)
}

const showNwHelp = "show nw wired|wireless"

func (cli *Cli) showNw(c *ishell.Context) {
	c.Println(showNwHelp)
}

const showNwWirelessHelp = "show nw wireless..."

func (cli *Cli) showNwWireless(c *ishell.Context) {
	// Get AccessPoint instances
	apObjParams, err := cli.restReadParams("Device.WiFi.AccessPoint.")
	if err != nil {
		log.Println("restErr:", err)
		cli.lastCmdErr = err
		return
	}
	for _, apObj := range apObjParams {
		log.Println("AP path:", apObj.Path)
		for _, apParam := range apObj.Params {
			if apParam.Name == "Alias" {
				c.Printf("%-25s : %-12s\n", "AP Name", apParam.Value)
			}
			if apParam.Name == "Enable" {
				c.Printf("%-25s : %-12s\n", "Enable", apParam.Value)
			}
			if apParam.Name == "Status" {
				c.Printf("%-25s : %-12s\n", "Status", apParam.Value)
			}
			if apParam.Name == "SSIDReference" {
				ssidObjs, err := cli.restReadParams(apParam.Value)
				if err != nil {
					log.Println("restErr:", err)
					continue
				}
				log.Println("SSID path:", ssidObjs[0].Path)
				for _, ssidParam := range ssidObjs[0].Params {
					if ssidParam.Name == "Name" {
						c.Printf("%-25s : %-12s\n", "SSID Name", ssidParam.Value)
					}
					if ssidParam.Name == "Status" {
						c.Printf("%-25s : %-12s\n", "SSID Status", ssidParam.Value)
					}
					if ssidParam.Name == "LowerLayers" {
						log.Println("ssidParam.LowerLayers", ssidParam.Value)
						radioObjs, err := cli.restReadParams(ssidParam.Value)
						if err != nil {
							log.Println("restErr:", err)
							continue
						}
						log.Println("Radio path:", radioObjs[0].Path)
						for _, radioParam := range radioObjs[0].Params {
							if radioParam.Name == "Channel" {
								c.Printf("%-25s : %-12s\n", "Radio Channel", radioParam.Value)
							}
							if radioParam.Name == "Enable" {
								c.Printf("%-25s : %-12s\n", "Radio Enable", radioParam.Value)
							}
							if radioParam.Name == "TransmitPower" {
								c.Printf("%-25s : %-12s\n", "Radio TxPower", radioParam.Value)
							}
						}
					}
				}
			}
		}
	}
	c.Println("-------------------------------------------------")
	cli.lastCmdErr = nil
}

const showNwWiredHelp = "show nw wired..."

func (cli *Cli) showNwWired(c *ishell.Context) {
	// Get IP Interface Paramlist
	intfObjs, err := cli.restReadParams("Device.IP.Interface.")
	if err != nil {
		log.Println("Error:", err)
		cli.lastCmdErr = err
		return
	}

	for _, intfObj := range intfObjs {
		log.Println("IntfPath:", intfObj.Path)
		for _, intfParam := range intfObj.Params {
			if intfParam.Name == "Alias" {
				c.Printf("%-25s : %-12s\n", "Interface Name", intfParam.Value)
			}
			if intfParam.Name == "Enable" {
				c.Printf("%-25s : %-12s\n", "Enable", intfParam.Value)
			}
			if intfParam.Name == "Status" {
				c.Printf("%-25s : %-12s\n", "Status", intfParam.Value)
			}
			if intfParam.Name == "IPv4Enable" {
				c.Printf("%-25s : %-12s\n", "IPv4 Status", intfParam.Value)
			}
			if intfParam.Name == "IPv6Enable" {
				c.Printf("%-25s : %-12s\n", "IPv6 Status", intfParam.Value)
			}
			if intfParam.Name == "LowerLayers" {
				c.Printf("%-25s : %-12s\n", "Phy", intfParam.Value)
			}
			if intfParam.Name == "Type" {
				c.Printf("%-25s : %-12s\n", "Type", intfParam.Value)
			}
		}

		/*
			ipv4addrPath := intfPath + "IPv4Address.\\d*.$"
			//log.Println("ipv4addrpath:", ipv4addrPath)
			addrs, err2 := cli.dbGetInstancesByRegex(ipv4addrPath)
			if err2 != nil {
				log.Println("Error:", err2)
				cli.lastCmdErr = err2
			}

			var ipv4path string
			for _, addr := range addrs {
				ipv4path = addr.path
				//log.Println("ipv4Path:", ipv4path)
				c.Printf(" %-24s : %-12s\n", "IPv4 Address Enable", ip[ipv4path+"Enable"])
				c.Printf(" %-24s : %-12s\n", "IPv4 Address Status", ip[ipv4path+"Status"])
				c.Printf(" %-24s : %-12s\n", "IPv4 Address", ip[ipv4path+"IPAddress"])
				c.Printf(" %-24s : %-12s\n", "Subnet", ip[ipv4path+"SubnetMask"])
				c.Printf(" %-24s : %-12s\n", "Type", ip[ipv4path+"AddressingType"])
			}
		*/
		c.Println("-------------------------------------------------")
	}
	cli.lastCmdErr = nil
}
