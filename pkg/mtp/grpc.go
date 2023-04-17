package mtp

import (
	"context"
	"errors"
	"log"
	"net"
	"time"

	"github.com/n4-networks/openusp/pkg/parser"
	"github.com/n4-networks/openusp/pkg/pb/bbf/usp_msg"
	"github.com/n4-networks/openusp/pkg/pb/mtpgrpc"
	"google.golang.org/grpc"
)

type instanceResp struct {
	epId      string
	msgId     string
	instances []*instance
}

func (m *Mtp) GrpcServer(port string, exit chan int32) {
	addr := ":" + port
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	} else {
		log.Printf("Starting Grpc Server at: %s", port)
		grpcServer := grpc.NewServer()
		mtpgrpc.RegisterMtpGrpcServer(grpcServer, m)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Grpc server failed to serve: %v", err)
		}
	}
	log.Fatalf("Grpc Server is exiting...")
	exit <- GRPC_SERVER
}

func (m *Mtp) MtpGetInfo(ctx context.Context, p *mtpgrpc.None) (*mtpgrpc.MtpInfoData, error) {
	ret := &mtpgrpc.MtpInfoData{}
	ret.Version = getVer()
	return ret, nil
}

/* USP releated services */
func (m *Mtp) MtpGetParamReq(ctx context.Context, p *mtpgrpc.MtpGetParamReqData) (*mtpgrpc.MtpReqResult, error) {
	log.Printf("GetParamReqTx: AgetnId: %v, MsgId: %v\n", p.AgentId, p.MsgId)
	var paths []string
	paths = append(paths, p.Path)
	ret := &mtpgrpc.MtpReqResult{IsSuccess: false}

	log.Printf("GetParamReqTx: Path: %v\n", p.Path)
	uspMsg, err := parser.CreateUspGetReqMsg(paths, p.MsgId)
	if err != nil {
		log.Println("Could not prepare Usp  msg of type: SET")
		ret.ErrMsg = err.Error()
		return ret, err
	}
	mtpIntf, err := m.getAgentMtp(p.AgentId)
	if err != nil {
		log.Println("Could not get agentMtp Interface for agent:", p.AgentId)
		ret.ErrMsg = err.Error()
		return ret, err
	}
	if err := m.sendUspMsgToAgent(p.AgentId, uspMsg, mtpIntf); err != nil {
		log.Println("Error in sending get msg to agent")
		ret.ErrMsg = err.Error()
		return ret, err
	}
	if cacheErr, err := m.cacheGetError(p.AgentId, p.MsgId); err == nil {
		ret.ErrMsg = cacheErr.Msg
		return ret, nil
	}
	ret.IsSuccess = true
	return ret, nil
}
func (m *Mtp) MtpSetParamReq(ctx context.Context, p *mtpgrpc.MtpSetParamReqData) (*mtpgrpc.MtpSetParamResData, error) {
	log.Printf("SetParamReqTx: AgetnId: %v, MsgId: %v\n", p.AgentId, p.MsgId)
	uspMsg, err := createSetMsg(p.Path, p.Param, p.Value, p.MsgId)
	log.Printf("SetParamReqRx: Path: %v, Param: %v, Value: %v\n", p.Path, p.Param, p.Value)
	ret := &mtpgrpc.MtpSetParamResData{IsSuccess: false}
	if err != nil {
		log.Println("Could not prepare Usp  msg of type: SET")
		ret.ErrMsg = "Could not prepare USP msg"
		return ret, err
	}
	mtpIntf, err := m.getAgentMtp(p.AgentId)
	if err != nil {
		log.Println("Could not get agentMtp Interface for agent:", p.AgentId)
		ret.ErrMsg = err.Error()
		return ret, err
	}
	log.Println("UspSetParam: Formed stomp msg with agendId:", p.AgentId)
	if err := m.sendUspMsgToAgent(p.AgentId, uspMsg, mtpIntf); err != nil {
		log.Println("Error in sending get msg to agent")
		ret.ErrMsg = "Could not send stomp msg"
		return ret, err
	}
	time.Sleep(time.Second)

	// Read msg error data from Cache
	if cacheErr, err := m.cacheGetError(p.AgentId, p.MsgId); err == nil {
		ret.ErrMsg = cacheErr.Msg
		return ret, nil
	}
	// No hit on Msg error cache which means received add response from agent
	// Read add response data from cache now
	cacheSetRes, err := m.cacheGetParamSetResult(p.AgentId, p.MsgId)
	if err != nil {
		return ret, errors.New("No response from agent")
	}
	ret.IsSuccess = cacheSetRes.OpIsSuccess
	if !cacheSetRes.OpIsSuccess {
		ret.ErrMsg = cacheSetRes.OpErrStr
	}
	ret.AgentId = p.AgentId
	ret.MsgId = p.MsgId
	ret.IsSuccess = true
	return ret, nil
}

func (m *Mtp) MtpGetInstancesReq(ctx context.Context, p *mtpgrpc.MtpGetInstancesReqData) (*mtpgrpc.MtpReqResult, error) {
	log.Printf("GetInstancesReqTx: AgetnId: %v, MsgId: %v\n", p.AgentId, p.MsgId)
	var objPaths []string
	objPaths = append(objPaths, p.Path)
	log.Printf("GetInstancesReqTx: Path: %v\n", p.Path)
	uspMsg, err := parser.CreateUspGetInstancesMsg(objPaths, p.FirstLevelOnly, p.MsgId)
	ret := &mtpgrpc.MtpReqResult{IsSuccess: false}
	if err != nil {
		log.Println("Could not prepare Usp  msg of type: GET INSTANCE")
		ret.ErrMsg = err.Error()
		return ret, err
	}
	mtpIntf, err := m.getAgentMtp(p.AgentId)
	if err != nil {
		log.Println("Could not get agentMtp Interface for agent:", p.AgentId)
		ret.ErrMsg = err.Error()
		return ret, err
	}
	if err := m.sendUspMsgToAgent(p.AgentId, uspMsg, mtpIntf); err != nil {
		log.Println("Error in sending get msg to agent")
		ret.ErrMsg = err.Error()
		return ret, err
	}
	if cacheErr, err := m.cacheGetError(p.AgentId, p.MsgId); err == nil {
		ret.ErrMsg = cacheErr.Msg
		return ret, nil
	}

	ret.IsSuccess = true
	return ret, nil
}

func (m *Mtp) MtpAddInstanceReq(ctx context.Context, req *mtpgrpc.MtpAddInstanceReqData) (*mtpgrpc.MtpAddInstanceResData, error) {
	log.Printf("AddInstanceReqTx: AgetnId: %v, MsgId: %v\n", req.AgentId, req.MsgId)
	var createObjs []*usp_msg.Add_CreateObject
	for _, obj := range req.GetObjs() {
		createObj := &usp_msg.Add_CreateObject{}
		createObj.ObjPath = obj.GetPath()
		log.Printf("AddInstancesReqTx: Path: %v\n", createObj.ObjPath)
		for key, val := range obj.GetParams() {
			createParam := &usp_msg.Add_CreateParamSetting{}
			createParam.Param = key
			createParam.Value = val
			createParam.Required = true
			createObj.ParamSettings = append(createObj.ParamSettings, createParam)
		}
		log.Printf("AddInstancesReqTx: Params: %+v\n", createObj.ParamSettings)
		createObjs = append(createObjs, createObj)
	}

	ret := &mtpgrpc.MtpAddInstanceResData{IsSuccess: false}

	uspMsg, err := parser.CreateUspAddReqMsg(createObjs, req.MsgId)
	if err != nil {
		log.Println("Could not prepare Usp  msg of type: ADD")
		return ret, err
	}
	log.Println("UspSetParam: Formed stomp msg with agendId:", req.AgentId)

	mtpIntf, err := m.getAgentMtp(req.AgentId)
	if err != nil {
		log.Println("Could not get agentMtp Interface for agent:", req.AgentId)
		ret.ErrMsg = err.Error()
		return ret, err
	}
	if err := m.sendUspMsgToAgent(req.AgentId, uspMsg, mtpIntf); err != nil {
		log.Println("Error in sending Add Instance msg to agent")
		return ret, err
	}

	time.Sleep(time.Second)

	// Read msg error data from Cache
	if cacheErr, err := m.cacheGetError(req.AgentId, req.MsgId); err == nil {
		ret.ErrMsg = cacheErr.Msg
		return ret, nil
	}
	// No hit on Msg error cache which means received add response from agent
	// Read add response data from cache now
	cacheInst, err := m.cacheGetInstance(req.AgentId, req.MsgId)
	if err != nil {
		return ret, errors.New("No response from agent")
	}
	ret.AgentId = req.AgentId
	ret.MsgId = req.MsgId
	if cacheInst.OpIsSuccess {
		grpcInst := &mtpgrpc.MtpAddInstanceResData_Instance{}
		grpcInst.Path = cacheInst.Path
		grpcInst.UniqueKeys = cacheInst.UniqueKeys
		ret.Inst = append(ret.Inst, grpcInst)
		ret.IsSuccess = true
	} else {
		ret.IsSuccess = false
		ret.ErrMsg = cacheInst.OpErrStr
	}
	return ret, nil
}

func (m *Mtp) MtpOperateReq(ctx context.Context, p *mtpgrpc.MtpOperateReqData) (*mtpgrpc.MtpOperateResData, error) {
	log.Printf("OperatgeReqTx: AgetnId: %v, MsgId: %v\n", p.AgentId, p.MsgId)
	log.Printf("OperatgeReqTx: Cmd: %v, CmdKey: %v\n", p.Cmd, p.CmdKey)
	uspMsg, err := parser.CreateUspOperateReqMsg(p.Cmd, p.CmdKey, p.Resp, p.MsgId, p.Inputs)
	ret := &mtpgrpc.MtpOperateResData{IsSuccess: false}
	if err != nil {
		log.Println("Could not prepare Usp  msg of type: OPERATE")
		errMsg := &mtpgrpc.MtpOperateResData_ErrMsg{ErrMsg: err.Error()}
		ret.Resp = errMsg
		return ret, err
	}
	mtpIntf, err := m.getAgentMtp(p.AgentId)
	if err != nil {
		log.Println("Could not get agentMtp Interface for agent:", p.AgentId)
		errMsg := &mtpgrpc.MtpOperateResData_ErrMsg{ErrMsg: err.Error()}
		ret.Resp = errMsg
		return ret, err
	}
	log.Println("OperateReq: Formed stomp msg with agendId:", p.AgentId)
	if err := m.sendUspMsgToAgent(p.AgentId, uspMsg, mtpIntf); err != nil {
		log.Println("Error in sending Operate msg to agent")
		errMsg := &mtpgrpc.MtpOperateResData_ErrMsg{ErrMsg: err.Error()}
		ret.Resp = errMsg // err.Error()
		return ret, err
	}
	if cacheErr, err := m.cacheGetError(p.AgentId, p.MsgId); err == nil {
		//ret.ErrMsg = cacheErr.Msg
		errMsg := &mtpgrpc.MtpOperateResData_ErrMsg{ErrMsg: cacheErr.Msg}
		ret.Resp = errMsg // err.Error()
		return ret, nil
	}
	ret.IsSuccess = true
	return ret, nil
}

func (m *Mtp) MtpGetDatamodelReq(ctx context.Context, p *mtpgrpc.MtpGetDatamodelReqData) (*mtpgrpc.MtpReqResult, error) {
	log.Printf("GetDatamodelReqTx: AgetnId: %v, MsgId: %v\n", p.AgentId, p.MsgId)
	log.Printf("GetDatamodelReqTx: Path: %v\n", p.Path)
	var paths []string
	paths = append(paths, p.Path)
	uspMsg, err := parser.CreateUspGetSupportedDmMsg(paths, p.RetCmd, p.RetEvents, p.RetParams, p.MsgId)
	ret := &mtpgrpc.MtpReqResult{IsSuccess: false}
	if err != nil {
		log.Println("Could not prepare Usp msg of type: GET DATAMODEL")
		return ret, err
	}
	mtpIntf, err := m.getAgentMtp(p.AgentId)
	if err != nil {
		log.Println("Could not get agentMtp Interface for agent:", p.AgentId)
		ret.ErrMsg = err.Error()
		return ret, err
	}
	log.Println("GetDatamodelReqTx: Formed stomp msg with agendId:", p.AgentId)
	if err := m.sendUspMsgToAgent(p.AgentId, uspMsg, mtpIntf); err != nil {
		log.Println("Error in sending Get Datamodel msg to agent")
		return ret, err
	}
	ret.IsSuccess = true
	return ret, nil
}

func (m *Mtp) MtpDeleteInstanceReq(ctx context.Context, p *mtpgrpc.MtpDeleteInstanceReqData) (*mtpgrpc.MtpReqResult, error) {
	log.Printf("DeleteInstanceReqTx: AgetnId: %v, MsgId: %v\n", p.AgentId, p.MsgId)
	var objPaths []string
	objPaths = append(objPaths, p.ObjPath)
	log.Printf("DeleteInstanceReqTx: Path: %v\n", p.ObjPath)
	uspMsg, err := parser.CreateUspDeleteReqMsg(objPaths, p.MsgId)
	ret := &mtpgrpc.MtpReqResult{IsSuccess: false}
	if err != nil {
		log.Println("Could not prepare Usp msg of type: DELETE INSTANCE")
		return ret, err
	}
	mtpIntf, err := m.getAgentMtp(p.AgentId)
	if err != nil {
		log.Println("Could not get agentMtp Interface for agent:", p.AgentId)
		ret.ErrMsg = err.Error()
		return ret, err
	}
	log.Println("DeleteInstanceReqTx: Formed stomp msg with agendId:", p.AgentId)
	if err := m.sendUspMsgToAgent(p.AgentId, uspMsg, mtpIntf); err != nil {
		log.Println("Error in sending delete instance msg to agent")
		return ret, err
	}
	ret.IsSuccess = true
	return ret, nil
}

func (m *Mtp) MtpGetAgentMsgs(ctx context.Context, p *mtpgrpc.MtpGetAgentMsgsData) (*mtpgrpc.MtpReqResult, error) {
	ret := &mtpgrpc.MtpReqResult{IsSuccess: false}
	if err := m.StompReceiveUspMsgFromAgentWithTimer(1); err != nil {
		log.Println("Error in receiving msg from agent, err:", err)
		return ret, err
	}
	ret.IsSuccess = true
	return ret, nil
}
