package apiserver

type devInfo struct {
	productClass string
	manufacturer string
	modelName    string
}

type object struct {
	path   string
	params map[string]string
}

type dmCmd struct {
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

type DmObject struct {
	Path          string    `json:"path"`
	MultiInstance bool      `json:"multi_instance"`
	Access        string    `json:"access"`
	Params        []dmParam `json:"params"`
	Events        []dmEvent `json:"events"`
	Cmds          []dmCmd   `json:"cmds"`
}

type param struct {
	path  string
	value string
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
type cfgInstance struct {
	path   string
	level  int
	key    string
	params map[string]string
}
type cfgParamNode struct {
	path   string
	params map[string]string
}

type addInstInfo struct {
	path   string
	parent string
	params map[string]string
}
