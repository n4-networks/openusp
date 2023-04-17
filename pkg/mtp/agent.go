package mtp

import (
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/n4-networks/usp/pkg/parser"
	"github.com/n4-networks/usp/pkg/pb/bbf/usp_msg"
)

type agentMtpIntf interface {
	sendMsg([]byte) error
	getMsgCnt() uint64
	incMsgCnt()
}

type agentInitData struct {
	epId    string
	mtpIntf agentMtpIntf
	params  map[string]string
}

type agentHandler struct {
	mtpMap    map[string]agentMtpIntf
	rxChannel chan agentInitData
}

func (m *Mtp) agentHandlerInit() error {
	m.agentH.mtpMap = make(map[string]agentMtpIntf, 1000) // max map size of 1000 agent at a time
	m.agentH.rxChannel = make(chan agentInitData, 100)    // buffered channel with queue size of 100
	return nil
}

func (m *Mtp) agentInitThread(initData *agentInitData) {
	agentId := initData.epId
	params := initData.params
	mtpIntf := initData.mtpIntf

	// Clear param and instance table
	path := "Device."
	if err := m.dbDeleteInstancesByRegex(agentId, path); err != nil {
		log.Println("Could not clear instances data after agent's boot for agent:", agentId)
		return
	}
	paths := []string{"Device."}
	if err := m.dbDeleteParams(agentId, paths); err != nil {
		log.Println("Could not clear params data after agent's boot for agent:", agentId)
		return
	}

	var err error
	var uspMsg []byte
	var msgId string

	// send get datamodel to agent
	msgId = "AGENT_REINIT_GET_DM_" + strconv.FormatUint(mtpIntf.getMsgCnt(), 10)
	uspMsg, err = parser.CreateUspGetSupportedDmMsg(paths, true, true, true, msgId)
	if err != nil {
		log.Println("Could not prepare Usp msg of type: GET DATAMODEL")
		return
	}
	if err := m.sendUspMsgToAgent(agentId, uspMsg, mtpIntf); err != nil {
		log.Println("Error in sending GET DATAMODEL msg to agent:", agentId)
		return
	}

	// send add instance requests
	log.Println("Setting default configs from database")
	//paramMap, _ := strToMapWithTwoDelims(n.evt.params["ParameterMap"], ",", ":")
	if err := m.setAgentDefaultConfig(agentId, params, mtpIntf); err != nil {
		log.Println("Error in setting default config", err)
	}

	// send get param to agent
	msgId = "AGENT_REINIT_GET_PARAM_" + strconv.FormatUint(mtpIntf.getMsgCnt(), 10)
	if uspMsg, err = parser.CreateUspGetReqMsg(paths, msgId); err != nil {
		log.Println("Could not prepare Usp msg of type: GET ")
		return
	}
	if err := m.sendUspMsgToAgent(agentId, uspMsg, mtpIntf); err != nil {
		log.Println("Error in sending GET msg to agent:", agentId)
		return
	}

	// get instance request to the agent
	msgId = "AGENT_REINIT_GET_INSTANCE_" + strconv.FormatUint(mtpIntf.getMsgCnt(), 10)
	if uspMsg, err = parser.CreateUspGetInstancesMsg(paths, false, msgId); err != nil {
		log.Println("Could not prepare Usp msg of type: GET Instances")
		return
	}
	if err := m.sendUspMsgToAgent(agentId, uspMsg, mtpIntf); err != nil {
		log.Println("Error in sending GET Instances msg to agent:", agentId)
		return
	}
	log.Println("Adding mtpMap for agent:", agentId)
	m.agentH.mtpMap[agentId] = mtpIntf
}

func (m *Mtp) sendUspMsgToAgent(agentId string, uspMsg []byte, mtpIntf agentMtpIntf) error {
	controllerId := m.Cfg.Usp.EndpointId
	uspRecord, err := parser.CreateNewPlainTextRecord(&agentId, &controllerId, nil, nil, uspMsg)
	if err != nil {
		log.Println("Could not convert USP msg to USP record: ", err)
		return err
	}
	if err := mtpIntf.sendMsg(uspRecord); err != nil {
		return err
	}
	mtpIntf.incMsgCnt()

	return nil
}

func (m *Mtp) setAgentDefaultConfig(agentId string, notifyParams map[string]string, mtpIntf agentMtpIntf) error {
	// Read cfg instance table
	var devInfo agentDeviceInfo
	var msgId string

	devInfo.productClass = notifyParams["Device.DeviceInfo.ProductClass"]
	devInfo.manufacturer = notifyParams["Device.DeviceInfo.ManufacturerOUI"]
	devInfo.modelName = notifyParams["Device.DeviceInfo.SerialNumber"]
	insts, err := m.dbGetCfgInstances(&devInfo)
	if err != nil {
		log.Println(err)
		return err
	}

	for _, inst := range insts {
		var createObjs []*usp_msg.Add_CreateObject
		createObj := &usp_msg.Add_CreateObject{}
		createObj.ObjPath = inst.path
		log.Println("Default Config adding instance for:", inst.path)
		for key, val := range inst.params {
			createParam := &usp_msg.Add_CreateParamSetting{}
			createParam.Param = key
			createParam.Value = val
			log.Printf("Default Config add param %v : %v\n", key, val)
			createParam.Required = true
			createObj.ParamSettings = append(createObj.ParamSettings, createParam)
		}
		createObjs = append(createObjs, createObj)
		//msgId = "DEF_CFG_ADD_INSTANCE_" + strconv.FormatInt(int64(inst.level), 10)
		msgId = "AGENT_DEFCFG_GET_INSTANCE_" + strconv.FormatUint(mtpIntf.getMsgCnt(), 10)

		uspMsg, err := parser.CreateUspAddReqMsg(createObjs, msgId)
		if err != nil {
			log.Println("Could not prepare Usp msg of type: ADD")
			continue
		}
		if err := m.sendUspMsgToAgent(agentId, uspMsg, mtpIntf); err != nil {
			log.Println("Error in sending Add Instance msg to agent")
		}
	}

	log.Println("Default config: set parameters")

	paramNodes, err := m.dbGetCfgParamNodes(&devInfo)
	if err != nil {
		log.Println(err)
		return err
	}
	var objs []*usp_msg.Set_UpdateObject
	for _, paramNode := range paramNodes {
		obj := &usp_msg.Set_UpdateObject{}
		obj.ObjPath = paramNode.path
		log.Println("Default config set operation path:", paramNode.path)

		var pss []*usp_msg.Set_UpdateParamSetting
		for name, val := range paramNode.params {
			ps := &usp_msg.Set_UpdateParamSetting{
				Param:    name,
				Value:    val,
				Required: true,
			}
			pss = append(pss, ps)
			log.Printf("Default config set param %v : %v\n", name, val)
		}
		obj.ParamSettings = pss
		objs = append(objs, obj)
	}
	msgId = "AGENT_DEFCFG_SET_PARAM_" + strconv.FormatUint(mtpIntf.getMsgCnt(), 10)
	uspMsg, err := parser.CreateUspSetReqMsg(objs, msgId)
	if err != nil {
		log.Println("Could not prepare Usp msg of type: SET")
		return err
	}
	if err := m.sendUspMsgToAgent(agentId, uspMsg, mtpIntf); err != nil {
		log.Println("Error in sending Add Instance msg to agent")
		return err
	}
	return nil
}

func strToMapWithTwoDelims(s string, delim1 string, delim2 string) (map[string]string, error) {
	tok := strings.Split(strings.Replace(strings.Trim(s, "{}"), "\"", "", -1), delim1)
	m := make(map[string]string)
	for _, s1 := range tok {
		t := strings.Split(s1, delim2)
		m[t[0]] = t[1]
	}
	return m, nil
}

func (m *Mtp) getAgentMtp(epId string) (agentMtpIntf, error) {
	// return agentMtpIntf if available in agent handle
	if mtpIntf, ok := m.agentH.mtpMap[epId]; ok {
		log.Println("Found agent MTP intf in agent handler")
		return mtpIntf, nil
	}
	// Get Agent MTP information from database
	log.Println("Finding out agent MTP info from DB")
	path := "Device.LocalAgent.MTP."
	instances, err := m.dbGetInstancesByRegex(epId, path)
	if err != nil {
		return nil, err
	}
	log.Println("Number of agent MTP instances found:", len(instances))

	paramMap, err := m.dbGetParamsByRegex(epId, path)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	//log.Printf("paramMap:%+v\n", paramMap)

	for _, inst := range instances {
		//log.Println("instPath:", inst.path)
		//log.Println("paramMap[Enable]:", paramMap[inst.path+"Enable"])
		//log.Println("paramMap[Protocol]:", paramMap[inst.path+"Protocol"])
		if paramMap[inst.path+"Enable"] == "true" {
			switch paramMap[inst.path+"Protocol"] {
			case "STOMP":
				aStomp := &agentStomp{}
				aStomp.conn = m.connH.stomp.Conn
				aStomp.destQueue = paramMap[inst.path+"STOMP.Destination"]
				m.agentH.mtpMap[epId] = aStomp
				return aStomp, nil
			case "CoAP":
				aCoap := &agentCoap{}
				aCoap.port = paramMap[inst.path+"CoAP.Port"]
				aCoap.path = paramMap[inst.path+"CoAP.Path"]
				aCoap.isEncrypted = paramMap[inst.path+"CoAP.IsEncrypted"]
				return aCoap, nil
			}
		}
	}
	log.Println("Could not resolve agent mtp info from DB")
	return nil, errors.New("Agent MTP not found")
}
