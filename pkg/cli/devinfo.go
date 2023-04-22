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
	"log"

	"github.com/abiosoft/ishell"
)

type devInfo struct {
	productClass string
	manufacturer string
	modelName    string
}
type isSetType struct {
	productClass bool
	manufacturer bool
	modelName    bool
	epId         bool
}

func (cli *Cli) registerNounsDevInfo() {
	devInfoCmds := []noun{
		{"show", "devinfo", showDevInfoHelp, cli.showDevInfo},
		{"set", "devinfo", setDevInfoHelp, cli.setDevInfo},
		{"remove", "devinfo", removeDevInfoHelp, cli.removeDevInfo},
		{"add", "devinfo", addDevInfoHelp, cli.addDevInfo},
		{"add.devinfo", "firmware", addDevInfoFirmwareHelp, cli.addDevInfoFirmware},

		{"showcfg", "devinfo", showCfgDevInfoHelp, cli.showCfgDevInfo},
		{"operate", "devinfo", operateDevInfoHelp, cli.operateDevInfo},
		{"operate.devinfo", "firmware", operateDevInfoFirmwareHelp, cli.operateDevInfoFirmware},
		{"operate.devinfo.firmware", "download", operateDevInfoFirmwareDownloadHelp, cli.operateDevInfoFirmwareDownload},
		{"operate.devinfo.firmware", "activate", operateDevInfoFirmwareActivateHelp, cli.operateDevInfoFirmwareActivate},
	}
	cli.registerNouns(devInfoCmds)
}

const operateDevInfoHelp = "operate devinfo firmware..."

func (cli *Cli) operateDevInfo(c *ishell.Context) {
	c.Printf(operateDevInfoHelp)
}

const operateDevInfoFirmwareHelp = "operate devinfo firmware download|activate..."

func (cli *Cli) operateDevInfoFirmware(c *ishell.Context) {
	c.Printf(operateDevInfoFirmwareHelp)
}

const operateDevInfoFirmwareDownloadHelp = "operate devinfo firmware download img-id <id|name> url <url-path> username <string> password <string> auto-activate <true|false> checksum-algo <SHA-1|SHA-224|SHA-256> checksum <value>"

func (cli *Cli) operateDevInfoFirmwareDownload(c *ishell.Context) {
	if err := cli.checkDefault(); err != nil {
		cli.lastCmdErr = err
		return
	}
	am, err := getMapFromArgs(c.Args)
	if err != nil {
		c.Println("Error in parsing input parameters")
		cli.lastCmdErr = err
		return
	}
	path := "Device.DeviceInfo.FirmwareImage."

	if id, ok := am["img-id"]; ok {
		if isDigit(id) {
			path = path + id + ".Download()"
		} else {
			path = path + "[Alias==\"" + id + "\"].Download()."
		}
	}
	inputs := make(map[string]string)

	if url, ok := am["url"]; ok {
		inputs["URL"] = url
	} else {
		c.Println("Must provide URL of image location")
		cli.lastCmdErr = errors.New("no URL path found")
		return
	}

	if userName, ok := am["username"]; ok {
		inputs["Username"] = userName
	}

	if pass, ok := am["password"]; ok {
		inputs["Password"] = pass
	}

	if autoActivate, ok := am["auto-activate"]; ok {
		inputs["AutoActivate"] = autoActivate
	}

	if checksumAlgo, ok := am["checksum-algo"]; ok {
		inputs["CheckSumAlgorithm"] = checksumAlgo
	}

	if checksum, ok := am["checksum"]; ok {
		inputs["CheckSum"] = checksum
	}

	log.Println("Firmware Download Path:", path)
	log.Printf("Params: %+v\n", inputs)

	// TODO: retrieve command from arg and pass it here
	if err := cli.restOperateCmd(path, inputs); err != nil {
		log.Println("Error in  executing rest Operate method", err)
		return
	}
}

const operateDevInfoFirmwareActivateHelp = "operate devinfo firmware activate img-id <id|name>"

func (cli *Cli) operateDevInfoFirmwareActivate(c *ishell.Context) {
	if err := cli.checkDefault(); err != nil {
		cli.lastCmdErr = err
		return
	}
	if c.Args[0] != "img-id" {
		c.Println("Please provide img-id")
		cli.lastCmdErr = errors.New("Image id not provided")
		return
	}

	path := "Device.DeviceInfo.FirmwareImage."
	imgId := c.Args[1]

	if isDigit(imgId) {
		path = path + imgId + ".Activate()"
	} else {
		path = path + "[Alias==\"" + imgId + "\"].Activate()."
	}
	log.Println("Firmware Activate Path:", path)

	// TODO: retrieve command from arg and pass it here
	if err := cli.restOperateCmd(path, nil); err != nil {
		log.Println("Error in  executing RESt Operate cmd", err)
		return
	}
}

const addDevInfoHelp = "add devinfo firmware..."

func (cli *Cli) addDevInfo(c *ishell.Context) {
	c.Printf(addDevInfoHelp)
}

const addDevInfoFirmwareHelp = "add devinfo firmware alias <string> available <true|false>"

func (cli *Cli) addDevInfoFirmware(c *ishell.Context) {
	if err := cli.checkDefault(); err != nil {
		c.Println(err)
		cli.lastCmdErr = err
		return
	}
	dest := destMtp
	instInfo, err := parseAddDevInfoFirmwareArgs(c.Args, dest)
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

func parseAddDevInfoFirmwareArgs(args []string, dest destType) (*addInstInfo, error) {
	// Receive arguments
	argLen := len(args)
	var requiredArgs int

	if dest == destMtp {
		requiredArgs = 2
	} else {
		requiredArgs = 4
	}
	if argLen < requiredArgs {
		return nil, errors.New("Wrong input")
	}
	am, _ := getMapFromArgs(args) // argMap

	// Validate Inputs and form param map

	params := make(map[string]string)

	if am["available"] == "true" {
		params["Available"] = "true"
	} else {
		params["Available"] = "true"
	}

	if alias, ok := am["alias"]; ok {
		params["Alias"] = alias
	} else if dest == destDb {
		return nil, errors.New("Alias is must for cfg")
	}

	parent := "Device.DeviceInfo."
	path := parent + "FirmwareImage."
	info := &addInstInfo{
		path:   path,
		parent: parent,
		params: params,
	}
	return info, nil
}

const removeDevInfoHelp = "remove devinfo firmeware <id|name>"

func (cli *Cli) removeDevInfo(c *ishell.Context) {
	if err := cli.removeInst(c, "devinfo"); err != nil {
		c.Println(err)
		cli.lastCmdErr = err
		return
	}
	cli.lastCmdErr = nil
	return
}

const setDevInfoHelp = "set devinfo <firmware|image|loc|logfile|vendorcfg> id <number|alias> availble <true|false> "

func (cli *Cli) setDevInfo(c *ishell.Context) {
	if err := cli.setParam(c.Args, "devinfo"); err != nil {
		c.Println("Error:", err)
		cli.lastCmdErr = err
		return
	}
	cli.lastCmdErr = nil
}

const showDevInfoHelp = "show devinfo <vendorcfg|memory|proc|temp|net|cpu|logfile|loc|imagefile|firmware> <id>"

func (cli *Cli) showDevInfo(c *ishell.Context) {
	if err := cli.showParams(c, "devinfo"); err != nil {
		c.Println(err)
		cli.lastCmdErr = err
		return
	}
	cli.lastCmdErr = nil
}

const showCfgDevInfoHelp = "showcfg devinfo"

func (cli *Cli) showCfgDevInfo(c *ishell.Context) {
	if !cli.agent.isSet.productClass {
		c.Println("Agent Dev Info is not configured, use setcfg devinfo to configure")
		cli.lastCmdErr = errors.New("Agent dev info is not set")
		return
	}
	c.Printf("  %-23s : %-12s\n", "Product Class", cli.agent.dev.productClass)
	c.Printf("  %-23s : %-12s\n", "Manufacturer", cli.agent.dev.manufacturer)
	c.Printf("  %-23s : %-12s\n", "Model Name", cli.agent.dev.modelName)
	c.Println("-------------------------------------------------")
	cli.lastCmdErr = nil
}
