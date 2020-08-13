package common

import (
	"encoding/json"
	"fmt"

	"github.com/golang/glog"
)

const ()

type ErrorMsg struct {
	code   uint
	errmsg string
}

func JsonMsg(arg ...interface{}) string {
	a, err := json.Marshal(arg)
	if err != nil {
		glog.Error("Marshal json msg error")
	}

	fmt.Println(string(a))
	return string(a)
}

func ErrJsonMsg(code uint, errmsg string) string {
	emsg := ErrorMsg{
		code:   code,
		errmsg: errmsg,
	}
	a, err := json.Marshal(emsg)
	if err != nil {
		glog.Error("Marshal json msg error")
	}
	return string(a)
}
