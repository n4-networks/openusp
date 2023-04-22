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
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/abiosoft/ishell"
)

type MsgType byte

const (
	MsgTypeNotify MsgType = iota
	MsgTypeGet
	MsgTypeSet
	MsgTypeOperate
	MsgTypeAdd
	MsgTypeDel
	MsgTypeGetDm
)

var (
	notifyCount  int
	getCount     int
	setCount     int
	addCount     int
	delCount     int
	getDmCount   int
	defaultCount int
)

type cmdPathInfo struct {
	rootPath         string
	instPath         string
	dmPath           string
	regexPath        string
	hasMultiInstance bool
	level            int
	cmd              string
	params           map[string]string
}

func parseArgs(root string, args []string) (*cmdPathInfo, error) {
	// example show ip <intf|port|twamp> [id] <stats|ipv4addr|ipv6addr|ipv6prefix|twamp> [id]",

	argLen := len(args)
	log.Println("args:", args)
	log.Println("ArgLen:", argLen)

	key := root
	rootObj, ok := nodeSchema[key]
	if !ok {
		return nil, errors.New("Wrong object type")
	}
	var rootPath, regexPath, dmPath, instPath string
	var level int // object level in the TR181 path hierarchy
	if root == "device" {
		rootPath = rootObj.pathName
		regexPath = rootPath
		level = 0
	} else {
		rootPath = "Device." + rootObj.pathName
		regexPath = "Device\\." + rootObj.pathName
		level = 1
	}
	dmPath = rootPath
	instPath = rootPath

	var hasInst bool = false
	var params map[string]string
	var cmd string
	var hasMulti bool = false

	for i := 0; i < argLen; i++ {
		key = key + "." + args[i]
		o, ok := nodeSchema[key]
		if !ok {
			// Check if the next arg is a cmd
			log.Println("args:", args[i])
			var okk bool
			if cmd, okk = cmdSchema[args[i]]; okk {
				log.Println("Its a command:", cmd)
				i++
				//return nil, errors.New("Its a command")
			}
			params = make(map[string]string)
			for j := i; j < argLen; j = j + 2 {
				if (j + 1) < argLen {
					params[args[j]] = args[j+1]
				}
			}
			break
			//return nil, errors.New("Wrong sub-object type")
		}
		regexPath = regexPath + "\\." + o.pathName
		dmPath = dmPath + "." + o.pathName
		instPath = instPath + "." + o.pathName
		level++
		//log.Println("path:", instPath)
		if o.multiInstance {
			//dmPath = dmPath + ".{i}"
			hasMulti = true
			hasInst = false
			if i < argLen-1 {
				hasInst = isDigit(args[i+1])
			}
			if hasInst {
				regexPath = regexPath + "\\." + args[i+1]
				instPath = instPath + "." + args[i+1]
				level++
				i++
			} else {
				regexPath = regexPath + "\\." + "\\d*"
				//instPath = instPath + "." + "*"
				level++
			}
		}
		//log.Println("After multi check path:", instPath)
	}
	regexPath = regexPath + "\\."
	instPath = instPath + "."
	rootPath = rootPath + "."
	dmPath = dmPath + "."
	c := &cmdPathInfo{
		instPath:         instPath,
		rootPath:         rootPath,
		regexPath:        regexPath,
		dmPath:           dmPath,
		hasMultiInstance: hasMulti,
		params:           params,
		cmd:              cmd,
		level:            level,
	}
	log.Printf("RegexPath: %v, DmPath: %v InstPath: %v\n", c.regexPath, c.dmPath, c.instPath)
	log.Printf("RootPath: %v, level: %v Cmd: %v\n", c.rootPath, c.level, c.cmd)
	log.Printf("Params: %+v\n", c.params)

	return c, nil
}

func getMapFromArgs(args []string) (map[string]string, error) {

	if len(args) < 2 {
		log.Println("Insufficient inputs to parse")
		return nil, errors.New("Insufficient inputs to prase")
	}
	var argsMap map[string]string
	argsMap = make(map[string]string)
	argLen := len(args)
	for i := 0; i < argLen-1; i = i + 2 {
		if args[i] == "id" {
			argsMap[args[i-2]+"_id"] = args[i+1]
		} else {
			argsMap[args[i]] = args[i+1]
		}
	}
	return argsMap, nil
}

func (cli *Cli) checkDefault() error {
	if !cli.agent.isSet.epId {
		return errors.New("Agent EndpointId is not set, use set agentid")
	}
	return nil
}

func (cli *Cli) checkCfgDevSet() error {
	if !cli.agent.isSet.productClass {
		return errors.New("Error: DeviceType: ProductClass Type is not defined")
	}
	if !cli.agent.isSet.manufacturer {
		return errors.New("Error: DeviceType: Manufacturer Name is not defined")
	}
	if !cli.agent.isSet.modelName {
		return errors.New("Error: DeviceType: Model Name is not defined")
	}
	return nil
}

func getPath(args []string) string {

	if len(args) == 0 {
		return "Device."
	} else {
		a := args[0]
		if a[len(a)-1:] == "." {
			return a
		}
		return a + "."
	}
}

func getMsgId(t MsgType) string {
	switch t {
	case MsgTypeGet:
		getCount++
		return "GET_" + strconv.Itoa(getCount)
	case MsgTypeGetDm:
		getDmCount++
		return "GETDM_" + strconv.Itoa(getDmCount)
	case MsgTypeSet:
		setCount++
		return "SET_" + strconv.Itoa(setCount)
	case MsgTypeAdd:
		addCount++
		return "ADD_" + strconv.Itoa(addCount)
	case MsgTypeDel:
		delCount++
		return "DEL_" + strconv.Itoa(delCount)
	case MsgTypeNotify:
		notifyCount++
		return "NOTIFY_" + strconv.Itoa(notifyCount)
	default:
		defaultCount++
		return "DEFAULT_" + strconv.Itoa(defaultCount)
	}
	return "InvalidMsgId"
}

type cliParamSchema struct {
	name         string
	isDigitValid bool
	def          string
	validList    []string
}
type cliParam struct {
	value   string
	isDigit bool
	errStr  string
}

func isDigit(str string) bool {
	if _, err := strconv.ParseInt(str, 0, 64); err != nil {
		return false
	}
	return true
}

func getCliParam(param string, s *cliParamSchema) (*cliParam, error) {
	p := &cliParam{}
	if s.isDigitValid {
		p.isDigit = isDigit(param)
		p.value = param
		return p, nil
	}
	for _, v := range s.validList {
		if v == param {
			p.value = param
			p.isDigit = false
			return p, nil
		}
	}
	p.value = s.def // Default value
	p.errStr = fmt.Sprintf("Wrong input, setting %s value to %s", s.name, s.def)
	return p, errors.New("Not found")
}

func (cli *Cli) registerNouns(nouns []noun) {
	var parent *ishell.Cmd
	var ok bool
	var cmdAbsPath string
	for _, n := range nouns {
		cmdAbsPath = n.parent + "." + n.name
		cmd := &ishell.Cmd{
			Name: n.name,
			Help: n.help,
			Func: n.cb,
		}
		parent, ok = cli.sh.cmds[n.parent]
		if !ok {
			log.Printf("Error in adding %v, parent not found:%v\n", n.name, n.parent)
			continue
		}
		parent.AddCmd(cmd)
		cli.sh.cmds[cmdAbsPath] = cmd
	}
}

func (cli *Cli) operateCmd(c *ishell.Context, rootObj string) error {
	if err := cli.checkDefault(); err != nil {
		return err
	}
	if len(c.Args) <= 0 {
		return errors.New("Please provide command name")
	}
	// operate ip intf|port|twamp <id|name> stats|ipv4addr|ipv6addr|ipv6prefix|twamp <id|name> command <cmd> inputs",
	cmdInfo, err := parseArgs(rootObj, c.Args)
	if err != nil {
		return err
	}
	log.Printf("RegexPath: %v, DmPath: %v InstPath: %v\n", cmdInfo.regexPath, cmdInfo.dmPath, cmdInfo.instPath)
	log.Printf("Cmd:  %v\n", cmdInfo.cmd)
	log.Printf("Params:  %v\n", cmdInfo.params)

	// TODO: retrieve command from arg and pass it here
	/*
		cmd := cmdInfo.instPath + cmdInfo.cmd
		if err := cli.MtpOperateReq(cmd, "none", true, cmdInfo.params); err != nil {
			return err
		}
	*/
	return nil
}

func (cli *Cli) showParams(c *ishell.Context, rootObj string) error {
	if err := cli.checkDefault(); err != nil {
		return err
	}
	// show ip <intf|port|twamp> [id] <stats|ipv4addr|ipv6addr|ipv6prefix|twamp> [id]",
	cmdInfo, err := parseArgs(rootObj, c.Args)
	if err != nil {
		return err
	}
	log.Printf("RegexPath: %v, DmPath: %v InstPath: %v\n", cmdInfo.regexPath, cmdInfo.dmPath, cmdInfo.instPath)
	objParams, err := cli.restReadParams(cmdInfo.instPath)
	if err != nil {
		log.Println("Err:", err)
		return err
	}
	for _, obj := range objParams {
		c.Printf("%-25s : %-12s\n", "Object Path", obj.Path)
		for _, p := range obj.Params {
			c.Printf(" %-24s : %-12s\n", p.Name, p.Value)
		}
		c.Println("-------------------------------------------------")
	}
	return nil
}

type objPathInfo struct {
	dm    string
	regex string
	inst  string
}

func (cli *Cli) showCfg(c *ishell.Context, rootObj string) error {
	if err := cli.checkCfgDevSet(); err != nil {
		return err
	}
	/*
		cmdInfo, err := parseArgs(rootObj, c.Args)
		if err != nil {
			return err
		}
				insts, err1 := cli.dbGetCfgInstancesByRegex(cmdInfo.instPath)
				if err1 != nil {
					return err1
				}
			c.Println("-------------------------------------------------")
			for _, inst := range insts {
				c.Printf("%-25s : %-12s\n", "Cfg Instance Path", inst.path)
				c.Printf("%-25s : %-12s\n", "Key", inst.key)
				for k, v := range inst.params {
					c.Printf(" %-24s : %-12s\n", k, v)
				}
			}
			paramNodes, err2 := cli.dbGetCfgParamNodesByRegex(cmdInfo.instPath)
			if err2 != nil {
				return err2
			}
			for _, paramNode := range paramNodes {
				c.Printf("%-25s : %-12s\n", "Cfg Param Path", paramNode.path)
				for k, v := range paramNode.params {
					c.Printf(" %-24s : %-12s\n", k, v)
				}
			}
	*/
	c.Println("-------------------------------------------------")
	return nil
}

func (cli *Cli) addInst(p *addInstInfo, dest destType) (string, error) {
	printAddInstInfo(p)
	if dest == destMtp {
		inst, err := cli.restAddInstance(p.path, p.params)
		if err != nil {
			log.Println("RestErr:", err)
			return "", err
		}
		// update db with the latest obj params
		if err := cli.restUpdateParams(p.parent); err != nil {
			return "", err
		}
		return inst.Path, nil
	}
	if dest == destDb {
		nodeLevel := strings.Count(p.path, ".")
		log.Println("Adding instance to DB with level:", nodeLevel)
		/*
			key := p.params["Alias"]
				if err := cli.dbWriteCfgInstance(p.path, nodeLevel, key, p.params); err != nil {
					return "", err
				}
		*/
	}
	return "", nil
}

type removeCmdPathInfo struct {
	rootPath string
	instPath string
}

func parseRemoveArgs(root string, args []string) (*removeCmdPathInfo, error) {
	argLen := len(args)
	log.Println("argLen:", argLen)

	if argLen <= 0 {
		return nil, errors.New("object name not provided")
	}
	rootNode, ok := nodeSchema[root]
	path := "Device." + rootNode.pathName + "."
	parentPath := path

	key := root
	var obj node
	var arg string

	for argCnt := 0; argCnt < argLen; argCnt++ {
		//log.Printf("key: %v path: %v\n", key, path)
		arg = args[argCnt]
		key = key + "." + arg

		if obj, ok = nodeSchema[key]; !ok {
			return nil, errors.New("object not found")
		} else {
			path = path + obj.pathName + "."
		}

		log.Println("argCnt:", argCnt)
		if obj.multiInstance {
			argCnt++
			if argCnt >= argLen {
				return nil, errors.New("instance id|name not found")
			}

			arg = args[argCnt]
			if isDigit(arg) || arg == "*" {
				path = path + arg + "."
			} else {
				path = path + "[Alias==\"" + arg + "\"]."
			}
		}
	}
	c := &removeCmdPathInfo{}
	c.instPath = path
	c.rootPath = parentPath
	log.Printf("Instance Path: %v, root path: %v\n", path, parentPath)

	return c, nil
}

type removeCfgCmdPathInfo struct {
	instPath string
	alias    string
}

// removecfg ip <intf|port> <alias> <ipv4addr|ipv6addr> <alias>
func parseRemoveCfgArgs(root string, args []string) (*removeCfgCmdPathInfo, error) {
	argLen := len(args)
	log.Println("argLen:", argLen)

	if argLen < 2 {
		return nil, errors.New("object name or id/name not provided")
	}
	var obj node
	var ok bool
	path := "Device."
	key := root
	obj, ok = nodeSchema[key]
	if !ok {
		return nil, errors.New("invalid root object")
	}
	path = path + obj.pathName + "."

	for i := 0; i < argLen-2; i++ {
		key = key + "." + args[i]
		obj, ok = nodeSchema[key]
		if !ok {
			return nil, errors.New("invalid object name")
		}
		if i < argLen-1 {
			id := args[i+1]
			if isDigit(id) {
				path = path + obj.pathName + "." + id + "."
			} else {
				path = path + obj.pathName + "." + "[Alias==\"" + args[i+1] + "\"]."
			}
			i++
		} else {
			return nil, errors.New("Alias not provided")
		}
	}
	key = key + "." + args[argLen-2]
	log.Println("key:", key)
	obj, ok = nodeSchema[key]
	if !ok {
		return nil, errors.New("Invalid last object name")
	}
	path = path + obj.pathName + "."

	c := &removeCfgCmdPathInfo{}
	c.instPath = path
	c.alias = args[argLen-1]

	return c, nil
}

func (cli *Cli) removeInst(c *ishell.Context, rootObj string) error {
	if err := cli.checkDefault(); err != nil {
		return err
	}
	cmd, err := parseRemoveArgs(rootObj, c.Args)
	if err != nil {
		return err
	}
	log.Println("Parsed all the path")
	log.Println("Initiating grpc")

	if err = cli.restDeleteInstances(cmd.instPath); err != nil {
		return err
	}
	// update the param collection for the same path
	if err = cli.restUpdateParams(cmd.rootPath); err != nil {
		return err
	}
	c.Println("Instance removed from:", cmd.instPath)
	return nil
}

// removecfg ip <intf|port> <id|name> <ipv4addr|ipv6addr> <id|name>
func (cli *Cli) removeCfgInst(c *ishell.Context, rootObj string) error {
	if err := cli.checkCfgDevSet(); err != nil {
		return err
	}
	pathInfo, err := parseRemoveCfgArgs(rootObj, c.Args)
	if err != nil {
		return err
	}
	log.Println("path:", pathInfo.instPath)
	log.Println("alias:", pathInfo.alias)
	/*
		if err := cli.dbDeleteCfgInstanceByKey(pathInfo.instPath, pathInfo.alias); err != nil {
			return err
		}
	*/

	c.Println("Instance removed from:", pathInfo.instPath)
	return nil
}

func (cli *Cli) updateDb(args []string, rootObj string) error {
	if err := cli.checkDefault(); err != nil {
		return err
	}
	cmdInfo, err := parseArgs(rootObj, args)
	if err != nil {
		return err
	}
	/*
		if strings.Contains(cmdInfo.instPath, "*") {
			return errors.New("Multiinstance object, instance id not found")
		}
	*/
	log.Printf("path: %v\n", cmdInfo.instPath)

	if err = cli.restUpdateParams(cmdInfo.instPath); err != nil {
		return err
	}

	if err = cli.restUpdateInstances(cmdInfo.rootPath); err != nil {
		return err
	}

	return nil
}

func (cli *Cli) setParam(args []string, rootObj string) error {
	if err := cli.checkDefault(); err != nil {
		return err
	}
	argLen := len(args)
	if argLen < 3 {
		return errors.New("Wrong input")
	}
	cmdInfo, err := parseArgs(rootObj, args[:argLen-2])
	if err != nil {
		return err
	}
	if strings.Contains(cmdInfo.instPath, "*") {
		return errors.New("Multiinstance object, instance id not found")
	}
	param := args[argLen-2]
	value := args[argLen-1]
	log.Printf("path: %v param: %v value: %v\n", cmdInfo.instPath, param, value)

	params := map[string]string{param: value}
	if err := cli.restSetParams(cmdInfo.instPath, params); err != nil {
		return err
	}
	return nil
}

func (cli *Cli) setCfgParam(args []string, rootObj string) error {
	if err := cli.checkCfgDevSet(); err != nil {
		return err
	}
	/*
		argLen := len(args)
		if argLen < 2 || argLen%2 != 0 {
			return errors.New("Wrong input")
		}
	*/
	cmdInfo, err := parseArgs(rootObj, args) //args[:argLen-2])
	if err != nil {
		return err
	}
	if strings.Contains(cmdInfo.instPath, "*") {
		return errors.New("Multi-instance object, instance id not found")
	}
	log.Printf("path: %v params: %+v\n", cmdInfo.instPath, cmdInfo.params)

	/*
		if params, err := cli.dbGetCfgParams(cmdInfo.instPath); err == nil {
			for k, v := range params {
				cmdInfo.params[k] = v
			}
		}

		if err := cli.dbWriteCfgParamNode(cmdInfo.instPath, cmdInfo.params); err != nil {
			return err
		}
	*/
	return nil
}
func printAddInstInfo(info *addInstInfo) {
	log.Println("Path:", info.path)
	log.Println("Parent:", info.parent)
	log.Printf("Param:%+v\n", info.params)
}
