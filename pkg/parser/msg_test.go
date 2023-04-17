package parser

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/n4-networks/usp/pkg/pb/bbf/usp_msg"
)

func usage(args []string) {

	fmt.Printf("usage:\nTo launch Server\n\t -s -p <portNumber>\nTo launch Client\n\t -c -h <Server ip> -p <ServerPort> \n")
	fmt.Printf("Arguments passed were \n\t")
	fmt.Println(args)
}

var client *websocket.Conn = nil
var clientKill chan bool
var launchChan chan bool

func TestLaunchServer(t *testing.T) {

	ProcessUspMsg = ProcessAgentMsgs // Initialize testing of Agent Masgs

	if launchChan == nil {
		launchChan = make(chan bool)
	}

	go wsServer(launchChan, "8080")

	if <-launchChan == false {
		t.Errorf("Could not launch")
		return
	}

	if clientKill == nil {
		clientKill = make(chan bool)
	}

	success := make(chan bool)

	go wsClient("127.0.0.1", "8080", success, clientKill)

	t.Log("Waiting for client to connect")
	if <-success == false {
		t.Errorf("Error connecting websocket server")
	}
	t.Log("Connected to server via websocket, can continue testing")
}

var id int = 1

func sendCreatedMsg(msgStream []byte) error {

	if client != nil {

		log.Printf("Send USP Msg Via websocket")
		var sig string = "My Test Signature"
		var cert string = "My Test Certificate"
		var to string = "Toid"
		var from string = "FromId"
		//recStream, err := usprecords.CreateNewPlainTextRecord(&to, &from, []byte(sig), []byte(cert), msgStream)
		recStream, err := CreateNewPlainTextRecord(&to, &from, []byte(sig), []byte(cert), msgStream)

		if err != nil {

			return errors.New("Could not create USP record from the msg stream")
		}

		err = client.WriteMessage(websocket.BinaryMessage, recStream)

		if err == nil {

			return nil
		}
		return errors.New("Could not send operate msg to Agent")
	}
	return errors.New("unknown USP Send error, client not connected")
}

// TestcreateUspOperateReqMsg ...
func TestUspOperateReqMsg(t *testing.T) {

	log.Printf("Create Operate Req")

	byteStream, err := CreateUspOperateReqMsg("Device.Reboot()", "1234", true, strconv.Itoa(id), nil)
	if err == nil {
		err = sendCreatedMsg(byteStream)
		if err == nil {
			return
		}
		t.Errorf("Send Operate msg failed")
		return
	}
	t.Error("Create Operate Msg failed")
}

func TestCreateUspGetReqMsg(t *testing.T) {

	var params []string = make([]string, 5)
	params[0] = "Device."
	params[1] = "Device."
	params[2] = "Device."
	params[3] = "Device."
	params[4] = "Device."

	id++

	byteStream, err := CreateUspGetReqMsg(params, strconv.Itoa(id))

	if err == nil {
		err = sendCreatedMsg(byteStream)
		if err == nil {
			return
		}
		t.Errorf(err.Error())
		return
	}

	t.Errorf(err.Error())
}

// createUspSetReqMsg
func TestCreateUspSetReqMsg(t *testing.T) {
	log.Printf("Create Operate Req")

	id++

	var setObj []*usp_msg.Set_UpdateObject = make([]*usp_msg.Set_UpdateObject, 3)
	var ob1 usp_msg.Set_UpdateObject
	var pr1 usp_msg.Set_UpdateParamSetting
	var pr11 usp_msg.Set_UpdateParamSetting
	var ob2 usp_msg.Set_UpdateObject
	var pr2 usp_msg.Set_UpdateParamSetting
	var pr22 usp_msg.Set_UpdateParamSetting
	var ob3 usp_msg.Set_UpdateObject
	var pr3 usp_msg.Set_UpdateParamSetting
	var pr33 usp_msg.Set_UpdateParamSetting

	pr1.Param = "test1"
	pr1.Required = true
	pr1.Value = "vTest1"

	pr11.Param = "test11"
	pr11.Required = true
	pr11.Value = "vTest11"

	pr2.Param = "test2"
	pr2.Required = true
	pr2.Value = "vTest2"

	pr22.Param = "test22"
	pr22.Required = true
	pr22.Value = "vTest22"

	pr3.Param = "test3"
	pr3.Required = true
	pr3.Value = "vTest3"

	pr33.Param = "test33"
	pr33.Required = true
	pr33.Value = "vTest33"

	ob1.ObjPath = "TestObj1"
	ob1.ParamSettings = make([]*usp_msg.Set_UpdateParamSetting, 2)
	ob1.ParamSettings[0] = &pr1
	ob1.ParamSettings[1] = &pr11

	ob2.ObjPath = "TestObj2"
	ob2.ParamSettings = make([]*usp_msg.Set_UpdateParamSetting, 2)
	ob2.ParamSettings[0] = &pr2
	ob2.ParamSettings[1] = &pr22

	ob3.ObjPath = "TestObj3"
	ob3.ParamSettings = make([]*usp_msg.Set_UpdateParamSetting, 2)
	ob3.ParamSettings[0] = &pr3
	ob3.ParamSettings[1] = &pr33

	setObj[0] = &ob1
	setObj[1] = &ob2
	setObj[2] = &ob3

	byteStream, err := CreateUspSetReqMsg(setObj, strconv.Itoa(id))

	if err == nil {

		err = sendCreatedMsg(byteStream)
		if err == nil {
			return
		}

		t.Errorf("Send Operate msg failed")
		return
	}
	t.Error("Create Operate Msg failed")
}

func TestAddObjects(t *testing.T) {

	addObj := make([]*usp_msg.Add_CreateObject, 3)

	var ob1 usp_msg.Add_CreateObject
	var pr1 usp_msg.Add_CreateParamSetting
	var pr11 usp_msg.Add_CreateParamSetting
	var ob2 usp_msg.Add_CreateObject
	var pr2 usp_msg.Add_CreateParamSetting
	var pr22 usp_msg.Add_CreateParamSetting
	var ob3 usp_msg.Add_CreateObject
	var pr3 usp_msg.Add_CreateParamSetting
	var pr33 usp_msg.Add_CreateParamSetting

	pr1.Param = "test1"
	pr1.Required = true
	pr1.Value = "vTest1"

	pr11.Param = "test11"
	pr11.Required = true
	pr11.Value = "vTest11"

	pr2.Param = "test2"
	pr2.Required = true
	pr2.Value = "vTest2"

	pr22.Param = "test22"
	pr22.Required = true
	pr22.Value = "vTest22"

	pr3.Param = "test3"
	pr3.Required = true
	pr3.Value = "vTest3"

	pr33.Param = "test33"
	pr33.Required = true
	pr33.Value = "vTest33"

	ob1.ObjPath = "TestingObjAdd1"
	ob1.ParamSettings = make([]*usp_msg.Add_CreateParamSetting, 2)
	ob1.ParamSettings[0] = &pr1
	ob1.ParamSettings[1] = &pr11

	ob2.ObjPath = "TestingObjAdd2"
	ob2.ParamSettings = make([]*usp_msg.Add_CreateParamSetting, 2)
	ob2.ParamSettings[0] = &pr2
	ob2.ParamSettings[1] = &pr22

	ob3.ObjPath = "TestingObjAdd3"
	ob3.ParamSettings = make([]*usp_msg.Add_CreateParamSetting, 2)
	ob3.ParamSettings[0] = &pr3
	ob3.ParamSettings[1] = &pr33

	addObj[0] = &ob1
	addObj[1] = &ob2
	addObj[2] = &ob3

	id++

	byteStream, err := CreateUspAddReqMsg(addObj, strconv.Itoa(id))

	if err == nil {
		err = sendCreatedMsg(byteStream)
		if err == nil {
			return
		}
		t.Errorf(err.Error())
		return
	}

	t.Errorf(err.Error())

}

func TestCreateUspDeleteReqMsg(t *testing.T) {
	var objs []string = make([]string, 5)
	objs[0] = "Device.ww"
	objs[1] = "Device.qq"
	objs[2] = "Device.aa"
	objs[3] = "Device.ss"
	objs[4] = "Device.zz"

	id++

	byteStream, err := CreateUspDeleteReqMsg(objs, strconv.Itoa(id))

	if err == nil {
		err = sendCreatedMsg(byteStream)
		if err == nil {
			return
		}
		t.Errorf(err.Error())
		return
	}

	t.Errorf(err.Error())
}

func TestCreateUspGetSupportedDmMsg(t *testing.T) {
	//
	var objs []string = make([]string, 5)
	objs[0] = "Device.ww"
	objs[1] = "Device.qq"
	objs[2] = "Device.aa"
	objs[3] = "Device.ss"
	objs[4] = "Device.zz"

	id++

	byteStream, err := CreateUspGetSupportedDmMsg(objs, true, true, true, strconv.Itoa(id))

	if err == nil {
		err = sendCreatedMsg(byteStream)
		if err == nil {
			return
		}
		t.Errorf(err.Error())
		return
	}

	t.Errorf(err.Error())
}

func TestCreateUspGetInstancesMsg(t *testing.T) {
	//
	var objs []string = make([]string, 5)
	objs[0] = "Device.ww"
	objs[1] = "Device.qq"
	objs[2] = "Device.aa"
	objs[3] = "Device.ss"
	objs[4] = "Device.zz"

	id++

	byteStream, err := CreateUspGetInstancesMsg(objs, true, strconv.Itoa(id))

	if err == nil {
		err = sendCreatedMsg(byteStream)
		if err == nil {
			return
		}
		t.Errorf(err.Error())
		return
	}

	t.Errorf(err.Error())

}

func TestCreateUspGetSupportedProtoMsg(t *testing.T) {

	byteStream, err := CreateUspGetSupportedProtoMsg("1.0,1.1", strconv.Itoa(id))

	if err == nil {
		err = sendCreatedMsg(byteStream)
		if err == nil {
			return
		}
		t.Errorf(err.Error())
		return
	}

	t.Errorf(err.Error())
}

func TestCreateUspNotifyEventMsg(t *testing.T) {

	var eventNotify usp_msg.Notify_Event_
	var event usp_msg.Notify_Event

	paramsMap := make(map[string]string, 3)
	paramsMap["event1"] = "Notification1"
	paramsMap["event2"] = "Notification2"
	paramsMap["event2"] = "Notification2"

	event.EventName = "EventNotification"
	event.ObjPath = "TestingEventObject"
	event.Params = paramsMap

	eventNotify.Event = &event
	byteStream, err := CreateUspNotifyEventMsg(&eventNotify, true, "SubscribeId1", strconv.Itoa(id))

	if err == nil {
		err = sendCreatedMsg(byteStream)
		if err == nil {
			return
		}
		t.Errorf(err.Error())
		return
	}

	t.Errorf(err.Error())
}

func TestCreateUspNotifyValueChangeEventMsg(t *testing.T) {

	var notifyValueChange usp_msg.Notify_ValueChange_
	var event usp_msg.Notify_ValueChange

	event.ParamPath = "TestingValChangeEventParam"
	event.ParamValue = "TestingValChangeEventParamValue"

	notifyValueChange.ValueChange = &event
	byteStream, err := CreateUspNotifyValueChangeMsg(&notifyValueChange, true, "SubscribeId1", strconv.Itoa(id))

	if err == nil {
		err = sendCreatedMsg(byteStream)
		if err == nil {
			return
		}
		t.Errorf(err.Error())
		return
	}

	t.Errorf(err.Error())
}

func TestCreateUspNotifyObjectCreationMsg(t *testing.T) {

	var notifyObjCreation usp_msg.Notify_ObjCreation
	var obj usp_msg.Notify_ObjectCreation
	uniqueKeys := make(map[string]string, 4)
	uniqueKeys["key1"] = "key11"
	uniqueKeys["key2"] = "key22"
	uniqueKeys["key3"] = "key33"
	uniqueKeys["key4"] = "key44"

	obj.ObjPath = "TestObjCreationPath"
	obj.UniqueKeys = uniqueKeys

	notifyObjCreation.ObjCreation = &obj

	byteStream, err := CreateUspNotifyObjectCreationMsg(&notifyObjCreation, true, "SubscribeId1", strconv.Itoa(id))

	if err == nil {
		err = sendCreatedMsg(byteStream)
		if err == nil {
			return
		}
		t.Errorf(err.Error())
		return
	}

	t.Errorf(err.Error())
}

func TestCreateUspNotifyObjectDeletionMsg(t *testing.T) {

	var notifyObjDeletion usp_msg.Notify_ObjDeletion
	var obj usp_msg.Notify_ObjectDeletion

	obj.ObjPath = "TestObjCreationPath"

	notifyObjDeletion.ObjDeletion = &obj

	byteStream, err := CreateUspNotifyObjectDeletionMsg(&notifyObjDeletion, true, "SubscribeId1", strconv.Itoa(id))

	if err == nil {
		err = sendCreatedMsg(byteStream)
		if err == nil {
			return
		}
		t.Errorf(err.Error())
		return
	}

	t.Errorf(err.Error())
}

func TestCreateUspNotifyOperationCompleteMsg(t *testing.T) {

	var notifyOprComplete usp_msg.Notify_OperComplete
	var oprComplete usp_msg.Notify_OperationComplete

	oprComplete.OperationResp = nil
	oprComplete.ObjPath = "TestOperationCompleteObjPath"

	notifyOprComplete.OperComplete = &oprComplete

	byteStream, err := CreateUspNotifyOperationCompleteMsg(&notifyOprComplete, true, "SubscribeId1", strconv.Itoa(id))

	if err == nil {
		err = sendCreatedMsg(byteStream)
		if err == nil {
			return
		}
		t.Errorf(err.Error())
		return
	}

	t.Errorf(err.Error())
}

func TestCreateUspNotifyOnBoardReqMsg(t *testing.T) {

	var notifyOnBoardReq usp_msg.Notify_OnBoardReq
	var onBoardReq usp_msg.Notify_OnBoardRequest

	onBoardReq.Oui = "01AB24"
	onBoardReq.ProductClass = "myname"
	onBoardReq.SerialNumber = "12341234"
	onBoardReq.AgentSupportedProtocolVersions = "1.0,1.1"

	notifyOnBoardReq.OnBoardReq = &onBoardReq

	byteStream, err := CreateUspNotifyOnboardReqMsg(&notifyOnBoardReq, true, "SubscribeId1", strconv.Itoa(id))

	if err == nil {

		err = sendCreatedMsg(byteStream)

		if err == nil {
			return
		}
		t.Errorf(err.Error())
		return
	}

	t.Errorf(err.Error())
}

func TestCreateErrorMsg(t *testing.T) {

	var errMsg usp_msg.Error
	errParams := make([]*usp_msg.Error_ParamError, 2)

	var param1 usp_msg.Error_ParamError
	var param2 usp_msg.Error_ParamError

	param1.ErrCode = 7001
	param1.ErrMsg = "Test Error Parameter 1"
	param1.ParamPath = "Device.ErrorParam Path 1"

	param2.ErrCode = 7002
	param2.ErrMsg = "Test Error Parameter 2"
	param2.ParamPath = "Device.ErrorParam Path 2"

	errParams[0] = &param1
	errParams[1] = &param2

	errMsg.ErrCode = 7000
	errMsg.ErrMsg = "Test Error String"
	errMsg.ParamErrs = errParams

	byteStream, err := CreateErrorMsg(&errMsg, strconv.Itoa(id))

	if err == nil {
		err = sendCreatedMsg(byteStream)
		if err == nil {
			return
		}
		t.Errorf(err.Error())
		return
	}

	t.Errorf(err.Error())
}

func TestChangeProcessingToController(t *testing.T) {
	ProcessUspMsg = ProcessControllerMsgs
}

func TestCreateUspAddRespMsg(t *testing.T) {
	//
	result := make([]*usp_msg.AddResp_CreatedObjectResult, 2)

	var res1 usp_msg.AddResp_CreatedObjectResult
	var res2 usp_msg.AddResp_CreatedObjectResult

	result[0] = &res1
	result[1] = &res2
	/***************/
	var opstat1 usp_msg.AddResp_CreatedObjectResult_OperationStatus

	var oprFailure usp_msg.AddResp_CreatedObjectResult_OperationStatus_OperFailure
	var operFailure usp_msg.AddResp_CreatedObjectResult_OperationStatus_OperationFailure
	res1.RequestedPath = "Test Add Resp result for path"
	res1.OperStatus = &opstat1

	operFailure.ErrCode = 1000
	operFailure.ErrMsg = "Test Add Resp Operation failure"

	opstat1.OperStatus = &oprFailure
	/******************/
	var opstat2 usp_msg.AddResp_CreatedObjectResult_OperationStatus

	var oprSuccess usp_msg.AddResp_CreatedObjectResult_OperationStatus_OperSuccess
	var operSuccess usp_msg.AddResp_CreatedObjectResult_OperationStatus_OperationSuccess

	opstat2.OperStatus = &oprSuccess
	oprSuccess.OperSuccess = &operSuccess

	operSuccess.InstantiatedPath = "Test Success Add Resp"

	paramErr := make([]*usp_msg.AddResp_ParameterError, 2)
	var paramErr1 usp_msg.AddResp_ParameterError
	var paramErr2 usp_msg.AddResp_ParameterError

	paramErr1.ErrCode = 1001
	paramErr1.ErrMsg = "Testing Add Resp Sucess but Param Error"
	paramErr1.Param = "Device.Tesing Onj.TestParam"

	paramErr2.ErrCode = 1002
	paramErr2.ErrMsg = "Testing Add Resp Sucess but Param 2 Error"
	paramErr2.Param = "Device.Tesing Onj.TestParam2"

	paramErr[0] = &paramErr1
	paramErr[1] = &paramErr2

	operSuccess.ParamErrs = paramErr
	operSuccess.UniqueKeys = make(map[string]string)
	operSuccess.UniqueKeys["key1"] = "vkey1"
	operSuccess.UniqueKeys["key2"] = "vkey2"

	id++
	byteStream, err := CreateUspAddRespMsg(strconv.Itoa(id), result)

	if err == nil {

		err = sendCreatedMsg(byteStream)

		if err == nil {
			return
		}
		t.Errorf("Send AddRespMsg Failed")
	}
	t.Errorf("Failed !!!! createUspAddRespMsg ")
}

func TestUspOperateResponseMsg(t *testing.T) {

	//createUspOperateRespMsg(id string, result []*usp_msg.OperateResp_OperationResult) ([]byte, error)
	result := make([]*usp_msg.OperateResp_OperationResult, 3)

	var opRes1 usp_msg.OperateResp_OperationResult
	var opRes2 usp_msg.OperateResp_OperationResult
	var opRes3 usp_msg.OperateResp_OperationResult

	result[0] = &opRes1
	result[1] = &opRes2
	result[2] = &opRes3

	var opResult1 usp_msg.OperateResp_OperationResult_ReqObjPath
	var opResult2 usp_msg.OperateResp_OperationResult_ReqOutputArgs
	var opResult3 usp_msg.OperateResp_OperationResult_CmdFailure

	opResult1.ReqObjPath = "opResult1 usp_msg.OperateResp_OperationResult_ReqObjPath"

	var args usp_msg.OperateResp_OperationResult_OutputArgs
	args.OutputArgs = make(map[string]string)

	args.OutputArgs["args1"] = "args usp_msg.OperateResp_OperationResult_OutputArgs"
	args.OutputArgs["args2"] = "args usp_msg.OperateResp_OperationResult_OutputArgs"

	opResult2.ReqOutputArgs = &args

	var cmdFailure usp_msg.OperateResp_OperationResult_CommandFailure
	cmdFailure.ErrCode = 1005
	cmdFailure.ErrMsg = "Test usp_msg.OperateResp_OperationResult_CommandFailure"

	opResult3.CmdFailure = &cmdFailure

	opRes1.OperationResp = &opResult1
	opRes1.ExecutedCommand = "Test Opr Command 1"

	opRes2.OperationResp = &opResult2
	opRes2.ExecutedCommand = "Test Opr Command 2"

	opRes3.OperationResp = &opResult3
	opRes3.ExecutedCommand = "Test Opr Command 3"

	id++

	byteStream, err := CreateUspOperateRespMsg(strconv.Itoa(id), result)

	if err == nil {

		err = sendCreatedMsg(byteStream)

		if err == nil {
			return
		}
		t.Errorf("Send OperateRespMsg Failed")
	}
	t.Errorf("Failed !!!! createUspOperateRespMsg ")
}

// createUspSetRespMsg
func TestCreateUspSetRespMsg(t *testing.T) {

	result := make([]*usp_msg.SetResp_UpdatedObjectResult, 2)

	var updateRes1 usp_msg.SetResp_UpdatedObjectResult
	var updateRes2 usp_msg.SetResp_UpdatedObjectResult

	result[0] = &updateRes1
	result[1] = &updateRes2

	var updateStatus1 usp_msg.SetResp_UpdatedObjectResult_OperationStatus
	var updateStatus2 usp_msg.SetResp_UpdatedObjectResult_OperationStatus

	var updatestatFail usp_msg.SetResp_UpdatedObjectResult_OperationStatus_OperFailure
	var updateStatSucc usp_msg.SetResp_UpdatedObjectResult_OperationStatus_OperSuccess

	var updatestatFailure usp_msg.SetResp_UpdatedObjectResult_OperationStatus_OperationFailure
	var updateStatSuccess usp_msg.SetResp_UpdatedObjectResult_OperationStatus_OperationSuccess

	updatestatFailure.ErrCode = 1001
	updatestatFailure.ErrMsg = "Test usp_msg.SetResp_UpdatedObjectResult_OperationStatus_OperationSuccess"

	upinstfail := make([]*usp_msg.SetResp_UpdatedInstanceFailure, 1)

	var updateinstFailure usp_msg.SetResp_UpdatedInstanceFailure
	upinstfail[0] = &updateinstFailure

	updateinstFailure.AffectedPath = "Test Affected PATH in update obj opr status"
	paramerr := make([]*usp_msg.SetResp_ParameterError, 1)
	var perr usp_msg.SetResp_ParameterError
	perr.ErrCode = 1000
	perr.ErrMsg = "Test set opr status failure param err"
	perr.Param = "Test set opr status failure param"

	updateinstFailure.ParamErrs = paramerr

	updatestatFail.OperFailure = &updatestatFailure
	updateStatSucc.OperSuccess = &updateStatSuccess

	updateStatus1.OperStatus = &updatestatFail
	updateStatus2.OperStatus = &updateStatSucc

	updateRes1.OperStatus = &updateStatus1
	updateRes1.RequestedPath = "Test usp_msg.SetResp_UpdatedObjectResult 1"

	updateRes2.OperStatus = &updateStatus2
	updateRes2.RequestedPath = "Test usp_msg.SetResp_UpdatedObjectResult 1"

	updateStatSuccess.UpdatedInstResults = make([]*usp_msg.SetResp_UpdatedInstanceResult, 2)
	var r1 usp_msg.SetResp_UpdatedInstanceResult
	var r2 usp_msg.SetResp_UpdatedInstanceResult

	updateStatSuccess.UpdatedInstResults[0] = &r1
	updateStatSuccess.UpdatedInstResults[1] = &r2
	r1.AffectedPath = "Test Set Obj Affected PATH"
	r1.ParamErrs = paramerr

	r1.UpdatedParams = make(map[string]string)
	r1.UpdatedParams["param"] = "mapParamVal"

	r2.AffectedPath = "Test Set Obj Affected PATH -1"
	r2.ParamErrs = paramerr

	r2.UpdatedParams = make(map[string]string)
	r2.UpdatedParams["param"] = "mapParamVal"
	id++
	byteStream, err := CreateUspSetRespMsg(strconv.Itoa(id), result)

	if err == nil {

		err = sendCreatedMsg(byteStream)

		if err == nil {
			return
		}
		t.Errorf("Send SetRespMsg Failed")
	}
	t.Errorf("Failed !!!! createUspSetRespMsg ")
}

func TestCreateUspGetRespMsg(t *testing.T) {

	result := make([]*usp_msg.GetResp_RequestedPathResult, 2)

	var r1 usp_msg.GetResp_RequestedPathResult
	var r2 usp_msg.GetResp_RequestedPathResult

	r1.ErrCode = 7000
	r1.ErrMsg = "Test GetResp_RequestedPathResult r1"
	r1.RequestedPath = "Test GetResp_RequestedPathResult r1 PATH"
	rr1 := make([]*usp_msg.GetResp_ResolvedPathResult, 2)
	var gr1 usp_msg.GetResp_ResolvedPathResult
	var gr2 usp_msg.GetResp_ResolvedPathResult

	rr1[0] = &gr1
	rr1[0] = &gr2

	r1.ResolvedPathResults = rr1

	r2.ErrCode = 7000
	r2.ErrMsg = "Test GetResp_RequestedPathResult r2"
	r2.RequestedPath = "Test GetResp_RequestedPathResult r2 PATH"
	r1r1 := make([]*usp_msg.GetResp_ResolvedPathResult, 2)
	var g1r1 usp_msg.GetResp_ResolvedPathResult
	var g1r2 usp_msg.GetResp_ResolvedPathResult

	r1r1[0] = &g1r1
	r1r1[0] = &g1r2

	r2.ResolvedPathResults = r1r1
	id++
	byteStream, err := CreateUspGetRespMsg(result, strconv.Itoa(id))

	if err == nil {

		err = sendCreatedMsg(byteStream)

		if err == nil {
			return
		}
		t.Errorf("Send GetRespMsg Failed")
	}
	t.Errorf("Failed !!!! createUspGetRespMsg ")
}

func TestCcreateUspDeleteRespMsg(t *testing.T) {

	result := make([]*usp_msg.DeleteResp_DeletedObjectResult, 2)
	var dr1 usp_msg.DeleteResp_DeletedObjectResult
	var dr2 usp_msg.DeleteResp_DeletedObjectResult
	result[0] = &dr1
	result[1] = &dr2

	var dros usp_msg.DeleteResp_DeletedObjectResult_OperationStatus

	var drosFail usp_msg.DeleteResp_DeletedObjectResult_OperationStatus_OperFailure
	var drosFailure usp_msg.DeleteResp_DeletedObjectResult_OperationStatus_OperationFailure

	drosFailure.ErrCode = 1000
	drosFailure.ErrMsg = "Test Del Resp result operation status failure"
	drosFail.OperFailure = &drosFailure
	dros.OperStatus = &drosFail
	dr1.OperStatus = &dros
	dr1.RequestedPath = "Test Del Resp path 1"

	var dros2 usp_msg.DeleteResp_DeletedObjectResult_OperationStatus

	var drosSucc usp_msg.DeleteResp_DeletedObjectResult_OperationStatus_OperSuccess
	var drosSuccess usp_msg.DeleteResp_DeletedObjectResult_OperationStatus_OperationSuccess

	var affectedPath []string = make([]string, 1)
	affectedPath[0] = "Test Del Resp success Affected path"
	drosSuccess.AffectedPaths = affectedPath

	upe := make([]*usp_msg.DeleteResp_UnaffectedPathError, 2)
	var upe1 usp_msg.DeleteResp_UnaffectedPathError
	var upe2 usp_msg.DeleteResp_UnaffectedPathError
	upe[0] = &upe1
	upe[1] = &upe2

	upe1.ErrCode = 1001
	upe1.ErrMsg = "Test Err DeleteResp_UnaffectedPathError 1"
	upe1.UnaffectedPath = "Test Err DeleteResp_UnaffectedPath 1"
	upe2.ErrCode = 1002
	upe2.ErrMsg = "Test Err DeleteResp_UnaffectedPathError 2"
	upe2.UnaffectedPath = "Test Err DeleteResp_UnaffectedPath 2"

	drosSuccess.UnaffectedPathErrs = upe

	drosSucc.OperSuccess = &drosSuccess
	dros.OperStatus = &drosSucc
	dr2.OperStatus = &dros2
	dr2.RequestedPath = "Test Del Resp path 2"

	id++
	byteStream, err := CreateUspDeleteRespMsg(result, strconv.Itoa(id))

	if err == nil {

		err = sendCreatedMsg(byteStream)

		if err == nil {
			return
		}
		t.Errorf("Send DeleteRespMsg Failed")
	}
	t.Errorf("Failed !!!! createUspDeleteRespMsg ")
}

func TestCcreateUspGetSupportedDmRespMsg(t *testing.T) {

	result := make([]*usp_msg.GetSupportedDMResp_RequestedObjectResult, 1)

	var r1 usp_msg.GetSupportedDMResp_RequestedObjectResult

	result[0] = &r1

	r1.ErrMsg = "Test usp_msg.GetSupportedDMResp_RequestedObjectResult Err Msg "
	r1.ReqObjPath = "Test path Supported DM resp"
	supObj := make([]*usp_msg.GetSupportedDMResp_SupportedObjectResult, 1)

	r1.SupportedObjs = supObj

	var so1 usp_msg.GetSupportedDMResp_SupportedObjectResult

	supObj[0] = &so1

	so1.IsMultiInstance = false
	supcmd := make([]*usp_msg.GetSupportedDMResp_SupportedCommandResult, 2)
	var supcmd1 usp_msg.GetSupportedDMResp_SupportedCommandResult
	var supcmd2 usp_msg.GetSupportedDMResp_SupportedCommandResult

	so1.SupportedCommands = supcmd

	supcmd1.CommandName = "Test usp_msg.GetSupportedDMResp_SupportedCommandResult 1"
	supcmd1.InputArgNames = make([]string, 2)
	supcmd1.InputArgNames[0] = "Input Arg 1"
	supcmd1.InputArgNames[1] = "Input Arg 2"
	supcmd1.OutputArgNames = make([]string, 2)
	supcmd1.OutputArgNames[0] = "output Arg 1"
	supcmd1.OutputArgNames[1] = "output Arg 2"

	supcmd2.CommandName = "Test usp_msg.GetSupportedDMResp_SupportedCommandResult 1"
	supcmd2.InputArgNames = make([]string, 2)
	supcmd2.InputArgNames[0] = "Input Arg 11"
	supcmd2.InputArgNames[1] = "Input Arg 22"
	supcmd2.OutputArgNames = make([]string, 2)
	supcmd2.OutputArgNames[0] = "output Arg 11"
	supcmd2.OutputArgNames[1] = "output Arg 22"

	so1.SupportedObjPath = "test SupportedObjPath 1"

	sparamreres := make([]*usp_msg.GetSupportedDMResp_SupportedParamResult, 2)

	var sparamreres1 usp_msg.GetSupportedDMResp_SupportedParamResult
	var sparamreres2 usp_msg.GetSupportedDMResp_SupportedParamResult

	sparamreres[0] = &sparamreres1
	sparamreres[1] = &sparamreres2
	sparamreres1.ParamName = "Test supported DM ParamMsg 1"
	sparamreres2.ParamName = "Test supported DM ParamMsg 2"

	var paramtype1 usp_msg.GetSupportedDMResp_ParamAccessType = 2
	var paramtype2 usp_msg.GetSupportedDMResp_ParamAccessType = 2

	sparamreres1.Access = paramtype1
	sparamreres2.Access = paramtype2

	so1.SupportedParams = sparamreres

	//var objtype usp_msg.GetSupportedDMResp_ObjAccessType
	//objtype = 2
	so1.Access = 2

	id++
	byteStream, err := CreateUspGetSupportedDmRespMsg(result, strconv.Itoa(id))

	if err == nil {

		err = sendCreatedMsg(byteStream)

		if err == nil {
			return
		}
		t.Errorf("Send GetSupportedDmRespMsg Failed")
	}
	t.Errorf("Failed !!!! createUspGetSupportedDmRespMsg ")
}

func TestCcreateUspGetInstancesRespMsg(t *testing.T) {

	result := make([]*usp_msg.GetInstancesResp_RequestedPathResult, 2)
	var r1 usp_msg.GetInstancesResp_RequestedPathResult
	var r2 usp_msg.GetInstancesResp_RequestedPathResult

	result[0] = &r1
	result[1] = &r2

	r1.ErrCode = 1001
	r1.ErrMsg = "Test GetInstancesRespMsg Err 1"
	r1.RequestedPath = "Test GetInstancesRespMsg Path 1"

	r2.ErrCode = 1001
	r2.ErrMsg = "Test GetInstancesRespMsg Err 2"
	r2.RequestedPath = "Test GetInstancesRespMsg Path 2"

	id++
	byteStream, err := CreateUspGetInstancesRespMsg(result, strconv.Itoa(id))

	if err == nil {

		err = sendCreatedMsg(byteStream)

		if err == nil {
			return
		}
		t.Errorf("Send GetInstancesRespMsg Failed")
	}
	t.Errorf("Failed !!!! createUspGetInstancesRespMsg ")
}

func TestCcreateUspGetSupportedProtoRespMsg(t *testing.T) {
	id++
	byteStream, err := CreateUspGetSupportedProtoRespMsg("1.0,1.1", strconv.Itoa(id))

	if err == nil {

		err = sendCreatedMsg(byteStream)

		if err == nil {
			return
		}
		t.Errorf("Send GetSupportedProtoRespMsg Failed")
	}
	t.Errorf("Failed !!!! createUspGetSupportedProtoRespMsg ")
}
func TestExit(t *testing.T) {

	time.Sleep(5 * time.Second)
	clientKill <- true
}

var ProcessUspMsg func([]byte) ([]byte, error)

func wsMsgHandler(conn *websocket.Conn) {

	for {
		// get msg
		_, byteStream, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("Agent received Msg from controller")

		// Decode usp records first

		//msgStream, err := usprecords.GetUspMsgStreamFromRecord(byteStream)
		msgStream, err := GetUspMsgStreamFromRecord(byteStream)

		if err != nil {
			log.Printf("Could not retrieve USP Msg from USP Record")
			continue
		}
		_, uspProcessErr := ProcessUspMsg(msgStream)

		if uspProcessErr != nil {
			log.Printf("Error ProcessUspMsg")
		}
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Hello World, HTTP server Home")
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
	}

	log.Println("Client Connected")

	// following code echos the received msg
	if err != nil {
		log.Println(err)
	}

	wsMsgHandler(ws)
	ws.Close()
}

func wsServer(ch chan bool, serverPort string) {

	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsEndpoint)

	port, err := strconv.Atoi(serverPort)

	if err != nil {
		usage(nil)
		return
	}

	if port <= 0 {

		log.Printf("Using default serverPort as 8080")
		port = 8080
	}
	if ch != nil {
		ch <- true
	}
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(int(port)), nil))
}

func wsClient(srv string, port string, success chan bool, exit chan bool) {
	var connErr error = nil

	for ok := true; ok; {

		log.Printf("Trying to connect to server using websocket")
		time.Sleep(time.Second)

		client, _, connErr = websocket.DefaultDialer.Dial("ws://"+srv+":"+port+"/ws", nil)
		if connErr != nil {
			log.Printf("Error!!!!!! dialing ws %v ", connErr)
		} else {
			if success != nil {
				success <- true
			}
			break
		}
	}

	if connErr == nil {
		<-exit // wait till we are asked to exit
	} else {
		if success != nil {
			success <- false
		}
	}

	log.Printf("web socket client exiting")
	client.Close()
	client = nil
}
