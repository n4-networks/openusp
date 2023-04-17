package rest

import (
	"context"
	"errors"
	"log"
	"strconv"

	"github.com/n4-networks/usp/pkg/pb/mtpgrpc"
	"google.golang.org/grpc"
)

var resultStr = map[bool]string{
	true:  "Passed",
	false: "Failed",
}

func (re *Rest) IsConnectedToMtp() bool {
	if re.mtp.grpcConn == nil {
		return false
	}
	return true
}

func (re *Rest) connectMtp() error {
	//log.Println("Connecting to MTP @", re.cfg.mtpAddr)

	if re.mtp.grpcConn != nil {
		re.mtp.grpcConn.Close()
	}

	ctx, cancel := context.WithTimeout(context.Background(), re.cfg.connTimeout)
	defer cancel()
	conn, err := grpc.DialContext(ctx, re.cfg.mtpAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return err
	}
	re.mtp.grpcIntf = mtpgrpc.NewMtpGrpcClient(conn)
	re.mtp.grpcConn = conn
	return nil
}

type MtpInfo struct {
	Version string `json:"version"`
}

func (re *Rest) getMtpInfoObj() (*MtpInfo, error) {
	return re.MtpGetInfo()
}

func (re *Rest) MtpGetInfo() (*MtpInfo, error) {
	if re.mtp.grpcIntf == nil {
		return nil, errors.New("MTP is not connected")
	}
	var none mtpgrpc.None
	res, err := re.mtp.grpcIntf.MtpGetInfo(context.Background(), &none)
	if err != nil {
		log.Println("gRPC error: ", err)
		return nil, err
	}
	info := &MtpInfo{}
	info.Version = res.GetVersion()
	return info, nil
}

func (re *Rest) MtpSetParamReq(epId string, path string, params map[string]string) error {
	if re.mtp.grpcIntf == nil {
		return errors.New("MTP is not connected")
	}
	var paramName, paramValue string
	for k, v := range params {
		paramName = k
		paramValue = v
		break
	}
	var in mtpgrpc.MtpSetParamReqData
	in.AgentId = epId
	in.MsgId = "SET" + strconv.FormatUint(re.mtp.incTxMsgCnt(), 10)
	in.Path = path
	in.Param = paramName
	in.Value = paramValue
	log.Println("Sending setparam request to MTP, path:", path)
	log.Printf("%v: %v\n", paramName, paramValue)
	res, err := re.mtp.grpcIntf.MtpSetParamReq(context.Background(), &in)
	if err != nil {
		log.Println("gRPC error: ", err)
	}
	if !res.GetIsSuccess() {
		log.Println("Error in executing MtpSetParamReq")
		return errors.New(res.GetErrMsg())
	}
	return err
}

func (re *Rest) MtpGetParamReq(epId string, path string) error {
	if re.mtp.grpcIntf == nil {
		return errors.New("MTP is not connected")
	}
	var in mtpgrpc.MtpGetParamReqData
	in.AgentId = epId
	in.MsgId = "GET" + strconv.FormatUint(re.mtp.incTxMsgCnt(), 10)
	in.Path = path
	log.Println("Sending getparam request to MTP, path:", path)
	out, err := re.mtp.grpcIntf.MtpGetParamReq(context.Background(), &in)
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

func (re *Rest) MtpGetInstancesReq(epId string, objPath string, firstLevelOnly bool) error {
	if re.mtp.grpcIntf == nil {
		return errors.New("MTP is not connected")
	}
	var in mtpgrpc.MtpGetInstancesReqData
	in.AgentId = epId
	in.MsgId = "GET_INST" + strconv.FormatUint(re.mtp.incTxMsgCnt(), 10)
	in.Path = objPath
	in.FirstLevelOnly = firstLevelOnly
	log.Println("Sending GetInstance request to MTP, path:", objPath)
	out, err := re.mtp.grpcIntf.MtpGetInstancesReq(context.Background(), &in)
	if err != nil {
		log.Println("gRPC error: ", err)
	}
	if !out.GetIsSuccess() {
		log.Printf("Error: %v", out.GetErrMsg())
		return errors.New(out.GetErrMsg())
	}
	return nil
}

func (re *Rest) MtpAddInstanceReq(epId string, objs []*object) ([]*Instance, error) {
	if re.mtp.grpcIntf == nil {
		return nil, errors.New("MTP is not connected")
	}
	var in mtpgrpc.MtpAddInstanceReqData

	in.AgentId = epId
	in.MsgId = "ADD" + strconv.FormatUint(re.mtp.incTxMsgCnt(), 10)
	for _, obj := range objs {
		addInstanceObj := &mtpgrpc.MtpAddInstanceReqData_Object{}
		addInstanceObj.Path = obj.path
		addInstanceObj.Params = obj.params
		in.Objs = append(in.Objs, addInstanceObj)
	}

	var instances []*Instance
	log.Println("Sending AddInstance request to MTP")
	out, err := re.mtp.grpcIntf.MtpAddInstanceReq(context.Background(), &in)
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

func (re *Rest) MtpOperateReq(epId string, cmd string, cmdKey string, resp bool, inputs map[string]string) error {
	if re.mtp.grpcIntf == nil {
		return errors.New("MTP is not connected")
	}
	var in mtpgrpc.MtpOperateReqData
	in.AgentId = epId
	in.MsgId = "OP" + strconv.FormatUint(re.mtp.incTxMsgCnt(), 10)
	in.Cmd = cmd
	in.CmdKey = cmdKey
	in.Resp = resp
	in.Inputs = inputs
	log.Println("Sending Operate request to MTP, cmd:", cmd)
	out, err := re.mtp.grpcIntf.MtpOperateReq(context.Background(), &in)
	if err != nil {
		log.Println("gRPC error: ", err)
	}
	if !out.GetIsSuccess() {
		log.Printf("Error: %v", out.GetErrMsg())
		return errors.New(out.GetErrMsg())
	}
	return nil
}

func (re *Rest) MtpGetDatamodelReq(epId string, path string) error {
	if re.mtp.grpcIntf == nil {
		return errors.New("MTP is not connected")
	}
	var in mtpgrpc.MtpGetDatamodelReqData
	in.AgentId = epId
	in.MsgId = "GET_DM" + strconv.FormatUint(re.mtp.incTxMsgCnt(), 10)
	in.Path = path
	in.RetCmd = true
	in.RetEvents = true
	in.RetParams = true
	log.Println("Sending Get Datamodel request to MTP, path:", path)
	out, err := re.mtp.grpcIntf.MtpGetDatamodelReq(context.Background(), &in)
	if err != nil {
		log.Println("gRPC error: ", err)
	}
	if !out.GetIsSuccess() {
		log.Printf("Error: %v", out.GetErrMsg())
		return errors.New(out.GetErrMsg())
	}
	return nil
}

func (re *Rest) MtpDeleteInstanceReq(epId string, objPath string) error {
	if re.mtp.grpcIntf == nil {
		return errors.New("MTP is not connected")
	}

	var in mtpgrpc.MtpDeleteInstanceReqData
	in.AgentId = epId
	in.MsgId = "DELETE_" + strconv.FormatUint(re.mtp.incTxMsgCnt(), 10)
	in.ObjPath = objPath
	log.Println("Sending Delete Instance request to MTP")
	out, err := re.mtp.grpcIntf.MtpDeleteInstanceReq(context.Background(), &in)
	if err != nil {
		log.Println("gRPC error: ", err)
	}
	if !out.GetIsSuccess() {
		log.Printf("Error: %v", out.GetErrMsg())
		return errors.New(out.GetErrMsg())
	}
	return nil
}

func (re *Rest) MtpGetAgentMsgs(epId string) error {
	if re.mtp.grpcIntf == nil {
		return errors.New("MTP is not connected")
	}
	var in mtpgrpc.MtpGetAgentMsgsData
	in.AgentId = epId
	log.Println("Sending get agent msg request to MTP")
	out, err := re.mtp.grpcIntf.MtpGetAgentMsgs(context.Background(), &in)
	if err != nil {
		log.Println("gRPC error: ", err)
	}
	if !out.GetIsSuccess() {
		log.Printf("Error: %v", out.GetErrMsg())
		return errors.New(out.GetErrMsg())
	}
	return nil
}
