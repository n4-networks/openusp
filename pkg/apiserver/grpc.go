package apiserver

import (
	"context"
	"errors"
	"log"
	"strconv"

	"github.com/n4-networks/openusp/pkg/pb/cntlrgrpc"
	"google.golang.org/grpc"
)

var resultStr = map[bool]string{
	true:  "Passed",
	false: "Failed",
}

func (as *ApiServer) IsConnectedToCntlr() bool {
	if as.grpcH.conn == nil {
		return false
	}
	return true
}

func (as *ApiServer) connectToController() error {
	//log.Println("Connecting to Controller @", as.cfg.cntlrAddr)

	if as.grpcH.conn != nil {
		as.grpcH.conn.Close()
	}

	ctx, cancel := context.WithTimeout(context.Background(), as.cfg.connTimeout)
	defer cancel()
	conn, err := grpc.DialContext(ctx, as.cfg.cntlrAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return err
	}
	as.grpcH.intf = cntlrgrpc.NewGrpcClient(conn)
	as.grpcH.conn = conn
	return nil
}

type CntlrInfo struct {
	Version string `json:"version"`
}

func (as *ApiServer) getCntlrInfoObj() (*CntlrInfo, error) {
	return as.GetCntlrInfo()
}

func (as *ApiServer) GetCntlrInfo() (*CntlrInfo, error) {
	if as.grpcH.intf == nil {
		return nil, errors.New("Controller is not connected")
	}
	var none cntlrgrpc.None
	res, err := as.grpcH.intf.GetInfo(context.Background(), &none)
	if err != nil {
		log.Println("gRPC error: ", err)
		return nil, err
	}
	info := &CntlrInfo{}
	info.Version = res.GetVersion()
	return info, nil
}

func (as *ApiServer) CntlrSetParamReq(epId string, path string, params map[string]string) error {
	if as.grpcH.intf == nil {
		return errors.New("Controller is not connected")
	}
	var paramName, paramValue string
	for k, v := range params {
		paramName = k
		paramValue = v
		break
	}
	var in cntlrgrpc.SetParamReqData
	in.AgentId = epId
	in.MsgId = "SET" + strconv.FormatUint(as.grpcH.incTxMsgCnt(), 10)
	in.Path = path
	in.Param = paramName
	in.Value = paramValue
	log.Println("Sending setparam request to Controller, path:", path)
	log.Printf("%v: %v\n", paramName, paramValue)
	res, err := as.grpcH.intf.SetParamReq(context.Background(), &in)
	if err != nil {
		log.Println("gRPC error: ", err)
	}
	if !res.GetIsSuccess() {
		log.Println("Error in executing CntlrSetParamReq")
		return errors.New(res.GetErrMsg())
	}
	return err
}

func (as *ApiServer) CntlrGetParamReq(epId string, path string) error {
	if as.grpcH.intf == nil {
		return errors.New("Controller is not connected")
	}
	var in cntlrgrpc.GetParamReqData
	in.AgentId = epId
	in.MsgId = "GET" + strconv.FormatUint(as.grpcH.incTxMsgCnt(), 10)
	in.Path = path
	log.Println("Sending getparam request to Controller , path:", path)
	out, err := as.grpcH.intf.GetParamReq(context.Background(), &in)
	if err != nil {
		log.Println("gRPC error: ", err)
	}
	if !out.GetIsSuccess() {
		errMsg := out.GetErrMsg()
		log.Printf("Result: %v\n", errMsg)
		return errors.New(errMsg)
	}
	return nil
}

func (as *ApiServer) CntlrGetInstancesReq(epId string, objPath string, firstLevelOnly bool) error {
	if as.grpcH.intf == nil {
		return errors.New("Controller is not connected")
	}
	var in cntlrgrpc.GetInstancesReqData
	in.AgentId = epId
	in.MsgId = "GET_INST" + strconv.FormatUint(as.grpcH.incTxMsgCnt(), 10)
	in.Path = objPath
	in.FirstLevelOnly = firstLevelOnly
	log.Println("Sending GetInstance request to Controller, path:", objPath)
	out, err := as.grpcH.intf.GetInstancesReq(context.Background(), &in)
	if err != nil {
		log.Println("gRPC error: ", err)
	}
	if !out.GetIsSuccess() {
		log.Printf("Error: %v", out.GetErrMsg())
		return errors.New(out.GetErrMsg())
	}
	return nil
}

func (as *ApiServer) CntlrAddInstanceReq(epId string, objs []*object) ([]*Instance, error) {
	if as.grpcH.intf == nil {
		return nil, errors.New("Controller is not connected")
	}
	var in cntlrgrpc.AddInstanceReqData

	in.AgentId = epId
	in.MsgId = "ADD" + strconv.FormatUint(as.grpcH.incTxMsgCnt(), 10)
	for _, obj := range objs {
		addInstanceObj := &cntlrgrpc.AddInstanceReqData_Object{}
		addInstanceObj.Path = obj.path
		addInstanceObj.Params = obj.params
		in.Objs = append(in.Objs, addInstanceObj)
	}

	var instances []*Instance
	log.Println("Sending AddInstance request to Controller")
	out, err := as.grpcH.intf.AddInstanceReq(context.Background(), &in)
	if err != nil {
		log.Println("gRPC error: ", err)
		return nil, err
	} else {
		log.Println("Result:", resultStr[out.GetIsSuccess()])
		if out.GetIsSuccess() {
			for _, grpcInst := range out.GetInst() {
				inst := &Instance{}
				inst.Path = grpcInst.GetPath()
				inst.UniqueKeys = grpcInst.GetUniqueKeys()
				instances = append(instances, inst)
			}
		} else {
			log.Println("Error from Agent:", out.GetErrMsg())
			return nil, errors.New(out.GetErrMsg())
		}
	}
	return instances, nil
}

func (as *ApiServer) CntlrOperateReq(epId string, cmd string, cmdKey string, resp bool, inputs map[string]string) error {
	if as.grpcH.intf == nil {
		return errors.New("Controller is not connected")
	}
	var in cntlrgrpc.OperateReqData
	in.AgentId = epId
	in.MsgId = "OP" + strconv.FormatUint(as.grpcH.incTxMsgCnt(), 10)
	in.Cmd = cmd
	in.CmdKey = cmdKey
	in.Resp = resp
	in.Inputs = inputs
	log.Println("Sending Operate request to Controller, cmd:", cmd)
	out, err := as.grpcH.intf.OperateReq(context.Background(), &in)
	if err != nil {
		log.Println("gRPC error: ", err)
	}
	if !out.GetIsSuccess() {
		log.Printf("Error: %v", out.GetErrMsg())
		return errors.New(out.GetErrMsg())
	}
	return nil
}

func (as *ApiServer) CntlrGetDatamodelReq(epId string, path string) error {
	if as.grpcH.intf == nil {
		return errors.New("Controller is not connected")
	}
	var in cntlrgrpc.GetDatamodelReqData
	in.AgentId = epId
	in.MsgId = "GET_DM" + strconv.FormatUint(as.grpcH.incTxMsgCnt(), 10)
	in.Path = path
	in.RetCmd = true
	in.RetEvents = true
	in.RetParams = true
	log.Println("Sending Get Datamodel request to Controller, path:", path)
	out, err := as.grpcH.intf.GetDatamodelReq(context.Background(), &in)
	if err != nil {
		log.Println("gRPC error: ", err)
	}
	if !out.GetIsSuccess() {
		log.Printf("Error: %v", out.GetErrMsg())
		return errors.New(out.GetErrMsg())
	}
	return nil
}

func (as *ApiServer) CntlrDeleteInstanceReq(epId string, objPath string) error {
	if as.grpcH.intf == nil {
		return errors.New("Controller is not connected")
	}

	var in cntlrgrpc.DeleteInstanceReqData
	in.AgentId = epId
	in.MsgId = "DELETE_" + strconv.FormatUint(as.grpcH.incTxMsgCnt(), 10)
	in.ObjPath = objPath
	log.Println("Sending Delete Instance request to Controller")
	out, err := as.grpcH.intf.DeleteInstanceReq(context.Background(), &in)
	if err != nil {
		log.Println("gRPC error: ", err)
	}
	if !out.GetIsSuccess() {
		log.Printf("Error: %v", out.GetErrMsg())
		return errors.New(out.GetErrMsg())
	}
	return nil
}

func (as *ApiServer) CntlrGetAgentMsgs(epId string) error {
	if as.grpcH.intf == nil {
		return errors.New("Controller is not connected")
	}
	var in cntlrgrpc.GetAgentMsgsData
	in.AgentId = epId
	log.Println("Sending get agent msg request to Controller")
	out, err := as.grpcH.intf.GetAgentMsgs(context.Background(), &in)
	if err != nil {
		log.Println("gRPC error: ", err)
	}
	if !out.GetIsSuccess() {
		log.Printf("Error: %v", out.GetErrMsg())
		return errors.New(out.GetErrMsg())
	}
	return nil
}
