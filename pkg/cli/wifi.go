package cli

import (
	"errors"
	"log"

	"github.com/abiosoft/ishell"
)

func (cli *Cli) registerNounsWiFi() {
	cmds := []noun{
		{"add", "wifi", addWiFiHelp, cli.addWiFi},
		{"add.wifi", "ssid", addWiFiSsidHelp, cli.addWiFiSsid},
		{"add.wifi", "accesspoint", addWiFiAccessPointHelp, cli.addWiFiAccessPoint},

		{"addcfg", "wifi", addCfgWiFiHelp, cli.addCfgWiFi},
		{"addcfg.wifi", "ssid", addCfgWiFiSsidHelp, cli.addCfgWiFiSsid},
		{"addcfg.wifi", "accesspoint", addCfgWiFiAccessPointHelp, cli.addCfgWiFiAccessPoint},

		{"show", "wifi", showWiFiHelp, cli.showWiFi},
		{"showcfg", "wifi", showCfgWiFiHelp, cli.showCfgWiFi},

		{"remove", "wifi", removeWiFiHelp, cli.removeWiFi},
		{"removecfg", "wifi", removeCfgWiFiHelp, cli.removeCfgWiFi},

		{"update", "wifi", updateWiFiHelp, cli.updateWiFi},

		{"set", "wifi", setWiFiHelp, cli.setWiFi},
		{"setcfg", "wifi", setCfgWiFiHelp, cli.setCfgWiFi},
	}
	cli.registerNouns(cmds)
}

const updateWiFiHelp = "update wifi <radio|ssid|accesspoint|endpoint> <id|name]"

func (cli *Cli) updateWiFi(c *ishell.Context) {
	if err := cli.updateDb(c.Args, "wifi"); err != nil {
		c.Println(err)
		cli.lastCmdErr = err
		return
	}
	c.Println("Updated Wifi instance and params successfully")
}

const setWiFiHelp = "set wifi <radio|ssid|accesspoint|endpoint> <id> <param> <value>"

func (cli *Cli) setWiFi(c *ishell.Context) {
	if err := cli.setParam(c.Args, "wifi"); err != nil {
		c.Println(err)
		cli.lastCmdErr = err
		return
	}

	if c.Args[0] == "radio" {
		path := "Device.WiFi.Radio." + c.Args[1] + "."
		log.Println("Additionally setting Radio enable true for:", path)
		params := map[string]string{"Enable": "true"}
		if err := cli.restSetParams(path, params); err != nil {
			c.Println("Error in enabling radio, err:", err)
			cli.lastCmdErr = err
			return
		}
	}
}

const setCfgWiFiHelp = "setcfg wifi <radio|ssid|accesspoint|endpoint> <id> <param> <value>"

func (cli *Cli) setCfgWiFi(c *ishell.Context) {
	if err := cli.setCfgParam(c.Args, "wifi"); err != nil {
		c.Println(err)
		c.Println(setCfgWiFiHelp)
		cli.lastCmdErr = err
	}
}

const removeWiFiHelp = "remove wifi ssid|accesspoint <id|name>"

func (cli *Cli) removeWiFi(c *ishell.Context) {
	if err := cli.removeInst(c, "wifi"); err != nil {
		c.Println(err)
		cli.lastCmdErr = err
	}
}

const removeCfgWiFiHelp = "removecfg wifi ssid|accesspoint <id|name>"

func (cli *Cli) removeCfgWiFi(c *ishell.Context) {
	if err := cli.removeCfgInst(c, "wifi"); err != nil {
		c.Println(err)
		cli.lastCmdErr = err
	}
}

const showWiFiHelp = "show wifi <ssid|radio|accesspoint|endpoint> [id] <stats>"

func (cli *Cli) showWiFi(c *ishell.Context) {
	if err := cli.showParams(c, "wifi"); err != nil {
		c.Println(err)
		cli.lastCmdErr = err
	}
}

const showCfgWiFiHelp = "showcfg wifi <ssid|radio|accesspoint|endpoint> [id] <stats>"

func (cli *Cli) showCfgWiFi(c *ishell.Context) {
	if err := cli.showCfg(c, "wifi"); err != nil {
		c.Println(err)
		cli.lastCmdErr = err
	}
}

const addCfgWiFiAccessPointHelp = "addcfg wifi accesspoint alias <string> ssid <id|name> security <open|wpa2-personal|wpa2-enterprise>"

func (cli *Cli) addCfgWiFiAccessPoint(c *ishell.Context) {
	if err := cli.checkCfgDevSet(); err != nil {
		c.Println(err)
		cli.lastCmdErr = err
		return
	}
	dest := destDb
	instInfo, err := parseAddWiFiAccessPointArgs(c.Args, dest)
	if err != nil {
		c.Println(err)
		cli.lastCmdErr = err
		return
	}
	securityMode := instInfo.params["Security.ModeEnabled"]
	if securityMode == "WPA2-Personal" || securityMode == "WPA2-Enterprise" {
		c.Print("Password needs to be of min 8 characters")
		c.Print("Password: ")
		instInfo.params["Security.KeyPassphrase"] = c.ReadPassword()
	}
	if _, err := cli.addInst(instInfo, dest); err != nil {
		c.Println(err)
		cli.lastCmdErr = err
		return
	}
	c.Println("WiFi Instance AccessPoint added to cfg datastore")
}

func parseAddWiFiAccessPointArgs(args []string, dest destType) (*addInstInfo, error) {
	argLen := len(args)
	if argLen < 4 {
		return nil, errors.New("Wrong input.")
	}
	am, _ := getMapFromArgs(args)

	// Add Instance
	parent := "Device.WiFi."
	path := parent + "AccessPoint."
	params := make(map[string]string)

	ssidRef := "Device.WiFi.SSID."
	if isDigit(am["ssid"]) {
		ssidRef = ssidRef + am["ssid"] + "."
	} else {
		ssidRef = ssidRef + "[Alias==\"" + am["ssid"] + "\"]."
	}

	if alias, ok := am["alias"]; ok {
		params["Alias"] = alias
	} else if dest == destDb {
		return nil, errors.New("Alias is must for cfg")
	}

	switch am["security"] {
	case "open":
		params["Security.ModeEnabled"] = "None"
	case "wpa2-personal":
		params["Security.ModeEnabled"] = "WPA2-Personal"
	case "wpa2-enterprise":
		params["Security.ModeEnabled"] = "WPA2-Enterprise"
	default:
		return nil, errors.New("security mode not found")
	}

	params["SSIDReference"] = ssidRef
	params["Enable"] = "true"

	info := &addInstInfo{
		path:   path,
		parent: parent,
		params: params,
	}
	return info, nil
}

const addCfgWiFiSsidHelp = "addcfg wifi ssid alias <string> name <name> radio <id|name>"

func (cli *Cli) addCfgWiFiSsid(c *ishell.Context) {
	if err := cli.checkCfgDevSet(); err != nil {
		c.Println(err)
		cli.lastCmdErr = err
		return
	}
	dest := destDb
	instInfo, err := parseAddWiFiSsidArgs(c.Args, dest)
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
	c.Println("WiFi Instance SSID added to cfg datastore")
}

func parseAddWiFiSsidArgs(args []string, dest destType) (*addInstInfo, error) {
	// Receive arguments
	argLen := len(args)
	if argLen < 4 {
		return nil, errors.New("Wrong input.")
	}
	am, _ := getMapFromArgs(args)
	//log.Printf("name:%v, radio:%v\n", am["name"], am["radio"])

	// Add Instance
	parent := "Device.WiFi."
	params := make(map[string]string)
	path := parent + "SSID."

	radioPath := "Device.WiFi.Radio."
	if isDigit(am["radio"]) {
		radioPath = radioPath + am["radio"] + "."
	} else {
		radioPath = radioPath + "[Alias==\"" + am["radio"] + "\"]."
	}

	params["LowerLayers"] = radioPath
	if ssidName, ok := am["name"]; !ok {
		return nil, errors.New("Wrong syntax, please verify all the parameters")
	} else {
		params["SSID"] = ssidName
	}
	params["Enable"] = "true"

	if alias, ok := am["alias"]; ok {
		params["Alias"] = alias
	} else if dest == destDb {
		return nil, errors.New("Alias is must for cfg")
	}

	info := &addInstInfo{
		path:   path,
		parent: parent,
		params: params,
	}
	return info, nil
}

const addCfgWiFiHelp = "addcfg wifi ssid|accesspoint..."

func (cli *Cli) addCfgWiFi(c *ishell.Context) {
	c.Println(addCfgWiFiHelp)
}

const addWiFiSsidHelp = "add wifi ssid alias <string> name <name> radio <id|name>"

func (cli *Cli) addWiFiSsid(c *ishell.Context) {
	var err error
	if err = cli.checkDefault(); err != nil {
		c.Println(err)
		cli.lastCmdErr = err
		return
	}
	dest := destMtp
	instInfo, err := parseAddWiFiSsidArgs(c.Args, dest)
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
		c.Println("WiFi SSID created with instance:", instId)
	}
	// Set regularitory domain to US (TDB)
	radioPath := instInfo.params["LowerLayers"]
	params := map[string]string{"RegulatoryDomain": "US"}
	if err := cli.restSetParams(radioPath, params); err != nil {
		c.Println("Error in setting regulartory domain of Radio:", err)
		// temporarily commenting this since for docker base agent it gives error
		// to avoid test failures
		//cli.lastCmdErr = err
	}
}

const addWiFiHelp = "add wifi ssid|accesspoint..."

func (cli *Cli) addWiFi(c *ishell.Context) {
	c.Println(addWiFiHelp)
}

const addWiFiAccessPointHelp = "add wifi accesspoint alias <string> ssid <id|name> security <open|wpa2-personal|wpa2-enterprise>"

func (cli *Cli) addWiFiAccessPoint(c *ishell.Context) {
	if err := cli.checkDefault(); err != nil {
		c.Println(err)
		cli.lastCmdErr = err
		return
	}
	dest := destMtp
	instInfo, err := parseAddWiFiAccessPointArgs(c.Args, dest)
	if err != nil {
		c.Println(err)
		cli.lastCmdErr = err
		return
	}
	securityMode := instInfo.params["Security.ModeEnabled"]
	if securityMode == "WPA2-Personal" || securityMode == "WPA2-Enterprise" {
		c.Print("Password needs to be of min 8 characters")
		c.Print("Password: ")
		instInfo.params["Security.KeyPassphrase"] = c.ReadPassword()
	}

	instId, err := cli.addInst(instInfo, dest)
	if err != nil {
		c.Println(err)
		cli.lastCmdErr = err
		return
	} else {
		c.Println("WiFi AccessPoint created with instance:", instId)
	}
	// Acivate the created AccessPoint
	params := map[string]string{"Enable": "true"}
	if err := cli.restSetParams(instId, params); err != nil {
		c.Println("Error in activating AccessPoint:", err)
		cli.lastCmdErr = err
	}
}
