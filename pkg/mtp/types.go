package mtp

type dmCommand struct {
	name    string
	inputs  []string
	outputs []string
}
type dmEvent struct {
	name string
	args []string
}
type dmParam struct {
	name   string
	access string
}

type dmObject struct {
	path          string
	multiInstance bool
	access        string
	params        []dmParam
	events        []dmEvent
	cmds          []dmCommand
}

type param struct {
	path  string
	value string
}

type instance struct {
	path       string
	uniqueKeys map[string]string
}

type cfgInstance struct {
	path   string
	level  int
	params map[string]string
}

type cfgParamNode struct {
	path   string
	params map[string]string
}

type NotifyType byte

const (
	NotifyEvent NotifyType = iota
	NotifyValueChange
	NotifyObjCreation
	NotifyObjDeletion
	NotifyOpComplete
	NotifyOnBoardReq
)

type event struct {
	path   string
	name   string
	params map[string]string
}

type valueChange struct {
	paramPath  string
	paramValue string
}

type objectCreation struct {
	path       string
	uniqueKeys map[string]string
}

type objectDeletion struct {
	path string
}
type opFailure struct {
	errCode uint32
	errMsg  string
}
type operationComplete struct {
	path       string
	cmdName    string
	cmdKey     string
	outArg     map[string]string
	cmdFailure opFailure
}
type onBoardReq struct {
	oui          string
	productClass string
	serialNum    string
	protoVer     string
}

type notification struct {
	subscriptionId string
	sendResp       bool
	nType          NotifyType
	evt            *event
	valChange      *valueChange
	objCreation    *objectCreation
	objDeletion    *objectDeletion
	opComplete     *operationComplete
	onBoard        *onBoardReq
}
type agentDeviceInfo struct {
	productClass string
	manufacturer string
	modelName    string
}
