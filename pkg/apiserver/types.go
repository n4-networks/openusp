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
