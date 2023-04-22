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

package cntlr

import "strconv"

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
		return "Del_" + strconv.Itoa(delCount)
	case MsgTypeNotify:
		notifyCount++
		return "NOTIFY_" + strconv.Itoa(notifyCount)
	default:
		defaultCount++
		return "DEFAULT_" + strconv.Itoa(defaultCount)
	}
	return "InvalidMsgId"
}
