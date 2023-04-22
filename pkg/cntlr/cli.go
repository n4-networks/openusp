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

package cntlr

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

func (c *Cntlr) Cli() {
	shell := ishell.New()
	shell.Println("OpenUSP Cntlr Cli")
	shell.SetPrompt("OpenUSP Cntlr>>")

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
		Func: c.cliSendGetMsg,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "getinstance",
		Help: "getinstance firstlevelonly (true/false), objpath",
		Func: c.cliSendGetInstanceMsg,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "setparam",
		Help: "set parameter value. ex: set path param value",
		Func: c.cliSendSetMsg,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "addinstance",
		Help: "add instance of an object and its parameters. ex: add objpath param value",
		Func: c.cliSendAddMsg,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "operate",
		Help: "operate msg to an USP Agent. ex: operate Device.Reboot()",
		Func: c.cliSendOperateMsg,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "getdm",
		Help: "getdm msg to an USP Agent. ex: getdm Device.DeviceInfo.",
		Func: c.cliSendGetSupportedDmMsg,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "delinstance",
		Help: "delete an instance of an multi-instance object. ex: delinstance Device.DeviceInfo.",
		Func: c.cliSendDeleteInstanceMsg,
	})

	/* Shibu: TODO
	shell.AddCmd(&ishell.Cmd{
		Name: "getresp",
		Help: "receive respone from agent, ex: getresp",
		Func: func(cli *ishell.Context) {
			c.StompReceiveUspMsgFromAgentWithTimer(1)
		},
	})
	*/

	shell.AddCmd(&ishell.Cmd{
		Name: "getver",
		Help: "Print MTP version",
		Func: func(cli *ishell.Context) {
			cli.Println("MTP Version:", getVer())
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "setlog",
		Help: SetLogHelp,
		Func: c.cliSetLog,
	})
	shell.Run()
}

func (c *Cntlr) cliSendGetMsg(cli *ishell.Context) {
	mtpIntf, err := c.getAgentMtp(cliAgentId)
	if err != nil {
		cli.Println("Could not get agent MTP interface")
		return
	}
	msgId := "MTPCLI_GET_PARAM_" + strconv.FormatUint(mtpIntf.GetMsgCnt(), 10)
	uspMsg, err := parser.CreateUspGetReqMsg(cli.Args, msgId)
	if err != nil {
		cli.Println("Could not prepare Usp msg of type: GET")
		return
	}
	if err := c.sendUspMsgToAgent(cliAgentId, uspMsg, mtpIntf); err != nil {
		cli.Println("Error in sending get msg to agent")
	}
}

func (c *Cntlr) cliSendGetInstanceMsg(cli *ishell.Context) {
	var firstLevelOnly bool
	switch cli.Args[0] {
	case "false":
		firstLevelOnly = false
	case "true":
		firstLevelOnly = true
	default:
		cli.Println("getinstance firstlevelonly(true/false) objpath")
		return
	}
	mtpIntf, err := c.getAgentMtp(cliAgentId)
	if err != nil {
		cli.Println("Could not get agent MTP interface")
		return
	}
	msgId := "MTPCLI_GET_INSTANCE_" + strconv.FormatUint(mtpIntf.GetMsgCnt(), 10)
	uspMsg, err := parser.CreateUspGetInstancesMsg(cli.Args[1:], firstLevelOnly, msgId)
	if err != nil {
		cli.Println("Could not prepare Usp  msg of type: GET_INSTANCES")
		return
	}
	if err := c.sendUspMsgToAgent(cliAgentId, uspMsg, mtpIntf); err != nil {
		cli.Println("Error in sending get msg to agent")
	}
}

func (c *Cntlr) cliSendSetMsg(cli *ishell.Context) {
	if len(cli.Args) < 3 {
		cli.Println("Error: Insufficient number of args (path/param/val)")
		return
	}
	mtpIntf, err := c.getAgentMtp(cliAgentId)
	if err != nil {
		cli.Println("Could not get agent MTP interface")
		return
	}
	msgId := "MTPCLI_SET_PARAM__" + strconv.FormatUint(mtpIntf.GetMsgCnt(), 10)
	uspMsg, err := createSetMsg(cli.Args[0], cli.Args[1], cli.Args[2], msgId)
	if err != nil {
		cli.Println("Could not prepare Usp  msg of type: SET")
		return
	}
	if err := c.sendUspMsgToAgent(cliAgentId, uspMsg, mtpIntf); err != nil {
		cli.Println("Error in sending get msg to agent")
	}
}

func (c *Cntlr) cliSendAddMsg(cli *ishell.Context) {
	if len(cli.Args) < 3 {
		cli.Println("Error: Insufficient number of args: cmd objpath paramName value)")
		return
	}
	mtpIntf, err := c.getAgentMtp(cliAgentId)
	if err != nil {
		cli.Println("Could not get agent MTP interface")
		return
	}
	msgId := "MTPCLI_ADD_INSTANCE__" + strconv.FormatUint(mtpIntf.GetMsgCnt(), 10)
	uspMsg, err := createAddMsg(cli.Args[0], cli.Args[1], cli.Args[2], msgId)
	if err != nil {
		cli.Println("Could not prepare Usp  msg of type: ADD")
		return
	}
	if err := c.sendUspMsgToAgent(cliAgentId, uspMsg, mtpIntf); err != nil {
		cli.Println("Error in sending get msg to agent")
	}
}

func (c *Cntlr) cliSendOperateMsg(cli *ishell.Context) {
	cmd := cli.Args[0] //[string("Device.Reboot()")
	cmdKey := "none"
	resp := true

	mtpIntf, err := c.getAgentMtp(cliAgentId)
	if err != nil {
		cli.Println("Could not get agent MTP interface")
		return
	}
	msgId := "MTPCLI_OPERATE_" + strconv.FormatUint(mtpIntf.GetMsgCnt(), 10)
	uspMsg, err := parser.CreateUspOperateReqMsg(cmd, cmdKey, resp, msgId, nil)
	if err != nil {
		cli.Println("Could not prepare Usp msg of type: OPERATE")
		return
	}
	if err := c.sendUspMsgToAgent(cliAgentId, uspMsg, mtpIntf); err != nil {
		cli.Println("Error in sending get msg to agent")
	}
}

func (c *Cntlr) cliSendGetProtoVersionMsg(cli *ishell.Context) {
	mtpIntf, err := c.getAgentMtp(cliAgentId)
	if err != nil {
		cli.Println("Could not get agent MTP interface")
		return
	}
	msgId := "MTPCLI_GET_PROTO_VER_" + strconv.FormatUint(mtpIntf.GetMsgCnt(), 10)
	uspMsg, err := parser.CreateUspGetSupportedProtoMsg(c.cfg.usp.protoVersion, msgId)
	if err != nil {
		cli.Println("Could not prepare Usp msg of type: GET_PROTO_VERSION")
		return
	}
	if err := c.sendUspMsgToAgent(cliAgentId, uspMsg, mtpIntf); err != nil {
		cli.Println("Error in sending get msg to agent")
	}
}

func (c *Cntlr) cliSendDeleteInstanceMsg(cli *ishell.Context) {
	mtpIntf, err := c.getAgentMtp(cliAgentId)
	if err != nil {
		cli.Println("Could not get agent MTP interface")
		return
	}
	msgId := "MTPCLI_DEL_INSTANCE_" + strconv.FormatUint(mtpIntf.GetMsgCnt(), 10)
	uspMsg, err := parser.CreateUspDeleteReqMsg(cli.Args, msgId)
	if err != nil {
		cli.Println("Could not prepare Usp msg of type: DELETE")
		return
	}
	if err := c.sendUspMsgToAgent(cliAgentId, uspMsg, mtpIntf); err != nil {
		cli.Println("Error in sending get msg to agent")
	}
}

func (c *Cntlr) cliSendGetSupportedDmMsg(cli *ishell.Context) {
	mtpIntf, err := c.getAgentMtp(cliAgentId)
	if err != nil {
		log.Println("Could not get agent mtpIntf")
		return
	}
	msgId := "MTPCLI_GET_DM_" + strconv.FormatUint(mtpIntf.GetMsgCnt(), 10)
	uspMsg, err := parser.CreateUspGetSupportedDmMsg(cli.Args, true, true, true, msgId)
	if err != nil {
		cli.Println("Could not prepare Usp msg of type: GET_SUPPORTED_DM")
		return
	}
	if err := c.sendUspMsgToAgent(cliAgentId, uspMsg, mtpIntf); err != nil {
		cli.Println("Error in sending get msg to agent")
	}
}

const SetLogHelp = "setlog <on|off> <all|date|time|long|short>"

func (c *Cntlr) cliSetLog(cli *ishell.Context) {
	if len(cli.Args) < 2 {
		cli.Println("setlog help")
		return
	}
	cmd := cli.Args[0]
	flag := cli.Args[1]
	switch flag {
	case "all":
		if cmd == "off" {
			log.SetFlags(0)
			log.SetOutput(ioutil.Discard)
			cli.Println("Switched off all logging flags")
			log.Println("If you are seeing this then something is broken...")
		} else if cmd == "on" {
			log.SetFlags(log.Lshortfile | log.Llongfile | log.Ldate | log.Ltime)
			log.SetOutput(os.Stdout)
			cli.Println("Switched on all logging flags")
			log.Println("This message from logging engine, its has been switched on now")
		} else {
			cli.Println("Invalid cmd. Syntax:", SetLogHelp)
		}
	case "date":
		if cmd == "off" {
			f := log.Flags() &^ log.Ldate
			log.SetFlags(f)
			cli.Println("Logging date flag has been switched off")
			log.Println("This msg through log should not have date")
		} else if cmd == "on" {
			f := log.Flags() | log.Ldate
			log.SetFlags(f)
			cli.Println("Logging date flag has been switched on")
			log.Println("This msg through log should have date")
		} else {
			cli.Println("Invalid cmd. Syntax:", SetLogHelp)
		}
	case "time":
		if cmd == "off" {
			f := log.Flags() &^ log.Ltime
			log.SetFlags(f)
			cli.Println("Logging time flag has been switched off")
			log.Println("This msg through log should not have time")
		} else if cmd == "on" {
			f := log.Flags() | log.Ltime
			log.SetFlags(f)
			cli.Println("Logging time flag has been switched on")
			log.Println("This msg through log should have time")
		} else {
			cli.Println("Invalid cmd. Syntax:", SetLogHelp)
		}
	case "long":
		if cmd == "off" {
			f := log.Flags() &^ log.Llongfile
			log.SetFlags(f)
			cli.Println("Logging long flag has been switched off")
			log.Println("This msg through log should not have long format")
		} else if cmd == "on" {
			f := log.Flags() &^ log.Lshortfile
			f = f | log.Llongfile
			log.SetFlags(f)
			cli.Println("Logging long flag has been switched on")
			log.Println("This msg through log should have long format")
		} else {
			cli.Println("Invalid cmd. Syntax:", SetLogHelp)
		}
	case "short":
		if cmd == "off" {
			f := log.Flags() &^ log.Lshortfile
			log.SetFlags(f)
			cli.Println("Logging short flag has been switched off")
			log.Println("This msg through log should not have short format")
		} else if cmd == "on" {
			f := log.Flags() | log.Lshortfile
			log.SetFlags(f)
			cli.Println("Logging short flag has been switched on")
			log.Println("This msg through log should have short format")
		} else {
			cli.Println("Invalid cmd. Syntax:", SetLogHelp)
		}
	default:
		cli.Println("Invalid cmd. Syntax:", SetLogHelp)
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
