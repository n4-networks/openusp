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
