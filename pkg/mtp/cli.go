package mtp

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/abiosoft/ishell"
	"github.com/n4-networks/openusp/pkg/parser"
	"github.com/n4-networks/openusp/pkg/pb/bbf/usp_msg"
)

var (
	cliAgentId string = "os::SSGspa-02:42:ac:11:00:05"
	agentIdSet bool   = false
)

func (m *Mtp) Cli() {
	shell := ishell.New()
	shell.Println("N4 MTP Cli")
	shell.SetPrompt("MTP>>")

	var ok bool
	if cliAgentId, ok = os.LookupEnv("AGENT_ID"); ok {
		agentIdSet = true
	} else {
		log.Println("Please set agent id through env AGENT_ID")
		return
	}

	shell.AddCmd(&ishell.Cmd{
		Name: "getparam",
		Help: "get msg to USP Agent.  ex: get path1, path2 ...",
		Func: m.cliSendGetMsg,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "getinstance",
		Help: "getinstance firstlevelonly (true/false), objpath",
		Func: m.cliSendGetInstanceMsg,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "setparam",
		Help: "set parameter value. ex: set path param value",
		Func: m.cliSendSetMsg,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "addinstance",
		Help: "add instance of an object and its parameters. ex: add objpath param value",
		Func: m.cliSendAddMsg,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "operate",
		Help: "operate msg to an USP Agent. ex: operate Device.Reboot()",
		Func: m.cliSendOperateMsg,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "getdm",
		Help: "getdm msg to an USP Agent. ex: getdm Device.DeviceInfo.",
		Func: m.cliSendGetSupportedDmMsg,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "delinstance",
		Help: "delete an instance of an multi-instance object. ex: delinstance Device.DeviceInfo.",
		Func: m.cliSendDeleteInstanceMsg,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "getresp",
		Help: "receive respone from agent, ex: getresp",
		Func: func(c *ishell.Context) {
			m.StompReceiveUspMsgFromAgentWithTimer(1)
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "getver",
		Help: "Print MTP version",
		Func: func(c *ishell.Context) {
			c.Println("MTP Version:", getVer())
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "setlog",
		Help: SetLogHelp,
		Func: m.cliSetLog,
	})
	shell.Run()
}

func (m *Mtp) cliSendGetMsg(c *ishell.Context) {
	mtpIntf, err := m.getAgentMtp(cliAgentId)
	if err != nil {
		c.Println("Could not get agent MTP interface")
		return
	}
	msgId := "MTPCLI_GET_PARAM_" + strconv.FormatUint(mtpIntf.getMsgCnt(), 10)
	uspMsg, err := parser.CreateUspGetReqMsg(c.Args, msgId)
	if err != nil {
		c.Println("Could not prepare Usp msg of type: GET")
		return
	}
	if err := m.sendUspMsgToAgent(cliAgentId, uspMsg, mtpIntf); err != nil {
		c.Println("Error in sending get msg to agent")
	}
}

func (m *Mtp) cliSendGetInstanceMsg(c *ishell.Context) {
	var firstLevelOnly bool
	switch c.Args[0] {
	case "false":
		firstLevelOnly = false
	case "true":
		firstLevelOnly = true
	default:
		c.Println("getinstance firstlevelonly(true/false) objpath")
		return
	}
	mtpIntf, err := m.getAgentMtp(cliAgentId)
	if err != nil {
		c.Println("Could not get agent MTP interface")
		return
	}
	msgId := "MTPCLI_GET_INSTANCE_" + strconv.FormatUint(mtpIntf.getMsgCnt(), 10)
	uspMsg, err := parser.CreateUspGetInstancesMsg(c.Args[1:], firstLevelOnly, msgId)
	if err != nil {
		c.Println("Could not prepare Usp  msg of type: GET_INSTANCES")
		return
	}
	if err := m.sendUspMsgToAgent(cliAgentId, uspMsg, mtpIntf); err != nil {
		c.Println("Error in sending get msg to agent")
	}
}

func (m *Mtp) cliSendSetMsg(c *ishell.Context) {
	if len(c.Args) < 3 {
		c.Println("Error: Insufficient number of args (path/param/val)")
		return
	}
	mtpIntf, err := m.getAgentMtp(cliAgentId)
	if err != nil {
		c.Println("Could not get agent MTP interface")
		return
	}
	msgId := "MTPCLI_SET_PARAM__" + strconv.FormatUint(mtpIntf.getMsgCnt(), 10)
	uspMsg, err := createSetMsg(c.Args[0], c.Args[1], c.Args[2], msgId)
	if err != nil {
		c.Println("Could not prepare Usp  msg of type: SET")
		return
	}
	if err := m.sendUspMsgToAgent(cliAgentId, uspMsg, mtpIntf); err != nil {
		c.Println("Error in sending get msg to agent")
	}
}

func (m *Mtp) cliSendAddMsg(c *ishell.Context) {
	if len(c.Args) < 3 {
		c.Println("Error: Insufficient number of args: cmd objpath paramName value)")
		return
	}
	mtpIntf, err := m.getAgentMtp(cliAgentId)
	if err != nil {
		c.Println("Could not get agent MTP interface")
		return
	}
	msgId := "MTPCLI_ADD_INSTANCE__" + strconv.FormatUint(mtpIntf.getMsgCnt(), 10)
	uspMsg, err := createAddMsg(c.Args[0], c.Args[1], c.Args[2], msgId)
	if err != nil {
		c.Println("Could not prepare Usp  msg of type: ADD")
		return
	}
	if err := m.sendUspMsgToAgent(cliAgentId, uspMsg, mtpIntf); err != nil {
		c.Println("Error in sending get msg to agent")
	}
}

func (m *Mtp) cliSendOperateMsg(c *ishell.Context) {
	cmd := c.Args[0] //[string("Device.Reboot()")
	cmdKey := "none"
	resp := true

	mtpIntf, err := m.getAgentMtp(cliAgentId)
	if err != nil {
		c.Println("Could not get agent MTP interface")
		return
	}
	msgId := "MTPCLI_OPERATE_" + strconv.FormatUint(mtpIntf.getMsgCnt(), 10)
	uspMsg, err := parser.CreateUspOperateReqMsg(cmd, cmdKey, resp, msgId, nil)
	if err != nil {
		c.Println("Could not prepare Usp msg of type: OPERATE")
		return
	}
	if err := m.sendUspMsgToAgent(cliAgentId, uspMsg, mtpIntf); err != nil {
		c.Println("Error in sending get msg to agent")
	}
}

func (m *Mtp) cliSendGetProtoVersionMsg(c *ishell.Context) {
	mtpIntf, err := m.getAgentMtp(cliAgentId)
	if err != nil {
		c.Println("Could not get agent MTP interface")
		return
	}
	msgId := "MTPCLI_GET_PROTO_VER_" + strconv.FormatUint(mtpIntf.getMsgCnt(), 10)
	uspMsg, err := parser.CreateUspGetSupportedProtoMsg(m.Cfg.Usp.ProtoVersion, msgId)
	if err != nil {
		c.Println("Could not prepare Usp msg of type: GET_PROTO_VERSION")
		return
	}
	if err := m.sendUspMsgToAgent(cliAgentId, uspMsg, mtpIntf); err != nil {
		c.Println("Error in sending get msg to agent")
	}
}

func (m *Mtp) cliSendDeleteInstanceMsg(c *ishell.Context) {
	mtpIntf, err := m.getAgentMtp(cliAgentId)
	if err != nil {
		c.Println("Could not get agent MTP interface")
		return
	}
	msgId := "MTPCLI_DEL_INSTANCE_" + strconv.FormatUint(mtpIntf.getMsgCnt(), 10)
	uspMsg, err := parser.CreateUspDeleteReqMsg(c.Args, msgId)
	if err != nil {
		c.Println("Could not prepare Usp msg of type: DELETE")
		return
	}
	if err := m.sendUspMsgToAgent(cliAgentId, uspMsg, mtpIntf); err != nil {
		c.Println("Error in sending get msg to agent")
	}
}

func (m *Mtp) cliSendGetSupportedDmMsg(c *ishell.Context) {
	mtpIntf, err := m.getAgentMtp(cliAgentId)
	if err != nil {
		log.Println("Could not get agent mtpIntf")
		return
	}
	msgId := "MTPCLI_GET_DM_" + strconv.FormatUint(mtpIntf.getMsgCnt(), 10)
	uspMsg, err := parser.CreateUspGetSupportedDmMsg(c.Args, true, true, true, msgId)
	if err != nil {
		c.Println("Could not prepare Usp msg of type: GET_SUPPORTED_DM")
		return
	}
	if err := m.sendUspMsgToAgent(cliAgentId, uspMsg, mtpIntf); err != nil {
		c.Println("Error in sending get msg to agent")
	}
}

const SetLogHelp = "setlog <on|off> <all|date|time|long|short>"

func (m *Mtp) cliSetLog(c *ishell.Context) {
	if len(c.Args) < 2 {
		c.Println("setlog help")
		return
	}
	cmd := c.Args[0]
	flag := c.Args[1]
	switch flag {
	case "all":
		if cmd == "off" {
			log.SetFlags(0)
			log.SetOutput(ioutil.Discard)
			c.Println("Switched off all logging flags")
			log.Println("If you are seeing this then something is broken...")
		} else if cmd == "on" {
			log.SetFlags(log.Lshortfile | log.Llongfile | log.Ldate | log.Ltime)
			log.SetOutput(os.Stdout)
			c.Println("Switched on all logging flags")
			log.Println("This message from logging engine, its has been switched on now")
		} else {
			c.Println("Invalid cmd. Syntax:", SetLogHelp)
		}
	case "date":
		if cmd == "off" {
			f := log.Flags() &^ log.Ldate
			log.SetFlags(f)
			c.Println("Logging date flag has been switched off")
			log.Println("This msg through log should not have date")
		} else if cmd == "on" {
			f := log.Flags() | log.Ldate
			log.SetFlags(f)
			c.Println("Logging date flag has been switched on")
			log.Println("This msg through log should have date")
		} else {
			c.Println("Invalid cmd. Syntax:", SetLogHelp)
		}
	case "time":
		if cmd == "off" {
			f := log.Flags() &^ log.Ltime
			log.SetFlags(f)
			c.Println("Logging time flag has been switched off")
			log.Println("This msg through log should not have time")
		} else if cmd == "on" {
			f := log.Flags() | log.Ltime
			log.SetFlags(f)
			c.Println("Logging time flag has been switched on")
			log.Println("This msg through log should have time")
		} else {
			c.Println("Invalid cmd. Syntax:", SetLogHelp)
		}
	case "long":
		if cmd == "off" {
			f := log.Flags() &^ log.Llongfile
			log.SetFlags(f)
			c.Println("Logging long flag has been switched off")
			log.Println("This msg through log should not have long format")
		} else if cmd == "on" {
			f := log.Flags() &^ log.Lshortfile
			f = f | log.Llongfile
			log.SetFlags(f)
			c.Println("Logging long flag has been switched on")
			log.Println("This msg through log should have long format")
		} else {
			c.Println("Invalid cmd. Syntax:", SetLogHelp)
		}
	case "short":
		if cmd == "off" {
			f := log.Flags() &^ log.Lshortfile
			log.SetFlags(f)
			c.Println("Logging short flag has been switched off")
			log.Println("This msg through log should not have short format")
		} else if cmd == "on" {
			f := log.Flags() | log.Lshortfile
			log.SetFlags(f)
			c.Println("Logging short flag has been switched on")
			log.Println("This msg through log should have short format")
		} else {
			c.Println("Invalid cmd. Syntax:", SetLogHelp)
		}
	default:
		c.Println("Invalid cmd. Syntax:", SetLogHelp)
	}
}

func createSetMsg(path string, paramName string, paramVal string, msgId string) ([]byte, error) {
	var pss []*usp_msg.Set_UpdateParamSetting
	ps := usp_msg.Set_UpdateParamSetting{
		Param:    paramName,
		Value:    paramVal,
		Required: true,
	}
	pss = append(pss, &ps)
	var objs []*usp_msg.Set_UpdateObject
	obj := usp_msg.Set_UpdateObject{}
	obj.ObjPath = path
	obj.ParamSettings = pss
	objs = append(objs, &obj)
	return parser.CreateUspSetReqMsg(objs, msgId)
}

func createAddMsg(path string, paramName string, paramVal string, msgId string) ([]byte, error) {
	var pss []*usp_msg.Add_CreateParamSetting
	ps := usp_msg.Add_CreateParamSetting{
		Param:    paramName,
		Value:    paramVal,
		Required: true,
	}
	pss = append(pss, &ps)
	var objs []*usp_msg.Add_CreateObject
	obj := usp_msg.Add_CreateObject{}
	obj.ObjPath = path
	obj.ParamSettings = pss
	objs = append(objs, &obj)
	return parser.CreateUspAddReqMsg(objs, msgId)
}
