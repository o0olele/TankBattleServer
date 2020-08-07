package env

import (
	"encoding/json"
	"io/ioutil"

	"github.com/golang/glog"
)

var configData map[string]map[string]string

func Load(path string) bool {
	file, err := ioutil.ReadFile(path)
	if nil != err {
		glog.Error("[Config] Read Error, Path=", path, ",", err)
		return false
	}

	err = json.Unmarshal(file, &configData)
	if nil != err {
		glog.Error("[Config] Analys Error, Path=", path, ",", err)
		return false
	}

	return true
}

func Get(table, key string) string {
	t, ok := configData[table]
	if !ok {
		return ""
	}

	v, ok := t[key]
	if !ok {
		return ""
	}

	return v
}
