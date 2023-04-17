package mtp

import (
	"errors"
	"log"

	"github.com/n4-networks/usp/pkg/db"
	"github.com/n4-networks/usp/pkg/pb/bbf/usp_msg"
)

type uspMsgErr struct {
	code uint32
	msg  string
}

type uspMsgData struct {
	id        string
	mType     usp_msg.Header_MsgType
	err       *uspMsgErr
	params    []*param
	protoVer  string
	instances []*instance
	paths     []string
	notify    *notification
	dms       []*db.DmObject
}

var (
	DmObjAccessString = map[int32]string{
		0: "OBJ_READ_ONLY",
		1: "OBJ_ADD_DELETE",
		2: "OBJ_ADD_ONLY",
		3: "OBJ_DELETE_ONLY",
	}
	DmObjAccessInt = map[string]int32{
		"OBJ_READ_ONLY":   0,
		"OBJ_ADD_DELETE":  1,
		"OBJ_ADD_ONLY":    2,
		"OBJ_DELETE_ONLY": 3,
	}
	DmParamAccessString = map[int32]string{
		0: "PARAM_READ_ONLY",
		1: "PARAM_READ_WRITE",
		2: "PARAM_WRITE_ONLY",
	}
	DmParamAccessInt = map[string]int32{
		"PARAM_READ_ONLY":  0,
		"PARAM_READ_WRITE": 1,
		"PARAM_WRITE_ONLY": 2,
	}
)

func (m *Mtp) processRxUspMsg(epId string, mData *uspMsgData) error {

	if mData == nil {
		log.Println("Can not proceed, msgData is nil")
		return errors.New("msgData has not been initialized")
	}

	switch mData.mType {
	case usp_msg.Header_ERROR:
		if mData.err != nil {
			cErr := &CError{}
			cErr.Msg = mData.err.msg
			cErr.Code = mData.err.code
			m.cacheSetError(epId, mData.id, cErr)
		} else {
			log.Println("mData.err is not initialized")
			return errors.New("mData.err not initialized")
		}
	case usp_msg.Header_GET_RESP:
		if mData.err != nil {
			cErr := &CError{}
			cErr.Msg = mData.err.msg
			cErr.Code = mData.err.code
			m.cacheSetError(epId, mData.id, cErr)
		}
		log.Println("Received GET_RESP msg, writing to DB")
		if len(mData.params) > 0 {
			if err := m.dbWriteParams(epId, mData.params); err != nil {
				return err
			}
		}
	case usp_msg.Header_GET_INSTANCES_RESP:
		//log.Println("Received New Instances, writing to DB")
		if mData.err != nil {
			cErr := &CError{}
			cErr.Msg = mData.err.msg
			cErr.Code = mData.err.code
			m.cacheSetError(epId, mData.id, cErr)
		}
		if len(mData.instances) > 0 {
			if err := m.dbWriteInstances(epId, mData.instances); err != nil {
				return err
			}
		}
	case usp_msg.Header_ADD_RESP:
		cInst := &CInstance{}
		if mData.err != nil {
			cInst.OpIsSuccess = false
			cInst.OpErrStr = mData.err.msg
		} else {
			cInst.OpIsSuccess = true
			cInst.Path = mData.instances[0].path
			cInst.UniqueKeys = mData.instances[0].uniqueKeys
		}
		m.cacheSetInstance(epId, mData.id, cInst)
		if len(mData.instances) > 0 {
			if err := m.dbWriteInstances(epId, mData.instances); err != nil {
				return err
			}
		}
	case usp_msg.Header_SET_RESP:
		cSetResult := &CParamSetResult{}
		if mData.err != nil {
			cSetResult.OpIsSuccess = false
			cSetResult.OpErrStr = mData.err.msg
		} else {
			log.Println("Set Response result: Success")
			cSetResult.OpIsSuccess = true
			cSetResult.Path = mData.params[0].path
		}
		m.cacheSetParamSetResult(epId, mData.id, cSetResult)
		if len(mData.instances) > 0 {
			if err := m.dbWriteInstances(epId, mData.instances); err != nil {
				return err
			}
		}
	case usp_msg.Header_GET_SUPPORTED_DM_RESP:
		if len(mData.dms) > 0 {
			if err := m.dbWriteDatamodels(mData.dms); err != nil {
				return err
			}
		}
	case usp_msg.Header_DELETE_RESP:
		if mData.err != nil {
			cErr := &CError{}
			cErr.Msg = mData.err.msg
			m.cacheSetError(epId, mData.id, cErr)
		}
		if len(mData.paths) > 0 {
			if err := m.dbDeleteInstances(epId, mData.paths); err != nil {
				return err
			}
			if err := m.dbDeleteParams(epId, mData.paths); err != nil {
				return err
			}
		}
		// Notify response is being handled by MtpReceiveThread (e.g. StompReceive or coapReceive)
		//case usp_msg.Header_NOTIFY:
		//m.handleNotifyReq(epId, mData.id, mData.notify, agentMtp)
	}
	return nil
}

func parseUspMsg(rData *uspRecordData) (*uspMsgData, error) {
	mData := &uspMsgData{}
	h := rData.msg.GetHeader()

	mData.id = h.GetMsgId()
	log.Println("Rx USP MsgId: ", mData.id)

	mData.mType = h.GetMsgType()
	log.Println("Rx USP MsgType: ", mData.mType)

	resp := rData.msg.GetBody().GetResponse()

	switch mData.mType {
	case usp_msg.Header_ERROR:
		log.Println("USP Msg Type: Error")
		mData.err = &uspMsgErr{}
		mData.err.code = rData.msg.GetBody().GetError().GetErrCode()
		mData.err.msg = rData.msg.GetBody().GetError().GetErrMsg()
		log.Println("Error code:", mData.err.code)
		log.Println("Error msg:", mData.err.msg)

	case usp_msg.Header_GET_RESP:
		log.Println("USP Msg Type: GET Response")
		params, mErr, _ := processGetResp(resp.GetGetResp())
		mData.params = params
		if mErr != nil {
			mData.err = mErr
		}

	case usp_msg.Header_GET_SUPPORTED_PROTO_RESP:
		log.Println("USP Msg Type: SUPPORTED PROTO Response")
		mData.protoVer = resp.GetGetSupportedProtocolResp().GetAgentSupportedProtocolVersions()
		log.Println("Agent Supported Proto version: ", mData.protoVer)

	case usp_msg.Header_OPERATE_RESP:
		log.Println("USP Msg Type: OPERATE Response")
		processOperateResp(resp.GetOperateResp())

	case usp_msg.Header_GET_INSTANCES_RESP:
		log.Println("USP Msg Type: GET INSTANCES Response")
		instances, _ := processGetInstancesResp(resp.GetGetInstancesResp())
		mData.instances = instances

	case usp_msg.Header_ADD_RESP:
		log.Println("USP Msg Type: ADD Response")
		instances, mErr, _ := processAddResp(resp.GetAddResp())
		mData.instances = instances
		mData.err = mErr

	case usp_msg.Header_SET_RESP:
		log.Println("USP Msg Type: SET Response")
		params, mErr, _ := processSetResp(resp.GetSetResp())
		mData.params = params
		mData.err = mErr

	case usp_msg.Header_DELETE_RESP:
		log.Println("USP Msg Type: DELETE Response")
		paths, mErr, _ := processDeleteResp(resp.GetDeleteResp())
		mData.paths = paths
		mData.err = mErr

	case usp_msg.Header_NOTIFY:
		log.Println("USP Msg Type: NOTIFY Request from Agent")
		notify, _ := processNotify(rData.msg.GetBody().GetRequest())
		mData.notify = notify

	case usp_msg.Header_NOTIFY_RESP:
		log.Println("USP Msg Type: NOTIFY Response")

	case usp_msg.Header_GET_SUPPORTED_DM_RESP:
		log.Println("SUPPORTED_DM_RESP")
		dmObjs, mErr, err := processGetSupportedDmResp(resp.GetGetSupportedDmResp(), rData.fromId)
		if mErr != nil {
			mData.err = mErr
			return nil, err
		}
		mData.dms = dmObjs
		//if err := m.dbWriteDatamodels(dmObjs); err != nil {
	//		return nil, err
	//	}

	default:
		log.Fatalln("Invalid Msg Type in Incoming USP Msg")
		return nil, errors.New("Wrong USP msg type")
	}
	return mData, nil
}

func processNotify(r *usp_msg.Request) (*notification, error) {

	notify := &notification{}

	if n := r.GetNotify(); n != nil {
		notify.subscriptionId = n.GetSubscriptionId()
		notify.sendResp = n.GetSendResp()
		if eMsg := n.GetEvent(); eMsg != nil {
			notify.nType = NotifyEvent

			evt := &event{}
			evt.name = eMsg.GetEventName()
			log.Println("Received an event of:", evt.name)
			evt.path = eMsg.GetObjPath()
			evt.params = eMsg.GetParams()
			notify.evt = evt
			return notify, nil
		}
		if valCMsg := n.GetValueChange(); valCMsg != nil {
			notify.nType = NotifyValueChange

			valC := &valueChange{}
			valC.paramPath = valCMsg.GetParamPath()
			log.Println("Received value change Notification for:", valC.paramPath)
			valC.paramValue = valCMsg.GetParamValue()
			notify.valChange = valC
			return notify, nil
		}
		if objCMsg := n.GetObjCreation(); objCMsg != nil {
			notify.nType = NotifyObjCreation

			objC := &objectCreation{}
			objC.path = objCMsg.GetObjPath()
			log.Println("Received Obj creation Notification for:", objC.path)
			objC.uniqueKeys = objCMsg.GetUniqueKeys()
			notify.objCreation = objC
			return notify, nil
		}
		if objDMsg := n.GetObjDeletion(); objDMsg != nil {
			notify.nType = NotifyObjDeletion

			objD := &objectDeletion{}
			objD.path = objDMsg.GetObjPath()
			log.Println("Received Obj deletion Notification for:", objD.path)

			notify.objDeletion = objD
			return notify, nil
		}
		if opCMsg := n.GetOperComplete(); opCMsg != nil {
			notify.nType = NotifyOpComplete

			opC := &operationComplete{}
			opC.path = opCMsg.GetObjPath()
			opC.cmdName = opCMsg.GetCommandName()
			opC.cmdKey = opCMsg.GetCommandKey()
			if args := opCMsg.GetReqOutputArgs(); args != nil {
				opC.outArg = args.GetOutputArgs()
			} else if cmdF := opCMsg.GetCmdFailure(); cmdF != nil {
				opC.cmdFailure.errCode = cmdF.GetErrCode()
				opC.cmdFailure.errMsg = cmdF.GetErrMsg()
			}
			notify.opComplete = opC
			return notify, nil
		}
		if obMsg := n.GetOnBoardReq(); obMsg != nil {
			notify.nType = NotifyOnBoardReq

			ob := &onBoardReq{}
			ob.oui = obMsg.GetOui()
			ob.productClass = obMsg.GetProductClass()
			ob.serialNum = obMsg.GetSerialNumber()
			ob.protoVer = obMsg.GetAgentSupportedProtocolVersions()
			notify.onBoard = ob
			return notify, nil
		}

	}
	return nil, errors.New("Invalid NotificationType")
}

func processGetInstancesResp(r *usp_msg.GetInstancesResp) ([]*instance, error) {
	var instances []*instance
	for _, res := range r.GetReqPathResults() {
		log.Println("Requested Path: ", res.GetRequestedPath())
		if e := res.GetErrCode(); e != 0 {
			log.Println("Error Code: ", e)
			eMsg := res.GetErrMsg()
			log.Println("Error Msg: ", eMsg)
			return nil, errors.New(eMsg)
		}
		for _, i := range res.GetCurrInsts() {
			inst := &instance{}
			inst.path = i.GetInstantiatedObjPath()
			//log.Println("Instantiated Obj Path: ", inst.path)
			uKeys := i.GetUniqueKeys()
			inst.uniqueKeys = make(map[string]string, len(uKeys))
			for key, val := range uKeys {
				//log.Printf("Unique key: %v, Value: %v\n", key, val)
				inst.uniqueKeys[key] = val
			}
			instances = append(instances, inst)
		}
	}
	return instances, nil
}

func processAddResp(r *usp_msg.AddResp) ([]*instance, *uspMsgErr, error) {
	var instances []*instance
	mErr := &uspMsgErr{}
	objResults := r.GetCreatedObjResults()
	for _, res := range objResults {
		path := res.GetRequestedPath()
		log.Println("Requested Path: ", path)
		if f := res.GetOperStatus().GetOperFailure(); f != nil {
			mErr.code = f.GetErrCode()
			mErr.msg = f.GetErrMsg()
			log.Println("Error Code: ", mErr.code)
			log.Println("Error Msg: ", mErr.msg)
			return nil, mErr, errors.New("Usp msg Error")
			//continue
		}
		if s := res.GetOperStatus().GetOperSuccess(); s != nil {
			inst := &instance{}
			inst.path = s.GetInstantiatedPath()
			log.Println("Object created at:", inst.path)
			inst.uniqueKeys = s.GetUniqueKeys()
			instances = append(instances, inst)
		}
	}
	return instances, nil, nil
}

func processGetSupportedDmResp(r *usp_msg.GetSupportedDMResp, epId string) ([]*db.DmObject, *uspMsgErr, error) {
	var dmObjs []*db.DmObject
	mErr := &uspMsgErr{}
	for _, res := range r.GetReqObjResults() {
		log.Println("ObjectPath: ", res.GetReqObjPath())
		if eCode := res.GetErrCode(); eCode != 0 {
			mErr.code = eCode
			mErr.msg = res.GetErrMsg()
			log.Println("Error Code: ", mErr.code)
			log.Println("Error Msg: ", mErr.msg)
			return nil, mErr, errors.New("USP msg error")
		}
		log.Println("DataModelInstUri: ", res.GetDataModelInstUri())

		for _, obj := range res.GetSupportedObjs() {
			dmObj := &db.DmObject{}
			dmObj.Path = obj.GetSupportedObjPath()
			dmObj.MultiInstance = obj.GetIsMultiInstance()
			dmObj.Access = DmObjAccessString[int32(obj.GetAccess())]
			cmds := obj.GetSupportedCommands()
			dmObj.Cmds = make([]db.DmCommand, len(cmds))
			for j, cmd := range cmds {
				dmObj.Cmds[j].Name = cmd.GetCommandName()
				log.Println("Command Name: ", cmd.GetCommandName())
				inputArgs := cmd.GetInputArgNames()
				dmObj.Cmds[j].Inputs = make([]string, len(inputArgs))
				for k, arg := range inputArgs {
					dmObj.Cmds[j].Inputs[k] = arg
					log.Println("Input Arg Name: ", arg)
				}
				outputArgs := cmd.GetOutputArgNames()
				dmObj.Cmds[j].Outputs = make([]string, len(outputArgs))
				for k, arg := range outputArgs {
					dmObj.Cmds[j].Outputs[k] = arg
					log.Println("Output Arg Name: ", arg)
				}
			}
			evts := obj.GetSupportedEvents()
			dmObj.Events = make([]db.DmEvent, len(evts))
			//log.Println("Number of events:", len(evts))
			for j, evt := range evts {
				log.Println("Event Name: ", evt.GetEventName())
				dmObj.Events[j].Name = evt.GetEventName()
				args := evt.GetArgNames()
				dmObj.Events[j].Args = make([]string, len(args))
				for k, name := range args {
					log.Println("Arg Name: ", name)
					dmObj.Events[j].Args[k] = name
				}
			}
			params := obj.GetSupportedParams()
			dmObj.Params = make([]db.DmParam, len(params))
			for j, param := range params {
				//log.Println("Param Name: ", param.GetParamName())
				dmObj.Params[j].Name = param.GetParamName()
				dmObj.Params[j].Access = DmParamAccessString[int32(param.GetAccess())]
			}
			dmObj.EndpointId = epId
			dmObjs = append(dmObjs, dmObj)
		}
	}
	return dmObjs, nil, nil
}

func processOperateResp(r *usp_msg.OperateResp) error {
	results := r.GetOperationResults()
	for _, opRes := range results {
		log.Println("Executed Command: ", opRes.GetExecutedCommand())
		for key, val := range opRes.GetReqOutputArgs().GetOutputArgs() {
			log.Printf("Output Args: Key: %v, Val: %v\n", key, val)
		}
	}
	return nil
}

func processDeleteResp(r *usp_msg.DeleteResp) ([]string, *uspMsgErr, error) {
	var affectedPaths []string
	mErr := &uspMsgErr{}

	objResults := r.GetDeletedObjResults()
	for _, res := range objResults {
		log.Println("Requested Path: ", res.GetRequestedPath())
		if f := res.GetOperStatus().GetOperFailure(); f != nil {
			mErr.code = f.GetErrCode()
			mErr.msg = f.GetErrMsg()
			log.Printf("Delete error: %v, msg: %v\n", mErr.code, mErr.msg)
			return nil, mErr, errors.New("USP msg error")
		}
		affectedPaths = res.GetOperStatus().GetOperSuccess().GetAffectedPaths()
	}
	return affectedPaths, nil, nil
}

func processSetResp(r *usp_msg.SetResp) ([]*param, *uspMsgErr, error) {
	var params []*param
	mErr := &uspMsgErr{}
	objResults := r.GetUpdatedObjResults()
	for _, res := range objResults {
		log.Println("Requested Path:", res.GetRequestedPath())
		if e := res.GetOperStatus().GetOperFailure(); e != nil {
			mErr.code = e.GetErrCode()
			mErr.msg = e.GetErrMsg()
			log.Println("Set operaion failed with ErrCode:", mErr.code)
			log.Println("ErrMsg:", mErr.msg)
			return nil, mErr, errors.New("Usp msg err")
		}
		if s := res.GetOperStatus().GetOperSuccess(); s != nil {
			instResults := s.GetUpdatedInstResults()
			for _, inst := range instResults {
				path := inst.GetAffectedPath()
				log.Println("Affected path:", path)
				paramMap := inst.GetUpdatedParams()
				for key, val := range paramMap {
					param := &param{}
					param.path = path + key
					param.value = val
					params = append(params, param)
				}
			}
		}
	}
	return params, nil, nil
}

func processGetResp(r *usp_msg.GetResp) ([]*param, *uspMsgErr, error) {
	var params []*param
	mErr := &uspMsgErr{}
	for _, reqPath := range r.GetReqPathResults() {
		log.Println("Requested Path: ", reqPath.GetRequestedPath())
		if eCode := reqPath.GetErrCode(); eCode != 0 {
			log.Println("Error Code: ", eCode)
			mErr.code = eCode
			mErr.msg = reqPath.GetErrMsg()
			log.Println("Error Msg: ", mErr.msg)
			return nil, mErr, errors.New("Usp Msg error")
		}
		resPathResults := reqPath.GetResolvedPathResults()
		for _, resPath := range resPathResults {
			path := resPath.GetResolvedPath()
			log.Println("Resolved Path: ", path)
			paramMap := resPath.GetResultParams()
			//log.Printf("Result Params:%+v\n", paramMap)
			for key, val := range paramMap {
				param := &param{}
				param.path = path + key
				param.value = val
				params = append(params, param)
			}
		}
	}
	return params, nil, nil
}
