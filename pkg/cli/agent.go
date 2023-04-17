package cli

import (
	"errors"
	"log"
	"os"

	"github.com/abiosoft/ishell"
)

type agentInfo struct {
	dev   devInfo
	epId  string
	isSet isSetType
}

func (cli *Cli) registerNounsAgent() {
	agentCmds := []noun{
		{"set", "agent", setAgentHelp, cli.setAgent},
		{"add", "agent", addAgentHelp, cli.addAgent},
		{"add.agent", "mtp", addAgentMtpHelp, cli.addAgentMtp},
		/*
			{"add.agent", "threshold", addAgentThresholdHelp, cli.addAgentThreshold},
			{"add.agent", "controller", addAgentControllerHelp, cli.addAgentController},
			{"add.agent.controller", "mtp", addAgentControllerMtpHelp, cli.addAgentControllerMtp},
			{"add.agent.controller.mtp", "coap", addAgentControllerMtpCoapHelp, cli.addAgentControllerMtpCoap},
			{"add.agent.controller.mtp", "stomp", addAgentControllerMtpStompHelp, cli.addAgentControllerMtpStomp},
			{"add.agent.controller.mtp", "websocket", addAgentControllerMtpWebSocketHelp, cli.addAgentControllerMtpWebSocket},
			{"add.agent.controller.mtp", "mqtt", addAgentControllerMtpMqttHelp, cli.addAgentControllerMtpMqtt},
			{"add.agent.controller.boot-param", "boot-param", addAgentControllerBootParamHelp, cli.addAgentControllerBootParam},
			{"add.agent.controller", "e2e", addAgentControllerE2eHelp, cli.addAgentControllerE2e},
			{"add.agent", "subscription", addAgentSubscriptionHelp, cli.addAgentSubscription},
			{"add.agent", "request", addAgentRequestHelp, cli.addAgentRequest},
			{"add.agent", "certificate", addAgentCertHelp, cli.addAgentCert},
			{"add.agent", "controller-trust", addAgentControllerTrustHelp, cli.addAgentControllerTrust},
			{"add.agent.controller-trust", "role", addAgentControllerTrustRoleHelp, cli.addAgentControllerTrustRole},
			{"add.agent.controller-trust.role", "permission", addAgentControllerTrustRolePermsHelp, cli.addAgentControllerTrustRolePerms},
			{"add.agent.controller-trust", "cred", addAgentControllerTrustCredHelp, cli.addAgentControllerTrustCred},
			{"add.agent.controller-trust", "challenge", addAgentControllerTrustChallengeHelp, cli.addAgentControllerTrustChallenge},
		*/

		{"unset", "agent", unsetAgentHelp, cli.unsetAgent},
		{"show", "agent", showAgentHelp, cli.showAgent},
	}
	cli.registerNouns(agentCmds)
}

const addAgentHelp = "add agent mtp|threshold|controller..."

func (cli *Cli) addAgent(c *ishell.Context) {
	c.Printf(addAgentHelp)
}

const addAgentMtpHelp = "add agent mtp alias <string> enable <true|false> protocol <stomp|coap|websocket> enable-mdns <true|false>"

func (cli *Cli) addAgentMtp(c *ishell.Context) {
	cli.lastCmdErr = errors.New("addAgentMtp Error")
	if err := cli.checkDefault(); err != nil {
		c.Println(err)
		cli.lastCmdErr = err
		return
	}
	dest := destMtp
	instInfo, err := parseAddAgentMtpArgs(c.Args, dest)
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

func parseAddAgentMtpArgs(args []string, dest destType) (*addInstInfo, error) {
	// Receive arguments
	argLen := len(args)
	var requiredArgs int

	if dest == destMtp {
		requiredArgs = 3
	} else {
		requiredArgs = 4
	}
	if argLen < requiredArgs {
		return nil, errors.New("Wrong input")
	}
	am, _ := getMapFromArgs(args) // argMap

	// Validate Inputs and form param map

	params := make(map[string]string)

	if enable, ok := am["enable"]; ok {
		params["Enable"] = enable
	} else {
		return nil, errors.New("Enable must be provided")
	}

	switch am["protocol"] {
	case "stomp":
		params["Protocol"] = "STOMP"
	case "coap":
		params["Protocol"] = "CoAP"
	case "websocket":
		params["Protocol"] = "WebSocket"
	case "mqtt":
		params["Protocol"] = "MQTT"
	default:
		return nil, errors.New("Unsupported protocol. Valid values: stomp|coap|websocket|mqtt")
	}

	if mdns, ok := am["enable-mdns"]; ok {
		params["EnableMDNS"] = mdns
	}

	if alias, ok := am["alias"]; ok {
		params["Alias"] = alias
	}

	parent := "Device.LocalAgent."
	path := parent + "MTP."
	info := &addInstInfo{
		path:   path,
		parent: parent,
		params: params,
	}
	return info, nil
}

const showAgentHelp = "show agent <cli-id|all-ids|mtp|threshold|controller|subsription|request|cert|controller-trust> [id] <coap|stomp|websocket|mqtt|mtp|boot-params|e2e|roles|creds|challenges> [id]"

func (cli *Cli) showAgent(c *ishell.Context) {
	obj := "all"
	if len(c.Args) > 0 {
		obj = c.Args[0]
	}
	if obj == "cli-id" {
		if !cli.agent.isSet.epId {
			c.Println("CLI Agent EndpointId is not configured, use set agentid to configure")
			cli.lastCmdErr = errors.New("Cli Agent Endpointid is not configured")
			return
		}
		c.Printf("%-25s : %-12s\n", "Cli Agent EndpointId", cli.agent.epId)
	} else if obj == "all-ids" {
		epIds, err := cli.restReadAgents()
		if err != nil {
			c.Println("No Valid Agent data found")
			cli.lastCmdErr = err
			return
		}
		for i, epId := range epIds {
			c.Printf("%-25s : %-12v\n", "Agent Number", i+1)
			c.Printf("%-25s : %-12s\n", "EndpointId", epId)
			c.Println("-------------------------------------------------")
		}

	} else if err := cli.showParams(c, "agent"); err != nil {
		c.Println(err)
		cli.lastCmdErr = err
		return
	}
	cli.lastCmdErr = nil
}

const setAgentHelp = "set agent cli-id <endpoint>"

func (cli *Cli) setAgent(c *ishell.Context) {
	if len(c.Args) < 2 {
		c.Println("Wrong input.", setAgentHelp)
		cli.lastCmdErr = errors.New("Wrong input, please provide id")
		return
	}
	if c.Args[0] == "cli-id" {
		epId := c.Args[1]
		cli.agent.epId = epId
		cli.agent.isSet.epId = true
		c.Println("Agent EndpointId set to:", epId)
		if err := cli.initCliWithAgentFactoryData(); err != nil {
			return
		}
	} else {
		c.Println("Wrong input", setAgentHelp)
		cli.lastCmdErr = errors.New("Wrong input, please provide id")
		return
	}
	cli.lastCmdErr = nil
}

const unsetAgentHelp = "unset agent id"

func (cli *Cli) unsetAgent(c *ishell.Context) {
	cli.agent.isSet.epId = false
	cli.lastCmdErr = nil
}

func (cli *Cli) initCliWithAgentParams() error {
	cli.agent.isSet.epId = false
	if epId, ok := os.LookupEnv("AGENT_ID"); !ok {
		log.Println("Please provide agent id either through env (AGENT_ID) or use set agentid command to set")
		return errors.New("Agent ID not found in env")
	} else {
		log.Println("Setting agent id to:", epId)
		cli.agent.epId = epId
		cli.agent.isSet.epId = true
	}
	if err := cli.initCliWithAgentFactoryData(); err != nil {
		return err
	}
	return nil
}
func (cli *Cli) initCliWithAgentFactoryData() error {
	cli.agent.isSet.productClass = false
	cli.agent.isSet.manufacturer = false
	cli.agent.isSet.modelName = false
	devInfo, err := cli.restReadParams("Device.DeviceInfo.")
	if err != nil {
		log.Println("RestErr:", err)
		return err
	}
	for _, param := range devInfo[0].Params {
		if param.Name == "ProductClass" {
			cli.agent.dev.productClass = param.Value
			cli.agent.isSet.productClass = true
			log.Println("Agent ProductClass", param.Value)
		}
		if param.Name == "ManufacturerOUI" {
			cli.agent.dev.manufacturer = param.Value
			cli.agent.isSet.manufacturer = true
			log.Println("Agent ManufacturerOUI", param.Value)
		}
		if param.Name == "SerialNumber" {
			cli.agent.dev.modelName = param.Value
			cli.agent.isSet.modelName = true
			log.Println("Agent SerialNumber", param.Value)
		}
	}
	return nil
}
