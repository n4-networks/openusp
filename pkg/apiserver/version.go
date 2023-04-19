package apiserver

import (
	"strconv"
)

type verType struct {
	major uint
	minor uint
}

var buildtime string

var (
	ver verType = verType{
		major: 1,
		minor: 0,
	}
)

func getVer() string {

	v := strconv.FormatUint(uint64(ver.major), 10) + "." + strconv.FormatUint(uint64(ver.minor), 10)
	if buildtime != "" {
		v = v + "." + buildtime
	}
	return v
}
